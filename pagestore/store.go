package pagestore

import (
	"errors"
	"sync"
)

var (
	ErrInvalidOffset = errors.New("offset must not be negative")
	ErrInvalidLimit  = errors.New("limit must be positive")
)

type Record struct {
	ID    int64
	Value string
}

type Store struct {
	mu      sync.RWMutex
	records []Record
	nextID  int64
}

func New() *Store {
	return &Store{nextID: 1}
}

func (s *Store) Append(value string) Record {
	s.mu.Lock()
	defer s.mu.Unlock()
	record := Record{ID: s.nextID, Value: value}
	s.nextID++
	s.records = append(s.records, record)
	return record
}

// List returns up to limit records starting at offset.
func (s *Store) List(offset, limit int) ([]Record, error) {
	if offset < 0 {
		return nil, ErrInvalidOffset
	}
	if limit <= 0 {
		return nil, ErrInvalidLimit
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if offset >= len(s.records) {
		return make([]Record, 0), nil
	}
	end := offset + limit
	if end > len(s.records) {
		end = len(s.records)
	}
	return append([]Record(nil), s.records[offset:end]...), nil
}
