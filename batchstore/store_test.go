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

func TestImportAtomicOnBadJSON(t *testing.T) {
	s := New(map[string]string{"theme": "dark"})
	input := strings.NewReader("{\"key\":\"theme\",\"value\":\"light\"}\n{bad json}\n")
	err := s.Import(input)
	if err == nil || !strings.Contains(err.Error(), "line 2") {
		t.Fatalf("Import error = %v", err)
	}
	if got, _ := s.Get("theme"); got != "dark" {
		t.Fatalf("theme should remain dark, got %q", got)
	}
}

func TestImportAtomicOnBlankKey(t *testing.T) {
	s := New(map[string]string{"theme": "dark"})
	input := strings.NewReader("{\"key\":\"theme\",\"value\":\"light\"}\n{\"key\":\"  \",\"value\":\"x\"}\n")
	err := s.Import(input)
	if err == nil || !strings.Contains(err.Error(), "line 2") || !strings.Contains(err.Error(), "blank key") {
		t.Fatalf("Import error = %v", err)
	}
	if got, _ := s.Get("theme"); got != "dark" {
		t.Fatalf("theme should remain dark, got %q", got)
	}
}

func TestImportAtomicNewRecordNotStored(t *testing.T) {
	s := New(map[string]string{"theme": "dark"})
	input := strings.NewReader("{\"key\":\"theme\",\"value\":\"light\"}\n{\"key\":\"layout\",\"value\":\"compact\"}\n{bad json}\n")
	err := s.Import(input)
	if err == nil || !strings.Contains(err.Error(), "line 3") {
		t.Fatalf("Import error = %v", err)
	}
	if got, _ := s.Get("theme"); got != "dark" {
		t.Fatalf("theme should remain dark, got %q", got)
	}
	if _, ok := s.Get("layout"); ok {
		t.Fatal("layout should not be stored")
	}
}

func TestImportAtomicErrorPreservesLineAndMessage(t *testing.T) {
	s := New(nil)
	input := strings.NewReader("{\"key\":\"a\",\"value\":\"1\"}\nnot json\n")
	err := s.Import(input)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "line 2") {
		t.Fatalf("error should contain line 2, got: %v", err)
	}
}

func TestImportAtomicMultipleValidLinesStillWork(t *testing.T) {
	s := New(map[string]string{"theme": "dark"})
	input := strings.NewReader("{\"key\":\"theme\",\"value\":\"light\"}\n{\"key\":\"lang\",\"value\":\"en\"}\n{\"key\":\"layout\",\"value\":\"compact\"}\n")
	if err := s.Import(input); err != nil {
		t.Fatal(err)
	}
	if got, _ := s.Get("theme"); got != "light" {
		t.Fatalf("theme = %q", got)
	}
	if got, _ := s.Get("lang"); got != "en" {
		t.Fatalf("lang = %q", got)
	}
	if got, _ := s.Get("layout"); got != "compact" {
		t.Fatalf("layout = %q", got)
	}
}
