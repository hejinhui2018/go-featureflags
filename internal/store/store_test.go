package store

import (
	"context"
	"testing"
)

func TestMemoryPutGet(t *testing.T) {
	m := NewMemory()
	if err := m.Put(context.Background(), " theme ", "dark"); err != nil {
		t.Fatal(err)
	}
	value, err := m.Get(context.Background(), "theme")
	if err != nil || value != "dark" {
		t.Fatalf("Get() = %q, %v", value, err)
	}
}

func TestMemoryRejectsBlankKey(t *testing.T) {
	if err := NewMemory().Put(context.Background(), "  ", "value"); err != ErrInvalidKey {
		t.Fatalf("Put() error = %v, want %v", err, ErrInvalidKey)
	}
}
