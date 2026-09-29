package middleware

import (
	"log"
	"os"

	"github.com/itstheanurag/in-memory/store"
)

type LoggingMiddleware struct {
	inner  store.Storer
	logger *log.Logger
	config LogConfig
}

type LogConfig struct {
	Get    bool
	Set    bool
	Delete bool
	Len    bool
	Keys   bool
	Incr   bool
	Decr   bool
}

func NewLoggingMiddleware(inner store.Storer, config LogConfig) *LoggingMiddleware {
	return &LoggingMiddleware{
		inner:  inner,
		logger: log.New(os.Stdout, "[log] ", log.Ltime|log.Lmicroseconds),
		config: config,
	}
}

func (m *LoggingMiddleware) Set(key, value string) error {
	if m.config.Set {
		m.logger.Printf("Set: key=%s", key)
	}

	return m.inner.Set(key, value)
}

func (m *LoggingMiddleware) Get(key string) (string, error) {
	if m.config.Get {
		m.logger.Printf("Get: key=%s", key)
	}

	val, err := m.inner.Get(key)

	if m.config.Get {
		if err != nil {
			m.logger.Printf("Get: key=%q -> miss (%v)", key, err)
		} else {
			m.logger.Printf("Get: key=%q -> hit", key)
		}
	}

	return val, err
}

func (m *LoggingMiddleware) Delete(key string) error {
	if m.config.Delete {
		m.logger.Printf("Delete: key=%s", key)
	}

	return m.inner.Delete(key)
}

func (m *LoggingMiddleware) Len() int {
	got := m.inner.Len()

	if m.config.Len {
		m.logger.Printf("Len: %d", got)
	}

	return got
}

func (m *LoggingMiddleware) Keys() []string {
	got := m.inner.Keys()

	if m.config.Keys {
		m.logger.Printf("Keys: %v", got)
	}

	return got
}

func (m *LoggingMiddleware) Incr(key string) (int, error) {
	got, err := m.inner.Incr(key)

	if m.config.Incr {
		if err != nil {
			m.logger.Printf("Incr: key=%s -> error=%v", key, err)
		} else {
			m.logger.Printf("Incr: key=%s -> %d", key, got)
		}
	}

	return got, err
}

func (m *LoggingMiddleware) Decr(key string) (int, error) {
	got, err := m.inner.Decr(key)

	if m.config.Decr {
		if err != nil {
			m.logger.Printf("Decr: key=%s -> error=%v", key, err)
		} else {
			m.logger.Printf("Decr: key=%s -> %d", key, got)
		}
	}

	return got, err
}
