package journal

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type record struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Store struct {
	mu     sync.RWMutex
	file   *os.File
	path   string
	values map[string]string
	write  func([]byte) (int, error)
}

func Open(path string) (*Store, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, err
	}
	s := &Store{file: f, path: path, values: make(map[string]string)}
	if err := s.load(); err != nil {
		_ = f.Close()
		return nil, err
	}
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		_ = f.Close()
		return nil, err
	}
	s.write = f.Write
	return s, nil
}

func (s *Store) load() error {
	if _, err := s.file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	scanner := bufio.NewScanner(s.file)
	for scanner.Scan() {
		var rec record
		if err := json.Unmarshal(scanner.Bytes(), &rec); err != nil {
			return fmt.Errorf("read journal: %w", err)
		}
		if strings.TrimSpace(rec.Key) == "" {
			return errors.New("read journal: blank key")
		}
		s.values[rec.Key] = rec.Value
	}
	return scanner.Err()
}

func (s *Store) Put(key, value string) error {
	if strings.TrimSpace(key) == "" {
		return errors.New("blank key")
	}
	b, err := json.Marshal(record{Key: key, Value: value})
	if err != nil {
		return err
	}
	b = append(b, '\n')

	s.mu.Lock()
	defer s.mu.Unlock()

	// Record the file position before writing so that on failure
	// we can roll back any partial data and keep the journal valid.
	pos, err := s.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("append journal: %w", err)
	}

	if _, err := s.write(b); err != nil {
		// A partial write may have left a truncated JSON line in
		// the file.  Truncate back to the pre-write position so
		// the journal stays readable on reopen.
		_ = s.file.Truncate(pos)
		_, _ = s.file.Seek(pos, io.SeekStart)
		return fmt.Errorf("append journal: %w", err)
	}
	s.values[key] = value
	return nil
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.values[key]
	return v, ok
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.file.Close()
}
