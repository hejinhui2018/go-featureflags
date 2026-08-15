package cursorstore

import (
	"errors"
	"reflect"
	"testing"
)

func seeded() *Store {
	s := New()
	s.Put("alpha", "1")
	s.Put("beta", "2")
	s.Put("gamma", "3")
	return s
}

func TestListFirstPageAndNextCursor(t *testing.T) {
	got, next, err := seeded().List("", 2)
	if err != nil {
		t.Fatal(err)
	}
	want := []Entry{{Key: "alpha", Value: "1"}, {Key: "beta", Value: "2"}}
	if !reflect.DeepEqual(got, want) || next != "beta" {
		t.Fatalf("got entries=%v next=%q", got, next)
	}
}

func TestListAfterCursor(t *testing.T) {
	got, next, err := seeded().List("beta", 2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []Entry{{Key: "gamma", Value: "3"}}) || next != "" {
		t.Fatalf("got entries=%v next=%q", got, next)
	}
}

func TestListLimitValidation(t *testing.T) {
	_, _, err := seeded().List("", 0)
	if !errors.Is(err, ErrInvalidLimit) {
		t.Fatalf("expected ErrInvalidLimit, got %v", err)
	}
}

func TestPutUpdatesAndKeepsOrder(t *testing.T) {
	s := seeded()
	s.Put("beta", "updated")
	s.Put("aardvark", "0")
	got, _, err := s.List("", 10)
	if err != nil {
		t.Fatal(err)
	}
	want := []Entry{{Key: "aardvark", Value: "0"}, {Key: "alpha", Value: "1"}, {Key: "beta", Value: "updated"}, {Key: "gamma", Value: "3"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

// Regression: a non-empty cursor that does not match any stored key must
// return ErrInvalidCursor with no entries and no next-page cursor.
func TestListNonExistentCursorReturnsError(t *testing.T) {
	got, next, err := seeded().List("nonexistent", 2)
	if !errors.Is(err, ErrInvalidCursor) {
		t.Fatalf("expected ErrInvalidCursor, got %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil entries, got %v", got)
	}
	if next != "" {
		t.Fatalf("expected empty next cursor, got %q", next)
	}
}

// Regression: a non-empty cursor against an empty store must also yield
// ErrInvalidCursor rather than silently returning the first page.
func TestListNonExistentCursorOnEmptyStore(t *testing.T) {
	s := New()
	got, next, err := s.List("ghost", 5)
	if !errors.Is(err, ErrInvalidCursor) {
		t.Fatalf("expected ErrInvalidCursor, got %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil entries, got %v", got)
	}
	if next != "" {
		t.Fatalf("expected empty next cursor, got %q", next)
	}
}
