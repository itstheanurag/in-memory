package ttl

import (
	"errors"
	"time"
)

var (
	ErrKeyExpired  = errors.New("key has expired")
	ErrKeyNotFound = errors.New("key not found")
	ErrEmptyKey    = errors.New("key is empty")
)

type TTLStore struct {
	data map[string]ttlEntry
}

type ttlEntry struct {
	value     string
	expiresAt time.Time
}

func NewTTLStore() *TTLStore {
	return &TTLStore{
		data: make(map[string]ttlEntry),
	}
}

func (s *TTLStore) Set(key, value string, ttl time.Duration) error {
	if key == "" {
		return ErrEmptyKey
	}

	s.data[key] = ttlEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}

	return nil
}

func (s *TTLStore) Get(key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}

	entry, ok := s.data[key]

	if !ok {
		return "", ErrKeyNotFound
	}

	if time.Now().After(entry.expiresAt) {
		delete(s.data, key)
		return "", ErrKeyExpired
	}

	return entry.value, nil
}

func (s *TTLStore) Delete(key string) {
	delete(s.data, key)
}

func (s *TTLStore) Len() int {
	now := time.Now()
	count := 0

	for key, entry := range s.data {
		if now.After(entry.expiresAt) {
			delete(s.data, key)
			continue
		}

		count++
	}

	return count
}
