package batchstore

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestImportRecords(t *testing.T) {
	s := New(map[string]string{"theme": "dark"})
	input := strings.NewReader("{\"key\":\"theme\",\"value\":\"light\"}\n{\"key\":\"layout\",\"value\":\"compact\"}\n")
	if err := s.Import(input); err != nil {
		t.Fatal(err)
	}
	if got, _ := s.Get("theme"); got != "light" {
		t.Fatalf("theme = %q", got)
	}
	if got, ok := s.Get("layout"); !ok || got != "compact" {
		t.Fatalf("layout = %q, %v", got, ok)
	}
}

func TestImportRejectsBlankKey(t *testing.T) {
	s := New(nil)
	err := s.Import(strings.NewReader("{\"key\":\"  \",\"value\":\"x\"}\n"))
	if err == nil || !strings.Contains(err.Error(), "blank key") {
		t.Fatalf("Import error = %v", err)
	}
	if _, ok := s.Get("  "); ok {
		t.Fatal("blank key was stored")
	}
}

func TestInvalidBatchLeavesStoreUnchanged(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "malformed JSON", input: "{\"key\":\"theme\",\"value\":\"light\"}\n{", want: "line 2"},
		{name: "blank key", input: "{\"key\":\"theme\",\"value\":\"light\"}\n{\"key\":\"  \",\"value\":\"x\"}\n", want: "blank key"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(map[string]string{"theme": "dark"})
			err := s.Import(strings.NewReader(tt.input))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Import error = %v", err)
			}
			if got, _ := s.Get("theme"); got != "dark" {
				t.Fatalf("theme after failed import = %q", got)
			}
			if _, ok := s.Get("layout"); ok {
				t.Fatal("failed batch created layout")
			}
		})
	}
}

func TestReaderFailureLeavesStoreUnchanged(t *testing.T) {
	readErr := errors.New("input interrupted")
	input := io.MultiReader(
		strings.NewReader("{\"key\":\"theme\",\"value\":\"light\"}\n"),
		errorReader{err: readErr},
	)
	s := New(map[string]string{"theme": "dark"})
	err := s.Import(input)
	if !errors.Is(err, readErr) {
		t.Fatalf("Import error = %v", err)
	}
	if got, _ := s.Get("theme"); got != "dark" {
		t.Fatalf("theme after reader failure = %q", got)
	}
}

type errorReader struct {
	err error
}

func (r errorReader) Read([]byte) (int, error) {
	return 0, r.err
}
