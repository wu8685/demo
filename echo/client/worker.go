package client

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type record func(int64, []byte)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type worker struct {
	name      string
	ratelimit int
	client    httpClient
	recordFn  record
}

func (w *worker) call(req *http.Request, stop <-chan struct{}) {
	buff := bytes.Buffer{}
	ticker := time.NewTicker(1. * time.Second)

	for {
		count := w.ratelimit
		failed := 0
		timestamp := time.Now().UnixNano() / 1e6
	LOOP:
		for {
			select {
			case <-ticker.C:
				fmt.Printf("%s another loop, request: %d, failed: %d\n", w.name, w.ratelimit-count, failed)
				break LOOP
			case <-stop:
				fmt.Printf("stop worker %s\n", w.name)
				return
			default:
				if count > 0 {
					count--
					resp, err := w.client.Do(req)
					if err != nil {
						NewRequestCount(req.Method, 0, false).Inc()
						bodyBuff := bytes.Buffer{}
						bodyBuff.ReadFrom(req.Body)
						fmt.Printf("%s fails to request %s %s %s: %s\n", w.name, req.Method, req.URL, bodyBuff.String(), err)
						failed++
					} else {
						NewRequestCount(req.Method, resp.StatusCode, true).Inc()
						buff.Reset()
						buff.ReadFrom(resp.Body)
						w.recordFn(timestamp, buff.Bytes())
					}
				}
			}
		}
	}
}
