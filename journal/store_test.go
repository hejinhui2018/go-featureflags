package journal

import (
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
