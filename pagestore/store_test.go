package pagestore

import (
	"errors"
	"testing"
)

func TestAppendAndListPages(t *testing.T) {
	s := New()
	s.Append("alpha")
	s.Append("beta")
	s.Append("gamma")

	page, err := s.List(1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(page) != 2 || page[0].ID != 2 || page[0].Value != "beta" || page[1].ID != 3 {
		t.Fatalf("unexpected page: %+v", page)
	}
}

func TestListClampsAtEnd(t *testing.T) {
	s := New()
	s.Append("alpha")
	s.Append("beta")

	page, err := s.List(1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(page) != 1 || page[0].Value != "beta" {
		t.Fatalf("unexpected page: %+v", page)
	}

	empty, err := s.List(2, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(empty) != 0 {
		t.Fatalf("expected empty page, got %+v", empty)
	}
}

func TestListValidatesArguments(t *testing.T) {
	s := New()
	if _, err := s.List(-1, 1); !errors.Is(err, ErrInvalidOffset) {
		t.Fatalf("expected ErrInvalidOffset, got %v", err)
	}
	if _, err := s.List(0, 0); !errors.Is(err, ErrInvalidLimit) {
		t.Fatalf("expected ErrInvalidLimit, got %v", err)
	}
}

func TestListReturnsCopy(t *testing.T) {
	s := New()
	s.Append("alpha")
	page, err := s.List(0, 1)
	if err != nil {
		t.Fatal(err)
	}
	page[0].Value = "changed"
	again, err := s.List(0, 1)
	if err != nil {
		t.Fatal(err)
	}
	if again[0].Value != "alpha" {
		t.Fatalf("caller changed stored value: %+v", again[0])
	}
}

func TestListOffsetBeyondEndReturnsEmpty(t *testing.T) {
	s := New()
	s.Append("alpha")
	s.Append("beta")

	page, err := s.List(5, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(page) != 0 {
		t.Fatalf("expected empty page, got %+v", page)
	}

	all, err := s.List(0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 || all[0].Value != "alpha" || all[1].Value != "beta" {
		t.Fatalf("out-of-range read changed records: %+v", all)
	}
}
