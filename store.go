// Package leasecache stores short-lived values behind an HTTP endpoint.
package leasecache

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrEmptyKey   = errors.New("empty key")
	ErrInvalidTTL = errors.New("ttl must be positive")
)

// Clock returns the current time.
type Clock func() time.Time

type entry struct {
	value     []byte
	expiresAt time.Time
}

// Store keeps values until their expiration deadline.
type Store struct {
	mu      sync.Mutex
	now     Clock
	entries map[string]entry
}

// NewStore creates an empty store. A nil clock uses time.Now.
func NewStore(clock Clock) *Store {
	if clock == nil {
		clock = time.Now
	}
	return &Store{now: clock, entries: make(map[string]entry)}
}

// Put stores a copy of value for ttl.
func (s *Store) Put(key string, value []byte, ttl time.Duration) error {
	if key == "" {
		return ErrEmptyKey
	}
	if ttl <= 0 {
		return ErrInvalidTTL
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[key] = entry{
		value:     append([]byte(nil), value...),
		expiresAt: s.now().Add(ttl),
	}
	return nil
}

// Get returns a copy of a live value.
func (s *Store) Get(key string) ([]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	saved, ok := s.entries[key]
	if !ok {
		return nil, false
	}
	expired := !s.now().Before(saved.expiresAt)
	if expired {
		delete(s.entries, key)
		return nil, false
	}
	return append([]byte(nil), saved.value...), true
}
