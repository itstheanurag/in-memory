package middleware

import "github.com/itstheanurag/in-memory/store"

type MetricsMiddleware struct {
	inner store.Storer

	getCalls    int
	setCalls    int
	deleteCalls int
	lenCalls    int
	keyCalls    int
}

func NewMetricsMiddleware(inner store.Storer) *MetricsMiddleware {
	return &MetricsMiddleware{inner: inner}
}

func (mm *MetricsMiddleware) Get(key string) (string, error) {
	mm.getCalls++
	val, err := mm.inner.Get(key)
	return val, err
}

func (mm *MetricsMiddleware) Set(key, value string) error {
	mm.setCalls++
	return mm.inner.Set(key, value)
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
