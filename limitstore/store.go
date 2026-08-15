package limitstore

import (
	"errors"
	"strings"
	"sync"
)

var (
	ErrCapacity   = errors.New("store capacity reached")
	ErrInvalidKey = errors.New("invalid key")
)

type Store struct {
	mu       sync.RWMutex
	capacity int
	values   map[string]string
}

func New(capacity int) *Store {
	if capacity < 0 {
		capacity = 0
	}
	return &Store{capacity: capacity, values: make(map[string]string)}
}

func (s *Store) Put(key, value string) error {
	if strings.TrimSpace(key) == "" {
		return ErrInvalidKey
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.values[key]
	if !exists && len(s.values) >= s.capacity {
		return ErrCapacity
	}
	s.values[key] = value
	return nil
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.values[key]
	return value, ok
}

func (s *Store) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.values[key]; !ok {
		return false
	}
	delete(s.values, key)
	return true
}

func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.values)
}
