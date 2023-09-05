package client

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "request_count",
		Help: "The total number of requests",
	}, []string{"method", "code", "succ"})
)

func NewRequestCount(method string, code int, succ bool) prometheus.Counter {
	return requestCount.WithLabelValues(method, fmt.Sprintf("%d", code), fmt.Sprintf("%t", succ))
}
