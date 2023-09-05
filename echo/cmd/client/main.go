package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/wu8685/demo/echo/client"
)

var method, url, body string
var qps int
var worker int
var timeout int64
var logPath string

func main() {
	f := flag.NewFlagSet("echo-client", flag.ExitOnError)

	f.StringVar(&method, "m", "GET", "HTTP call method")
	f.StringVar(&url, "url", "http://localhost:8080/healthz", "HTTP call url")
	f.StringVar(&body, "d", "", "HTTP call body")
	f.IntVar(&qps, "qps", 100, "HTTP call QPS")
	f.IntVar(&worker, "worker", 0, fmt.Sprintf("HTTP call worker number, if not provided, will be calculated by given QPS and %d QPS per worker", client.DefaultQPSperWorker))
	f.Int64Var(&timeout, "timeout", 500, "HTTP call timeout (unit milliSecond)")
	f.StringVar(&logPath, "log", "/app/log/qps.log", "QPS log file path")

	if err := f.Parse(os.Args[1:]); err != nil {
		fmt.Printf("fail to parse command line: %s\n", err)
		return
	}

	logPath = strings.TrimRight(logPath, "/")

	_, err := os.Stat(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			parts := strings.Split(logPath, "/")
			if err := os.MkdirAll(strings.Join(parts[:len(parts)-1], "/"), 0777); err != nil {
				log.Printf("fail to create path %s: %s\n", logPath, err)
			}
		} else {
			log.Printf("fail to find stat of log path %s: %s\n", logPath, err)
		}
	}

	logfile, err := os.OpenFile(logPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
	}
	defer logfile.Close()

	// remove the timestamp prefix in logs
	client.Logger = log.New(logfile, "", log.LstdFlags&^(log.Ldate|log.Ltime))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	stop := make(chan struct{})
	workerStop := client.Start(qps, worker, method, url, body, timeout, stop)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	sig := <-sigs
	fmt.Printf("handle signal %s\n", sig)
	close(stop)

	select {
	case <-workerStop:
	case <-time.After(2 * time.Second):
	}
	fmt.Printf("stop\n")
}
