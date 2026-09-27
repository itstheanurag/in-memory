package kv

import (
	"fmt"
	"sort"
)

// sential error in go
var ErrEmptyKey = fmt.Errorf("key cannot be empty")
var ErrStoreFull = fmt.Errorf("store is full, upgrade your tier")

type Store struct {
	data    map[string]string
	maxSize int
}

func NewStore(size int) *Store {
	return &Store{
		data:    make(map[string]string),
		maxSize: size, // let's make 0 size to be unlimited for the testing purposes.
	}
}

func (s *Store) Get(key string) (string, error) {

	if key == "" {
		return "", ErrEmptyKey
	}

	value, ok := s.data[key]

	if !ok {
		// return "", ErrKeyNotFound
		return "", fmt.Errorf("Key %s does not exists", key)
	}

	return value, nil
}

func (s *Store) Set(key, value string) error {
	if key == "" {
		return ErrEmptyKey
	}

	_, err := s.Get(key)

	if err == nil {
		s.data[key] = value
		return nil
	}

	if s.maxSize > 0 && s.Len() >= s.maxSize {
		return ErrStoreFull
	}

	s.data[key] = value
	return nil
}

func (s *Store) Delete(key string) error {
	if key == "" {
		return ErrEmptyKey
	}
	delete(s.data, key)
	return nil
}

func (s *Store) Keys() []string {
	keys := make([]string, 0, s.Len())

	for key := range s.data {
		keys = append(keys, key)
	}
	// sort the keys
	sort.Strings(keys)
	return keys
}

func (s *Store) Len() int {
	return len(s.data)
}
