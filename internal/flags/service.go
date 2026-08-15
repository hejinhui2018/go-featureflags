package flags

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrInvalidArgument = errors.New("tenant and flag name are required")

type Service struct {
	store Store
	mu    sync.RWMutex
	cache map[string]bool
}

func NewService(store Store) *Service {
	return &Service{store: store, cache: make(map[string]bool)}
}

func (s *Service) Enabled(ctx context.Context, tenant, name string) (bool, error) {
	tenant = strings.TrimSpace(tenant)
	name = strings.TrimSpace(name)
	if tenant == "" || name == "" {
		return false, ErrInvalidArgument
	}

	if err := ctx.Err(); err != nil {
		return false, err
	}

	key := tenant + ":" + name
	s.mu.RLock()
	cached, ok := s.cache[key]
	s.mu.RUnlock()
	if ok {
		return cached, nil
	}

	enabled, err := s.store.Get(ctx, tenant, name)
	if err != nil {
		return false, err
	}

	s.mu.Lock()
	s.cache[key] = enabled
	s.mu.Unlock()
	return enabled, nil
}
