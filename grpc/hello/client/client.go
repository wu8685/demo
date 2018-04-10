package main

import (
	"context"
	"flag"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	"github.com/wu8685/demo/grpc/hello/proto"
	"google.golang.org/grpc/metadata"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "testdata/ca.pem", "The file containning the CA root cert file")
	serverAddr         = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		var sn string
		if *serverHostOverride != "" {
			sn = *serverHostOverride
		}
		var creds credentials.TransportCredentials
		if *caFile != "" {
			var err error
			creds, err = credentials.NewClientTLSFromFile(*caFile, sn)
			if err != nil {
				grpclog.Fatalf("Failed to create TLS credentials %v", err)
			}
		} else {
			creds = credentials.NewClientTLSFromCert(nil, sn)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to connect server: %v", *serverAddr)
	}
	defer conn.Close()

	client := proto.NewPersonManagerClient(conn)

	log.Printf("start recording person")
	recordPerson(client)

	log.Printf("start chatting")
	chat(client)

	log.Printf("start listing")
	listPerson(client)

	log.Printf("start finding")
	findPerson(client)
}

// Unary RPC
func findPerson(client proto.PersonManagerClient) {
	res, err := client.FindPerson(context.Background(), &proto.Person{"Lilei", 18})
	if err != nil {
		log.Fatalf("err when finding person: %s", err.Error())
	}

	log.Printf("find person: %t, message: %s", res.Succ, res.Message)
}

// Server streaming RPC
func listPerson(client proto.PersonManagerClient) {
	stream, err := client.ListPersons(context.Background(), &proto.Empty{})
	if err != nil {
		log.Fatalf("err when call listPersons: %s", err.Error())
	}
	for {
		result, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("err when list person: %s", err.Error())
			return
		}

		log.Printf("list person: %t, message: %s", result.Succ, result.Message)
	}
}

// Client streaming RPC
func recordPerson(client proto.PersonManagerClient) {
	md := metadata.MD{}
	md["test-key"] = []string{"test-value"}
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	stream, err := client.RecordPerson(ctx, grpc.Header(&md))
	if err != nil {
		log.Fatalf("err when call recordPersonL: %s", err.Error())
	}

	persons := []*proto.Person{
		{"HanMeimei", 17},
		{"Lucy", 16},
		{"Luly", 16},
	}
	for _, p := range persons {
		err = stream.Send(p)
		if err != nil {
			log.Fatalf("err when send in recordPerson: %s", err.Error())
		}
	}

	result, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("err when record persons: %s", err.Error())
	}

	log.Printf("record person: %t, message: %s", result.Succ, result.Message)
}

// Bidirectional streaming RPC
func chat(client proto.PersonManagerClient) {
	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatalf("err when call chat: %s", err.Error())
	}

	fin := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(fin)
				break
			}
			if err != nil {
				log.Fatalf("err when receive in chat: %s", err.Error())
			}

			log.Printf("receive in chat: %t, message: %s", res.Succ, res.Message)
		}
	}()

	persons := []*proto.Person{
		{"Jim", 18},
		{"Tom", 18},
		{"Ann", 16},
	}
	for _, p := range persons {
		err = stream.Send(p)
		if err != nil {
			log.Fatalf("err when send in recordPerson: %s", err.Error())
		}
	}
	stream.CloseSend()
	<-fin
}
