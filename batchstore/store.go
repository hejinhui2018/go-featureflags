package batchstore

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
)

type Record struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Store struct {
	mu     sync.RWMutex
	values map[string]string
}

func New(initial map[string]string) *Store {
	values := make(map[string]string, len(initial))
	for key, value := range initial {
		values[key] = value
	}
	return &Store{values: values}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.values[key]
	return value, ok
}

func (s *Store) Import(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	line := 0
	for scanner.Scan() {
		line++
		var rec Record
		if err := json.Unmarshal(scanner.Bytes(), &rec); err != nil {
			return fmt.Errorf("line %d: %w", line, err)
		}
		if strings.TrimSpace(rec.Key) == "" {
			return fmt.Errorf("line %d: blank key", line)
		}
		s.mu.Lock()
		s.values[rec.Key] = rec.Value
		s.mu.Unlock()
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read import: %w", err)
	}
	return nil
}
