package ttlstore

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrEmptyKey   = errors.New("key must not be empty")
	ErrInvalidTTL = errors.New("ttl must not be negative")
	ErrNotFound   = errors.New("key not found")
)

type entry struct {
	value     string
	expiresAt time.Time
}

type Store struct {
	mu      sync.Mutex
	entries map[string]entry
	now     func() time.Time
}

func New(now func() time.Time) *Store {
	if now == nil {
		now = time.Now
	}
	return &Store{entries: make(map[string]entry), now: now}
}

// Put stores a value. A zero TTL keeps the value until it is deleted.
// A negative TTL is rejected with ErrInvalidTTL without modifying the store.
func (s *Store) Put(key, value string, ttl time.Duration) error {
	if key == "" {
		return ErrEmptyKey
	}
	if ttl < 0 {
		return ErrInvalidTTL
	}
	expiresAt := time.Time{}
	if ttl != 0 {
		expiresAt = s.now().Add(ttl)
	}
	s.mu.Lock()
	s.entries[key] = entry{value: value, expiresAt: expiresAt}
	s.mu.Unlock()
	return nil
}

func (s *Store) Get(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.entries[key]
	if !ok {
		return "", ErrNotFound
	}
	if !item.expiresAt.IsZero() && !s.now().Before(item.expiresAt) {
		delete(s.entries, key)
		return "", ErrNotFound
	}
	return item.value, nil
}

func (s *Store) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[key]; !ok {
		return false
	}
	delete(s.entries, key)
	return true
}

func (s *Store) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	for key, item := range s.entries {
		if !item.expiresAt.IsZero() && !now.Before(item.expiresAt) {
			delete(s.entries, key)
		}
	}
	return len(s.entries)
}