package tokenbucket

import (
	"errors"
	"testing"
)

func TestReserveConsumesCapacity(t *testing.T) {
	bucket := New(5)
	if err := bucket.Reserve(3); err != nil {
		t.Fatalf("reserve 3: %v", err)
	}
	if got := bucket.Available(); got != 2 {
		t.Fatalf("available after reserve = %d, want 2", got)
	}
	if err := bucket.Reserve(3); !errors.Is(err, ErrInsufficientTokens) {
		t.Fatalf("reserve beyond capacity error = %v, want ErrInsufficientTokens", err)
	}
	if got := bucket.Available(); got != 2 {
		t.Fatalf("available after failed reserve = %d, want 2", got)
	}
}

func TestReleaseDoesNotIncreaseAboveCapacity(t *testing.T) {
	bucket := New(4)
	if err := bucket.Reserve(3); err != nil {
		t.Fatalf("reserve 3: %v", err)
	}

	bucket.Release(5)
	if got := bucket.Available(); got != 4 {
		t.Fatalf("available after over-release = %d, want 4", got)
	}

	if err := bucket.Reserve(4); err != nil {
		t.Fatalf("reserve full capacity after over-release: %v", err)
	}
	if got := bucket.Available(); got != 0 {
		t.Fatalf("available after full reserve = %d, want 0", got)
	}
}
