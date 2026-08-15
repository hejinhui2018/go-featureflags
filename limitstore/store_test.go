package limitstore

import (
	"errors"
	"testing"
)

func TestPutAndGet(t *testing.T) {
	s := New(2)
	if err := s.Put("theme", "dark"); err != nil {
		t.Fatal(err)
	}
	if got, ok := s.Get("theme"); !ok || got != "dark" {
		t.Fatalf("Get(theme) = %q, %v", got, ok)
	}
}

func TestCapacityRejectsNewKey(t *testing.T) {
	s := New(1)
	if err := s.Put("theme", "dark"); err != nil {
		t.Fatal(err)
	}
	if err := s.Put("layout", "compact"); !errors.Is(err, ErrCapacity) {
		t.Fatalf("Put(layout) error = %v", err)
	}
}

func TestDeleteFreesCapacity(t *testing.T) {
	s := New(1)
	if err := s.Put("theme", "dark"); err != nil {
		t.Fatal(err)
	}
	if !s.Delete("theme") {
		t.Fatal("Delete(theme) = false")
	}
	if err := s.Put("layout", "compact"); err != nil {
		t.Fatal(err)
	}
}

func TestRejectsBlankKey(t *testing.T) {
	s := New(1)
	if err := s.Put("  ", "value"); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("Put(blank) error = %v", err)
	}
}
