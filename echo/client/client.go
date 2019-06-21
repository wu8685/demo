package client

import (
	"bytes"
	"fmt"
	"github.com/golang/glog"
	"net/http"
	"sync"
	"time"
)

const DefaultQPSperWorker = 50

func Start(qps int, workerNum *int, method, url, body string, stop <-chan struct{}) {
	num := 0
	avgQps := DefaultQPSperWorker
	if workerNum != nil {
		num = *workerNum
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

		worker := &worker{ratelimit: workerQps, client: &http.Client{}, recordFn: recordResponse}
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker.call(req, stop)
		}()
	}

	<-stop
	select {
	case <-time.After(2 * time.Second):
	default:
		wg.Wait()
	}
	return
}

func recordResponse(content []byte) {
	glog.V(0).Infof("response: %s", content)
}
