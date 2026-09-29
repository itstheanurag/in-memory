package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/itstheanurag/in-memory/store"
)

type MetricsMiddleware struct {
	inner store.Storer

	mu sync.Mutex

	getCalls    int
	setCalls    int
	deleteCalls int
	lenCalls    int
	keyCalls    int
	incrCalls   int
	decrCalls   int

	getMisses        int
	totalGetLatency  time.Duration
	totalSetLatency  time.Duration
	totalIncrLatency time.Duration
	totalDecrLatency time.Duration
}

func NewMetricsMiddleware(inner store.Storer) *MetricsMiddleware {
	return &MetricsMiddleware{
		inner: inner,
	}
}

func (mm *MetricsMiddleware) Get(key string) (string, error) {
	start := time.Now()

	val, err := mm.inner.Get(key)
	latency := time.Since(start)

	mm.mu.Lock()
	mm.getCalls++
	mm.totalGetLatency += latency

	if err != nil {
		mm.getMisses++
	}
	mm.mu.Unlock()

	return val, err
}

func (mm *MetricsMiddleware) Set(key, value string) error {
	start := time.Now()

	err := mm.inner.Set(key, value)
	latency := time.Since(start)

	mm.mu.Lock()
	mm.setCalls++
	mm.totalSetLatency += latency
	mm.mu.Unlock()

	return err
}

func (mm *MetricsMiddleware) Delete(key string) error {
	err := mm.inner.Delete(key)

	mm.mu.Lock()
	mm.deleteCalls++
	mm.mu.Unlock()

	return err
}

func (mm *MetricsMiddleware) Len() int {
	length := mm.inner.Len()

	mm.mu.Lock()
	mm.lenCalls++
	mm.mu.Unlock()

	return length
}

func (mm *MetricsMiddleware) Keys() []string {
	keys := mm.inner.Keys()

	mm.mu.Lock()
	mm.keyCalls++
	mm.mu.Unlock()

	return keys
}

func (mm *MetricsMiddleware) Incr(key string) (int, error) {
	start := time.Now()

	value, err := mm.inner.Incr(key)
	latency := time.Since(start)

	mm.mu.Lock()
	mm.incrCalls++
	mm.totalIncrLatency += latency
	mm.mu.Unlock()

	return value, err
}

func (mm *MetricsMiddleware) Decr(key string) (int, error) {
	start := time.Now()

	value, err := mm.inner.Decr(key)
	latency := time.Since(start)

	mm.mu.Lock()
	mm.decrCalls++
	mm.totalDecrLatency += latency
	mm.mu.Unlock()

	return value, err
}

func (mm *MetricsMiddleware) Report() {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	fmt.Printf("GET calls: %d, (missed %d)\n", mm.getCalls, mm.getMisses)
	fmt.Printf("SET calls: %d, (latency: %s)\n", mm.setCalls, mm.totalSetLatency)
	fmt.Printf("DELETE calls: %d\n", mm.deleteCalls)
	fmt.Printf("LEN calls: %d\n", mm.lenCalls)
	fmt.Printf("KEYS calls: %d\n", mm.keyCalls)
	fmt.Printf("INCR calls: %d\n", mm.incrCalls)
	fmt.Printf("DECR calls: %d\n", mm.decrCalls)

	if mm.getCalls > 0 {
		avg := mm.totalGetLatency / time.Duration(mm.getCalls)
		fmt.Printf("AVG GET latency: %s\n", avg)
	}

	if mm.setCalls > 0 {
		avg := mm.totalSetLatency / time.Duration(mm.setCalls)
		fmt.Printf("AVG SET latency: %s\n", avg)
	}

	if mm.incrCalls > 0 {
		avg := mm.totalIncrLatency / time.Duration(mm.incrCalls)
		fmt.Printf("AVG INCR latency: %s\n", avg)
	}

	if mm.decrCalls > 0 {
		avg := mm.totalDecrLatency / time.Duration(mm.decrCalls)
		fmt.Printf("AVG DECR latency: %s\n", avg)
	}
}
