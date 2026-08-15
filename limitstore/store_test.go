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

func TestUpdateAtCapacity(t *testing.T) {
	s := New(2)
	if err := s.Put("theme", "dark"); err != nil {
		t.Fatal(err)
	}
	if err := s.Put("density", "comfortable"); err != nil {
		t.Fatal(err)
	}
	if err := s.Put("theme", "light"); err != nil {
		t.Fatalf("update existing key: %v", err)
	}
	if got, _ := s.Get("theme"); got != "light" {
		t.Fatalf("theme = %q", got)
	}
	if s.Len() != 2 {
		t.Fatalf("Len() = %d", s.Len())
	}
	if err := s.Put("language", "en"); !errors.Is(err, ErrCapacity) {
		t.Fatalf("new key error = %v", err)
	}
	if _, ok := s.Get("language"); ok {
		t.Fatal("full store accepted a new key")
	}
}

func TestUpdateAfterDeleteAndRefill(t *testing.T) {
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
	if err := s.Put("layout", "wide"); err != nil {
		t.Fatalf("update refilled store: %v", err)
	}
	if got, _ := s.Get("layout"); got != "wide" {
		t.Fatalf("layout = %q", got)
	}
}

func TestZeroCapacityRejectsNewKey(t *testing.T) {
	s := New(0)
	if err := s.Put("theme", "dark"); !errors.Is(err, ErrCapacity) {
		t.Fatalf("Put error = %v", err)
	}
}
