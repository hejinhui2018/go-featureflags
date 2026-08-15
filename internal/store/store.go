package store

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrInvalidKey = errors.New("invalid key")
	ErrNotFound   = errors.New("record not found")
)

type Memory struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewMemory() *Memory {
	return &Memory{data: make(map[string]string)}
}

func (m *Memory) Put(ctx context.Context, key, value string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return ErrInvalidKey
	}
	m.mu.Lock()
	m.data[key] = value
	m.mu.Unlock()
	return nil
}

func (m *Memory) Get(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return "", ErrInvalidKey
	}
	m.mu.RLock()
	value, ok := m.data[key]
	m.mu.RUnlock()
	if !ok {
		return "", ErrNotFound
	}
	return value, nil
}
