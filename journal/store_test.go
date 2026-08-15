package journal

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestPutGetAfterReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "records.log")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Put("theme", "dark"); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	s, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if got, ok := s.Get("theme"); !ok || got != "dark" {
		t.Fatalf("Get(theme) = %q, %v", got, ok)
	}
}

func TestRejectBlankKey(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "records.log"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if err := s.Put("  ", "value"); err == nil {
		t.Fatal("Put accepted a blank key")
	}
}

func TestPutFailureDoesNotCorruptJournal(t *testing.T) {
	path := filepath.Join(t.TempDir(), "records.log")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Put("theme", "dark"); err != nil {
		t.Fatal(err)
	}

	// Simulate a partial write (e.g. disk full): write only half
	// of the record bytes, then return an error.
	s.write = func(b []byte) (int, error) {
		n := len(b) / 2
		_, _ = s.file.Write(b[:n])
		return n, errors.New("disk full")
	}

	if err := s.Put("layout", "compact"); err == nil {
		t.Fatal("Put should have returned an error")
	}

	if err := s.Close(); err != nil {
		t.Fatal(err)
	}

	// Reopen: should succeed; theme preserved, layout absent.
	s, err = Open(path)
	if err != nil {
		t.Fatalf("Open after failed Put: %v", err)
	}
	defer s.Close()

	if got, ok := s.Get("theme"); !ok || got != "dark" {
		t.Fatalf(`Get(theme) = %q, %v; want "dark", true`, got, ok)
	}
	if _, ok := s.Get("layout"); ok {
		t.Fatal("layout should not exist after failed Put")
	}
}
