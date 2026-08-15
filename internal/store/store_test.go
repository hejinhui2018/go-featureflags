package store

import (
	"context"
	"errors"
	"testing"
)

func TestMemoryGetFound(t *testing.T) {
	memory := NewMemory(map[string]string{"alpha": "one"})
	value, err := memory.Get(context.Background(), " alpha ")
	if err != nil || value != "one" {
		t.Fatalf("Get() = (%q, %v), want (%q, nil)", value, err, "one")
	}
}

func TestMemoryGetMissing(t *testing.T) {
	memory := NewMemory(nil)
	_, err := memory.Get(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get() error = %v, want ErrNotFound", err)
	}
}

func TestMemoryGetCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := NewMemory(map[string]string{"alpha": "one"}).Get(ctx, "alpha")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Get() error = %v, want context.Canceled", err)
	}
}
