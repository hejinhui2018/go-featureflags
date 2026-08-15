// Package checkpointstore records the latest processed offset for each stream.
package checkpointstore

import (
	"errors"
	"sync"
)

var (
	ErrEmptyStream    = errors.New("stream is required")
	ErrNegativeOffset = errors.New("offset must not be negative")
	ErrStaleOffset    = errors.New("offset is older than the current checkpoint")
)

// Store keeps checkpoints in memory.
type Store struct {
	mu      sync.RWMutex
	offsets map[string]int64
}

// NewStore creates an empty checkpoint store.
func NewStore() *Store {
	return &Store{offsets: make(map[string]int64)}
}

// Advance records offset when it does not move the stream backwards.
// The returned boolean reports whether the stored offset changed.
func (s *Store) Advance(stream string, offset int64) (bool, error) {
	if stream == "" {
		return false, ErrEmptyStream
	}
	if offset < 0 {
		return false, ErrNegativeOffset
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	current, exists := s.offsets[stream]
	if exists && offset < current {
		return false, ErrStaleOffset
	}
	if exists && offset == current {
		return false, nil
	}
	s.offsets[stream] = offset
	return true, nil
}

// Current returns the stored offset for stream.
func (s *Store) Current(stream string) (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	offset, ok := s.offsets[stream]
	return offset, ok
}
