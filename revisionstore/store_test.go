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
	if updated.Revision != 2 {
		t.Fatalf("expected revision 2 after update, got %d", updated.Revision)
	}
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

func TestZeroRevisionOnExistingKeyReturnsConflict(t *testing.T) {
	s := New()

	// Create the record normally.
	created, err := s.Put("theme", "dark", 0)
	if err != nil {
		t.Fatal(err)
	}
	if created.Revision != 1 || created.Value != "dark" {
		t.Fatalf("unexpected created record: %+v", created)
	}

	// Calling Put with expectedRevision=0 on an existing key must return ErrConflict.
	_, err = s.Put("theme", "light", 0)
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict for zero-revision on existing key, got %v", err)
	}

	// The original record must remain unchanged in both value and revision.
	got, err := s.Get("theme")
	if err != nil {
		t.Fatal(err)
	}
	if got.Value != "dark" {
		t.Fatalf("expected original value %q, got %q", "dark", got.Value)
	}
	if got.Revision != 1 {
		t.Fatalf("expected original revision %d, got %d", 1, got.Revision)
	}

	// After a successful update with matching revision, a subsequent call
	// with expectedRevision=0 must still return ErrConflict.
	updated, err := s.Put("theme", "light", created.Revision)
	if updated.Revision != 2 {
		t.Fatalf("expected revision 2 after update, got %d", updated.Revision)
	}
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.Put("theme", "amber", 0)
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict for zero-revision after update, got %v", err)
	}

	// Verify the record is still at the updated value and revision.
	got, err = s.Get("theme")
	if err != nil {
		t.Fatal(err)
	}
	if got.Value != "light" {
		t.Fatalf("expected value %q after conflict, got %q", "light", got.Value)
	}
	if got.Revision != 2 {
		t.Fatalf("expected revision %d after conflict, got %d", 2, got.Revision)
	}
}

func TestNewKeyWithZeroRevisionStillCreates(t *testing.T) {
	s := New()
	created, err := s.Put("brand-new", "first", 0)
	if err != nil {
		t.Fatalf("expected successful creation for new key, got %v", err)
	}
	if created.Revision != 1 || created.Value != "first" {
		t.Fatalf("unexpected created record: %+v", created)
	}

	// A second brand-new key should also be created with revision 1.
	created2, err := s.Put("another-new", "second", 0)
	if err != nil {
		t.Fatalf("expected successful creation for second new key, got %v", err)
	}
	if created2.Revision != 1 || created2.Value != "second" {
		t.Fatalf("unexpected created record: %+v", created2)
	}
}

func TestMultipleUpdatesAndConflict(t *testing.T) {
	s := New()
	// Create key "config" with revision 1.
	r1, err := s.Put("config", "v1", 0)
	if err != nil {
		t.Fatal(err)
	}
	// Update to revision 2.
	r2, err := s.Put("config", "v2", r1.Revision)
	if err != nil {
		t.Fatal(err)
	}
	if r2.Revision != 2 {
		t.Fatalf("expected revision 2, got %d", r2.Revision)
	}
	// Update to revision 3.
	r3, err := s.Put("config", "v3", r2.Revision)
	if err != nil {
		t.Fatal(err)
	}
	if r3.Revision != 3 {
		t.Fatalf("expected revision 3, got %d", r3.Revision)
	}
	// Stale update using revision 2 should conflict.
	if _, err := s.Put("config", "v-stale", r2.Revision); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict for stale revision, got %v", err)
	}
	// Record should remain at revision 3 with value "v3".
	got, err := s.Get("config")
	if err != nil {
		t.Fatal(err)
	}
	if got.Value != "v3" || got.Revision != 3 {
		t.Fatalf("unexpected record after conflict: %+v", got)
	}
}