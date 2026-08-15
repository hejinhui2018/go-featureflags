package journal

import (
	"errors"
	"io"
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

func TestFailedAppendKeepsJournalReadable(t *testing.T) {
	tests := []struct {
		name      string
		writeErr  error
		wantError error
	}{
		{name: "write error", writeErr: errors.New("disk full")},
		{name: "short write", wantError: io.ErrShortWrite},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "records.log")
			s, err := Open(path)
			if err != nil {
				t.Fatal(err)
			}
			if err := s.Put("theme", "dark"); err != nil {
				t.Fatal(err)
			}
			s.write = func(data []byte) (int, error) {
				n := len(data) / 2
				if _, err := s.file.Write(data[:n]); err != nil {
					return n, err
				}
				return n, tt.writeErr
			}
			err = s.Put("layout", "compact")
			if err == nil {
				t.Fatal("Put unexpectedly succeeded")
			}
			if tt.writeErr != nil && !errors.Is(err, tt.writeErr) {
				t.Fatalf("Put error = %v, want wrapped %v", err, tt.writeErr)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("Put error = %v, want wrapped %v", err, tt.wantError)
			}
			if err := s.Close(); err != nil {
				t.Fatal(err)
			}

			reopened, err := Open(path)
			if err != nil {
				t.Fatalf("reopen after failed append: %v", err)
			}
			defer reopened.Close()
			if got, ok := reopened.Get("theme"); !ok || got != "dark" {
				t.Fatalf("Get(theme) = %q, %v", got, ok)
			}
			if _, ok := reopened.Get("layout"); ok {
				t.Fatal("failed value was persisted")
			}
		})
	}
}
