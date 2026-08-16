// Package appendledger stores contiguous, append-only records per stream.
package appendledger

import (
	"errors"
	"sync"
)

var (
	ErrEmptyStream    = errors.New("stream is required")
	ErrInvalidSeq     = errors.New("sequence must be positive")
	ErrEmptyValue     = errors.New("value is required")
	ErrSequenceGap    = errors.New("sequence gap")
	ErrRecordConflict = errors.New("record conflicts with existing value")
)

// Record is one stream entry.
type Record struct {
	Sequence int64  `json:"sequence"`
	Value    string `json:"value"`
}

type streamState struct {
	next    int64
	records map[int64]string
}

// Store keeps append-only records grouped by stream.
type Store struct {
	mu      sync.RWMutex
	streams map[string]streamState
}

// NewStore creates an empty append ledger.
func NewStore() *Store {
	return &Store{streams: make(map[string]streamState)}
}

// Append adds the next contiguous record. Repeating the same record is idempotent.
// The returned boolean reports whether a new record was written.
func (s *Store) Append(stream string, sequence int64, value string) (bool, error) {
	if stream == "" {
		return false, ErrEmptyStream
	}
	if sequence <= 0 {
		return false, ErrInvalidSeq
	}
	if value == "" {
		return false, ErrEmptyValue
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	state, exists := s.streams[stream]
	if !exists {
		state = streamState{next: 1, records: make(map[int64]string)}
	}
	if sequence < state.next {
		saved, ok := state.records[sequence]
		if !ok {
			return false, ErrSequenceGap
		}
		if saved != value {
			return false, ErrRecordConflict
		}
		s.streams[stream] = state
		return false, nil
	}
	if sequence > state.next {
		return false, ErrSequenceGap
	}
	state.records[sequence] = value
	state.next++
	s.streams[stream] = state
	return true, nil
}

// NextSequence returns the next sequence expected for stream.
func (s *Store) NextSequence(stream string) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state, ok := s.streams[stream]
	if !ok {
		return 1
	}
	return state.next
}
