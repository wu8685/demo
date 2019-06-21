package main

import (
	"flag"
	"fmt"
	"github.com/wu8685/demo/echo/client"
	"os"
)

var method, url, body string
var qps int
var worker *int

func main() {
	f := flag.NewFlagSet("echo-client", flag.ExitOnError)
	f.StringVar(&method, "m", "GET", "HTTP call method")
	f.StringVar(&url, "url", "http://localhost:8080/healthz", "HTTP call url")
	f.StringVar(&body, "d", "", "HTTP call body")
	f.IntVar(&qps, "qps", 100,"HTTP call QPS")
	f.Int("worker", 0, fmt.Sprintf("HTTP call worker number, if not provided, will be calculated by given QPS and %d QPS per worker", client.DefaultQPSperWorker))

	if err := f.Parse(os.Args[1:]); err != nil {
		fmt.Printf("fail to parse command line: %s\n", err)
		return
	}

	if f := flag.CommandLine.Lookup("worker"); f == nil {
		worker = nil
	}

	stop := make(chan struct{})
	client.Start(qps, worker, method, url, body, stop)
}
