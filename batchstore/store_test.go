package batchstore

import (
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
