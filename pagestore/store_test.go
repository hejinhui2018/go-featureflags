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

// TestListOffsetBeyondEnd is the regression test for the panic that occurred
// when offset was greater than the number of stored records. The store must
// return an empty page and a nil error instead of crashing.
func TestListOffsetBeyondEnd(t *testing.T) {
	s := New()
	s.Append("alpha")
	s.Append("beta")

	// offset far beyond the end - this used to panic.
	page, err := s.List(5, 2)
	if err != nil {
		t.Fatalf("expected nil error for out-of-range offset, got %v", err)
	}
	if len(page) != 0 {
		t.Fatalf("expected empty page, got %+v", page)
	}
}

// TestListOffsetEqualsRecordCount verifies that an offset exactly equal to
// the record count returns an empty page with nil error.
func TestListOffsetEqualsRecordCount(t *testing.T) {
	s := New()
	s.Append("alpha")
	s.Append("beta")

	page, err := s.List(2, 5)
	if err != nil {
		t.Fatalf("expected nil error when offset == record count, got %v", err)
	}
	if len(page) != 0 {
		t.Fatalf("expected empty page, got %+v", page)
	}
}

// TestListOffsetBeyondEndDoesNotAffectStorage verifies that an out-of-range
// List call has no side effects on the stored records.
func TestListOffsetBeyondEndDoesNotAffectStorage(t *testing.T) {
	s := New()
	s.Append("alpha")
	s.Append("beta")

	_, _ = s.List(100, 10)

	page, err := s.List(0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(page) != 2 {
		t.Fatalf("storage changed after out-of-range List: got %d records, want 2", len(page))
	}
	if page[0].Value != "alpha" || page[1].Value != "beta" {
		t.Fatalf("storage content changed: %+v", page)
	}
}
