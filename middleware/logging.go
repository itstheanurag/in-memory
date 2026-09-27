package middleware

import (
	"log"
	"os"

	"github.com/itstheanurag/in-memory/store"
)

type LogginMiddleWare struct {
	inner  store.Storer
	logger *log.Logger
}

func NewLoggingMiddleware(inner store.Storer) *LogginMiddleWare {
	return &LogginMiddleWare{
		inner:  inner,
		logger: log.New(os.Stdout, "[log] ", log.Ltime|log.Lmicroseconds),
	}
}

func (m *LogginMiddleWare) Set(key, value string) error {
	m.logger.Printf("Set: key=%s", key)
	return m.inner.Set(key, value)
}

func (m *LogginMiddleWare) Get(key string) (string, error) {
	m.logger.Printf("Get: key=%s", key)
	val, err := m.inner.Get(key)

	if err != nil {
		m.logger.Printf("Get: key=%q -> miss (%v)", key, err)
	} else {
		m.logger.Printf("Get: key=%q -> hit", key)
	}

	return val, err
}

func (m *LogginMiddleWare) Delete(key string) error {
	m.logger.Printf("Delete: key=%s", key)
	return m.inner.Delete(key)
}

func (m *LogginMiddleWare) Len() int {
	got := m.inner.Len()
	m.logger.Printf("Len: %d", got)
	return got
}

func (m *LogginMiddleWare) Keys() []string {
	got := m.inner.Keys()
	m.logger.Printf("Keys: %v", got)
	return got
}
