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

func NewMemory(initial map[string]string) *Memory {
	data := make(map[string]string, len(initial))
	for key, value := range initial {
		data[key] = value
	}
	return &Memory{data: data}
}

func (m *Memory) Get(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return "", ErrNotFound
	}

	m.mu.RLock()
	value, ok := m.data[key]
	m.mu.RUnlock()
	if !ok {
		return "", ErrNotFound
	}
	return value, nil
}
