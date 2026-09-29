package kv

import (
	"fmt"
	"maps"
	"sort"
	"strconv"
	"sync"
)

// Sentinel errors.
var (
	ErrEmptyKey  = fmt.Errorf("key cannot be empty")
	ErrStoreFull = fmt.Errorf("store is full, upgrade your tier")
)

type Store struct {
	mu      sync.RWMutex
	data    map[string]string
	maxSize int
}

func NewStore(size int) *Store {
	return &Store{
		data:    make(map[string]string),
		maxSize: size, // 0 means unlimited.
	}
}

func (s *Store) Get(key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("key %s does not exist", key)
	}

	return value, nil
}

func (s *Store) Set(key, value string) error {
	if key == "" {
		return ErrEmptyKey
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Updating an existing key does not affect capacity.
	if _, exists := s.data[key]; exists {
		s.data[key] = value
		return nil
	}

	// maxSize == 0 means unlimited.
	if s.maxSize > 0 && len(s.data) >= s.maxSize {
		return ErrStoreFull
	}

	s.data[key] = value

	return nil
}

func (s *Store) Delete(key string) error {
	if key == "" {
		return ErrEmptyKey
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)

	return nil
}

func (s *Store) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.data))

	for key := range s.data {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}

// Clone returns an independent copy of the store.
func (s *Store) Clone() *Store {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cp := &Store{
		data:    make(map[string]string, len(s.data)),
		maxSize: s.maxSize,
	}

	maps.Copy(cp.data, s.data)

	return cp
}

func (s *Store) Incr(key string) (int, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	current := 0

	if raw, ok := s.data[key]; ok {
		val, err := strconv.Atoi(raw)
		if err != nil {
			return 0, fmt.Errorf("INCR %q: %w", key, err)
		}

		current = val
	}

	current++

	s.data[key] = strconv.Itoa(current)

	return current, nil
}

func (s *Store) Decr(key string) (int, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	current := 0

	if raw, ok := s.data[key]; ok {
		val, err := strconv.Atoi(raw)
		if err != nil {
			return 0, fmt.Errorf("DECR %q: %w", key, err)
		}

		current = val
	}

	if current <= 0 {
		return 0, nil
	}

	current--

	s.data[key] = strconv.Itoa(current)

	return current, nil
}
