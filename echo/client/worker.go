package client

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type record func([]byte)

type worker struct {
	ratelimit int
	client    *http.Client
	recordFn  record
}

func (w *worker) call(req *http.Request, stop <-chan struct{}) {
	buff := bytes.Buffer{}
	for {
		count := w.ratelimit
		for {
			select {
			case <-time.After(1 * time.Second):
				fmt.Printf("timeout\n")
				break
			case <-stop:
				fmt.Printf("stop\n")
				return
			default:
				if count > 0 {
					resp, err := w.client.Do(req)
					if err != nil {
						bodyBuff := bytes.Buffer{}
						bodyBuff.ReadFrom(req.Body)
						fmt.Printf("Fail to request %s %s %s: %s\n", req.Method, req.URL, bodyBuff.String(), err)
					} else {
						buff.Reset()
						buff.ReadFrom(resp.Body)
						w.recordFn(buff.Bytes())
					}
				}
			}
		}
	}
}
