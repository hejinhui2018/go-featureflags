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

// TestListOffsetBeyondRecords ensures that an offset greater than the number
// of stored records returns an empty page with a nil error instead of
// panicking.
func TestListOffsetBeyondRecords(t *testing.T) {
	s := New()
	s.Append("alpha")
	s.Append("beta")

	page, err := s.List(5, 2)
	if err != nil {
		t.Fatalf("expected nil error for out-of-range offset, got %v", err)
	}
	if len(page) != 0 {
		t.Fatalf("expected empty page for out-of-range offset, got %+v", page)
	}

	// A much larger offset should also be safe.
	page2, err2 := s.List(1000, 10)
	if err2 != nil {
		t.Fatalf("expected nil error for very large offset, got %v", err2)
	}
	if len(page2) != 0 {
		t.Fatalf("expected empty page for very large offset, got %+v", page2)
	}
}

// TestListOffsetEqualsRecordCount ensures that an offset exactly equal to the
// record count returns an empty page with a nil error. This boundary case
// must continue to work after the fix.
func TestListOffsetEqualsRecordCount(t *testing.T) {
	s := New()
	s.Append("alpha")
	s.Append("beta")

	page, err := s.List(2, 5)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(page) != 0 {
		t.Fatalf("expected empty page, got %+v", page)
	}
}

// TestListNormalPaginationFromStart verifies normal pagination from the
// beginning of the store.
func TestListNormalPaginationFromStart(t *testing.T) {
	s := New()
	s.Append("a")
	s.Append("b")
	s.Append("c")
	s.Append("d")
	s.Append("e")

	page, err := s.List(0, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(page) != 3 {
		t.Fatalf("expected 3 records, got %d", len(page))
	}
	if page[0].Value != "a" || page[1].Value != "b" || page[2].Value != "c" {
		t.Fatalf("unexpected page: %+v", page)
	}
}

// TestListTailTruncation verifies that the last page returns only the
// remaining records.
func TestListTailTruncation(t *testing.T) {
	s := New()
	s.Append("a")
	s.Append("b")
	s.Append("c")

	page, err := s.List(2, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(page) != 1 {
		t.Fatalf("expected 1 record, got %d", len(page))
	}
	if page[0].Value != "c" {
		t.Fatalf("expected 'c', got %s", page[0].Value)
	}
}

// TestListNegativeOffset verifies that a negative offset returns
// ErrInvalidOffset.
func TestListNegativeOffset(t *testing.T) {
	s := New()
	s.Append("alpha")
	if _, err := s.List(-5, 2); !errors.Is(err, ErrInvalidOffset) {
		t.Fatalf("expected ErrInvalidOffset, got %v", err)
	}
}

// TestListNonPositiveLimit verifies that a non-positive limit returns
// ErrInvalidLimit.
func TestListNonPositiveLimit(t *testing.T) {
	s := New()
	s.Append("alpha")
	if _, err := s.List(0, -1); !errors.Is(err, ErrInvalidLimit) {
		t.Fatalf("expected ErrInvalidLimit, got %v", err)
	}
	if _, err := s.List(0, 0); !errors.Is(err, ErrInvalidLimit) {
		t.Fatalf("expected ErrInvalidLimit, got %v", err)
	}
}

// TestListDoesNotMutateStore verifies that calling List with various
// offsets does not modify the underlying store contents.
func TestListDoesNotMutateStore(t *testing.T) {
	s := New()
	s.Append("alpha")
	s.Append("beta")

	// Call List with an out-of-range offset.
	_, _ = s.List(5, 2)
	// Call List with a valid offset.
	_, _ = s.List(0, 2)

	page, err := s.List(0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(page) != 2 {
		t.Fatalf("store contents changed: expected 2 records, got %d", len(page))
	}
	if page[0].Value != "alpha" || page[1].Value != "beta" {
		t.Fatalf("store contents changed: %+v", page)
	}
}

// TestListEmptyStore verifies that List on an empty store returns an empty
// page without panicking for various offsets.
func TestListEmptyStore(t *testing.T) {
	s := New()

	page, err := s.List(0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(page) != 0 {
		t.Fatalf("expected empty page, got %+v", page)
	}

	page2, err2 := s.List(5, 2)
	if err2 != nil {
		t.Fatal(err2)
	}
	if len(page2) != 0 {
		t.Fatalf("expected empty page, got %+v", page2)
	}
}
