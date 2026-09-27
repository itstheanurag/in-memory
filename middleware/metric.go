package middleware

import (
	"fmt"
	"time"

	"github.com/itstheanurag/in-memory/store"
)

type MetricsMiddleware struct {
	inner store.Storer

	getCalls    int
	setCalls    int
	deleteCalls int
	lenCalls    int
	keyCalls    int

	getMisses       int
	totalGetLatency time.Duration
	totalSetLatency time.Duration
}

func NewMetricsMiddleware(inner store.Storer) *MetricsMiddleware {
	return &MetricsMiddleware{inner: inner}
}

func (mm *MetricsMiddleware) Get(key string) (string, error) {
	start := time.Now()
	mm.getCalls++
	val, err := mm.inner.Get(key)
	mm.totalGetLatency += time.Since(start) // accumulate the total Get latency

	if err != nil {
		mm.getMisses++
	}

	return val, err
}

func (mm *MetricsMiddleware) Set(key, value string) error {
	start := time.Now()
	err := mm.inner.Set(key, value)
	mm.setCalls++
	mm.totalSetLatency += time.Since(start) // accumulate the total set latency
	return err
}

func (mm *MetricsMiddleware) Delete(key string) error {
	mm.deleteCalls++
	return mm.inner.Delete(key)
}

func (mm *MetricsMiddleware) Len() int {
	mm.lenCalls++
	return mm.inner.Len()
}

func (mm *MetricsMiddleware) Keys() []string {
	mm.keyCalls++
	return mm.inner.Keys()
}

func (mm *MetricsMiddleware) Report() {
	fmt.Printf("GET calls: %d, (missed %d)\n", mm.getCalls, mm.getMisses)
	fmt.Printf("SET calls: %d, (latency: %s)\n", mm.setCalls, mm.totalSetLatency)
	fmt.Printf("DELETE calls: %d\n", mm.deleteCalls)
	fmt.Printf("LEN calls: %d\n", mm.lenCalls)

	if mm.getCalls > 0 {
		avg := mm.totalGetLatency / time.Duration(mm.getCalls)
		fmt.Printf("AVG GET latency: %s\n", avg)
	}

	if mm.setCalls > 0 {
		avg := mm.totalSetLatency / time.Duration(mm.setCalls)
		fmt.Printf("AVG SET latency: %s\n", avg)
	}
}
