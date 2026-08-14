package flags

import (
	"context"
	"errors"
	"sync"
)

var ErrNotFound = errors.New("flag not found")

type Store interface {
	Get(ctx context.Context, tenant, name string) (bool, error)
}

type MemoryStore struct {
	mu     sync.RWMutex
	values map[string]map[string]bool
}

func NewMemoryStore(values map[string]map[string]bool) *MemoryStore {
	cloned := make(map[string]map[string]bool, len(values))
	for tenant, tenantValues := range values {
		cloned[tenant] = make(map[string]bool, len(tenantValues))
		for name, enabled := range tenantValues {
			cloned[tenant][name] = enabled
		}
	}
	return &MemoryStore{values: cloned}
}

func (s *MemoryStore) Get(_ context.Context, tenant, name string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tenantValues, ok := s.values[tenant]
	if !ok {
		return false, ErrNotFound
	}
	enabled, ok := tenantValues[name]
	if !ok {
		return false, ErrNotFound
	}
	return enabled, nil
}
