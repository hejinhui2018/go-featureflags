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
