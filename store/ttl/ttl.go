package ttl

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	ErrKeyExpired  = errors.New("key has expired")
	ErrKeyNotFound = errors.New("key not found")
	ErrEmptyKey    = errors.New("key is empty")
)

type TTLStore struct {
	mu   sync.RWMutex
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

	s.mu.Lock()
	defer s.mu.Unlock()

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

	s.mu.Lock()
	defer s.mu.Unlock()

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
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}

func (s *TTLStore) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()

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

func (s *TTLStore) Incr(key string) (int, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	current := 0
	expiresAt := now.Add(time.Hour)

	if entry, ok := s.data[key]; ok {
		if now.After(entry.expiresAt) {
			delete(s.data, key)
			return 0, ErrKeyExpired
		}

		value, err := strconv.Atoi(entry.value)
		if err != nil {
			return 0, fmt.Errorf("INCR %q: %w", key, err)
		}

		current = value
		expiresAt = entry.expiresAt
	}

	current++

	s.data[key] = ttlEntry{
		value:     strconv.Itoa(current),
		expiresAt: expiresAt,
	}

	return current, nil
}

func (s *TTLStore) Decr(key string) (int, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	current := 0
	expiresAt := now.Add(time.Hour)

	if entry, ok := s.data[key]; ok {
		if now.After(entry.expiresAt) {
			delete(s.data, key)
			return 0, ErrKeyExpired
		}

		value, err := strconv.Atoi(entry.value)
		if err != nil {
			return 0, fmt.Errorf("DECR %q: %w", key, err)
		}

		current = value
		expiresAt = entry.expiresAt
	}

	if current <= 0 {
		return 0, nil
	}

	current--

	s.data[key] = ttlEntry{
		value:     strconv.Itoa(current),
		expiresAt: expiresAt,
	}

	return current, nil
}
