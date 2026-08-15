package revisionstore

import (
	"errors"
	"sync"
)

var (
	ErrEmptyKey        = errors.New("key must not be empty")
	ErrInvalidRevision = errors.New("revision must not be negative")
	ErrNotFound        = errors.New("record not found")
	ErrConflict        = errors.New("revision conflict")
)

type Record struct {
	Key      string
	Value    string
	Revision int64
}

type Store struct {
	mu      sync.RWMutex
	records map[string]Record
}

func New() *Store {
	return &Store{records: make(map[string]Record)}
}

// Put creates a record when expectedRevision is zero, or updates a record whose revision matches.
func (s *Store) Put(key, value string, expectedRevision int64) (Record, error) {
	if key == "" {
		return Record{}, ErrEmptyKey
	}
	if expectedRevision < 0 {
		return Record{}, ErrInvalidRevision
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	current, exists := s.records[key]
	if !exists {
		if expectedRevision != 0 {
			return Record{}, ErrNotFound
		}
		created := Record{Key: key, Value: value, Revision: 1}
		s.records[key] = created
		return created, nil
	}
	if expectedRevision == 0 {
		return Record{}, ErrConflict
	}
	if expectedRevision != 0 && expectedRevision != current.Revision {
		return Record{}, ErrConflict
	}
	updated := Record{Key: key, Value: value, Revision: current.Revision + 1}
	s.records[key] = updated
	return updated, nil
}

func (s *Store) Get(key string) (Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	record, ok := s.records[key]
	if !ok {
		return Record{}, ErrNotFound
	}
	return record, nil
}
