package revisionstore

import (
	"errors"
	"testing"
)

func TestCreateAndMatchingUpdate(t *testing.T) {
	s := New()
	created, err := s.Put("theme", "dark", 0)
	if err != nil {
		t.Fatal(err)
	}
	if created.Revision != 1 || created.Value != "dark" {
		t.Fatalf("unexpected created record: %+v", created)
	}
	updated, err := s.Put("theme", "light", created.Revision)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Revision != 2 || updated.Value != "light" {
		t.Fatalf("unexpected updated record: %+v", updated)
	}
}

func TestStaleRevisionReturnsConflict(t *testing.T) {
	s := New()
	created, err := s.Put("theme", "dark", 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Put("theme", "light", created.Revision); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Put("theme", "contrast", created.Revision); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
	got, err := s.Get("theme")
	if err != nil || got.Value != "light" || got.Revision != 2 {
		t.Fatalf("conflicting update changed record: record=%+v err=%v", got, err)
	}
}

func TestMissingRecordWithRevisionReturnsNotFound(t *testing.T) {
	s := New()
	if _, err := s.Put("missing", "value", 3); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestInputValidationAndGet(t *testing.T) {
	s := New()
	if _, err := s.Put("", "value", 0); !errors.Is(err, ErrEmptyKey) {
		t.Fatalf("expected ErrEmptyKey, got %v", err)
	}
	if _, err := s.Put("key", "value", -1); !errors.Is(err, ErrInvalidRevision) {
		t.Fatalf("expected ErrInvalidRevision, got %v", err)
	}
	if _, err := s.Get("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
