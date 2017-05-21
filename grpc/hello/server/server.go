package main

import (
	"fmt"
	"io"
	"flag"
	"net"
	"log"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/credentials"

	"demo/grpc/hello/proto"
)

var (
	port = flag.Int("port", 10000, "server port")
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "testdata/server1.pem", "The TLS cert file")
	keyFile    = flag.String("key_file", "testdata/server1.key", "The TLS key file")
)

func main() {
	flag.Parse()

	server := newMyServer()
	server.start(*port)
}

func newMyServer() *MyServer {
	return &MyServer{
		persons: map[string]*proto.Person{},
	}
}

type MyServer struct {
	grpcServer *grpc.Server
	persons map[string]*proto.Person
}

func (s *MyServer) start(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("fail to start my server: %s", err.Error())
	}

	var opts []grpc.ServerOption
	if *tls {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			grpclog.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	s.grpcServer = grpc.NewServer(opts...)
	proto.RegisterPersonManagerServer(s.grpcServer, s)

	log.Printf("start server :%d", port)
	s.grpcServer.Serve(lis)
}

// type 1: A simple RPC where the client sends a request to the server using the stub
// and waits for a response to come back, just like a normal function call.
func (s *MyServer) FindPerson(ctx context.Context, p *proto.Person) (*proto.Result, error) {
	if person, exist := s.persons[p.Name]; !exist {
		return &proto.Result{false, fmt.Sprintf("a person names %s does not exist", p.Name)}, nil
	} else {
		return &proto.Result{true, fmt.Sprintf("name: %s, age: %d", person.Name, person.Age)}, nil
	}
}
// type 2: A server-side streaming RPC where the client sends a request to the server
// and gets a stream to read a sequence of messages back. The client reads from the returned stream
// until there are no more messages. As you can see in our example,
// you specify a server-side streaming method by placing the stream keyword before the response type.
func (s *MyServer) ListPersons(_ *proto.Empty, stream proto.PersonManager_ListPersonsServer) error {
	for _, person := range s.persons {
		stream.Send(&proto.Result{true, fmt.Sprintf("name: %s, age: %d", person.Name, person.Age)})
	}
	return nil
}
// type 3: A client-side streaming RPC where the client writes a sequence of messages
// and sends them to the server, again using a provided stream.
// Once the client has finished writing the messages, it waits for the server to read them all
// and return its response. You specify a client-side streaming method
// by placing the stream keyword before the request type.
func (s *MyServer) RecordPerson(stream proto.PersonManager_RecordPersonServer) error {
	for {
		person, err := stream.Recv()
		if err == io.EOF {

			return stream.SendAndClose(&proto.Result{true, "finish recording"})
		}
		if err != nil {
			return err
		}

		s.persons[person.Name] = person
	}
}
// type 4: A bidirectional streaming RPC where both sides send a sequence of messages
// using a read-write stream. The two streams operate independently,
// so clients and servers can read and write in whatever order they like: for example,
// the server could wait to receive all the client messages before writing its responses,
// or it could alternately read a message then write a message, or some other combination of reads and writes.
// The order of messages in each stream is preserved. You specify this type of method by placing
// the stream keyword before both the request and the response.
func (s *MyServer) Chat(stream proto.PersonManager_ChatServer) error {
	for _, person := range s.persons {
		stream.Send(&proto.Result{true, fmt.Sprintf("name: %s, age: %d", person.Name, person.Age)})
	}

	for {
		person, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		s.persons[person.Name] = person
		stream.Send(&proto.Result{true, fmt.Sprintf("name: %s, age: %d", person.Name, person.Age)})
	}
}


