package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"sync"
	"time"
)

const DefaultQPSperWorker = 50

var Logger *log.Logger

func Start(qps int, workerNum int, method, url, body string, timeout int64, stop <-chan struct{}) chan struct{} {
	num := 0
	avgQps := DefaultQPSperWorker
	if workerNum > 0 {
		num = workerNum
		avgQps = qps / num
	} else {
		num = qps / DefaultQPSperWorker
		if qps%DefaultQPSperWorker > 0 {
			num++
		}
	}

	wg := sync.WaitGroup{}

	rest := qps
	buff := bytes.NewBuffer([]byte(body))
	req, err := http.NewRequest(method, url, buff)
	if err != nil {
		fmt.Printf("fail to create request: %s\n", err)
	}

	for i := 0; i < num; i++ {
		workerQps := rest
		if workerQps > avgQps {
			workerQps = avgQps
		}
		rest -= workerQps

		var client *http.Client
		if timeout > 0 {
			timeoutDuration := time.Duration(timeout) * time.Millisecond
			client = &http.Client{Timeout: timeoutDuration}
		} else {
			client = &http.Client{}
		}
		//client := &CmdClient{method: method, url: url, body: body}
		worker := &worker{name: fmt.Sprintf("worker-%d", i), ratelimit: workerQps, client: client, recordFn: recordResponse}
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker.call(req, stop)
		}()
	}

	workerStop := make(chan struct{})
	go func() {
		<-stop

		select {
		case <-time.After(1 * time.Second):
		default:
			wg.Wait()
		}

		close(workerStop)
	}()

	return workerStop
}

func recordResponse(timestamp int64, content []byte) {
	Logger.Printf("%d %s\n", timestamp, content)
}

type CmdClient struct {
	method string
	url    string
	body   string
}

func (c *CmdClient) Do(req *http.Request) (*http.Response, error) {
	var cmd *exec.Cmd
	if len(c.body) == 0 {
		cmd = exec.Command("curl", "-X", c.method, c.url)
	} else {
		cmd = exec.Command("curl", "-X", c.method, c.url, "-d", c.body)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	resp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer(output)),
	}

	return resp, nil
}
