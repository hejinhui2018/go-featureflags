package cursorstore

import (
	"errors"
	"sort"
)

var (
	ErrInvalidCursor = errors.New("invalid cursor")
	ErrInvalidLimit  = errors.New("invalid limit")
)

type Entry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Store struct {
	entries []Entry
}

func New() *Store { return &Store{} }

func (s *Store) Put(key, value string) {
	for i := range s.entries {
		if s.entries[i].Key == key {
			s.entries[i].Value = value
			return
		}
	}
	s.entries = append(s.entries, Entry{Key: key, Value: value})
	sort.Slice(s.entries, func(i, j int) bool { return s.entries[i].Key < s.entries[j].Key })
}

// List returns up to limit entries after the supplied key. An empty cursor starts at the first entry.
func (s *Store) List(after string, limit int) ([]Entry, string, error) {
	if limit <= 0 {
		return nil, "", ErrInvalidLimit
	}
	start := 0
	if after != "" {
		for i, entry := range s.entries {
			if entry.Key == after {
				start = i + 1
				break
			}
		}
	}
	if start >= len(s.entries) {
		return []Entry{}, "", nil
	}
	end := start + limit
	if end > len(s.entries) {
		end = len(s.entries)
	}
	result := append([]Entry(nil), s.entries[start:end]...)
	next := ""
	if end < len(s.entries) {
		next = s.entries[end-1].Key
	}
	return result, next, nil
}
