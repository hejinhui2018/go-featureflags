// Package quotareservoir stores bounded reservations by name.
package quotareservoir

import (
	"errors"
	"math"
	"sync"
)

var (
	ErrInvalidCapacity  = errors.New("capacity must be positive")
	ErrEmptyBatch       = errors.New("reservation batch is empty")
	ErrEmptyName        = errors.New("reservation name is required")
	ErrInvalidUnits     = errors.New("reservation units must be positive")
	ErrDuplicateName    = errors.New("reservation name appears more than once")
	ErrCapacityExceeded = errors.New("reservation capacity exceeded")
)

// Reservation describes one named allocation in a batch.
type Reservation struct {
	Name  string `json:"name"`
	Units int64  `json:"units"`
}

// Store tracks reservations against one fixed capacity.
type Store struct {
	mu           sync.RWMutex
	capacity     int64
	used         int64
	reservations map[string]int64
}

// NewStore creates an empty store with capacity.
func NewStore(capacity int64) (*Store, error) {
	if capacity <= 0 {
		return nil, ErrInvalidCapacity
	}
	return &Store{capacity: capacity, reservations: make(map[string]int64)}, nil
}

// ReserveBatch validates and applies a batch atomically.
func (s *Store) ReserveBatch(batch []Reservation) error {
	if len(batch) == 0 {
		return ErrEmptyBatch
	}

	seen := make(map[string]struct{}, len(batch))
	var total int64
	for _, reservation := range batch {
		if reservation.Name == "" {
			return ErrEmptyName
		}
		if reservation.Units <= 0 {
			return ErrInvalidUnits
		}
		if _, exists := seen[reservation.Name]; exists {
			return ErrDuplicateName
		}
		seen[reservation.Name] = struct{}{}
		if reservation.Units > math.MaxInt64-total {
			return ErrCapacityExceeded
		}
		total += reservation.Units
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if total > s.capacity-s.used {
		return ErrCapacityExceeded
	}

	for _, reservation := range batch {
		s.reservations[reservation.Name] += reservation.Units
		s.used += reservation.Units
	}
	return nil
}

// Used returns the total units currently reserved.
func (s *Store) Used() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.used
}

// Reservation returns the units reserved under name.
func (s *Store) Reservation(name string) (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	units, ok := s.reservations[name]
	return units, ok
}
