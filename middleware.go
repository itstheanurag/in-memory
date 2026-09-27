package main

import (
	"log"
	"os"

	"github.com/itstheanurag/in-memory/store"
)

type LogginMiddleWareStore struct {
	inner  store.Storer
	logger *log.Logger
}

func NewLoggingMiddleware(inner store.Storer) *LogginMiddleWareStore {
	return &LogginMiddleWareStore{
		inner:  inner,
		logger: log.New(os.Stdout, "[log] ", log.Ltime|log.Lmicroseconds),
	}
}

func (m *LogginMiddleWareStore) Set(Key, value string) error {

	return nil
}

func (m *LogginMiddleWareStore) Get(Key string) (string, error) {

}

func (m *LogginMiddleWareStore) Delete(Key string) error {

	return nil
}

func (m *LogginMiddleWareStore) Len() int {
	return m.inner.Len()
}

func (m *LogginMiddleWareStore) Keys() []string {
	return m.inner.Keys()
}
