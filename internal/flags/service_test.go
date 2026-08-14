package flags

import (
	"context"
	"errors"
	"testing"
)

func TestServiceEnabledCachesValue(t *testing.T) {
	store := &countingStore{value: true}
	service := NewService(store)

	for i := 0; i < 2; i++ {
		enabled, err := service.Enabled(context.Background(), "tenant-a", "checkout")
		if err != nil {
			t.Fatalf("Enabled() error = %v", err)
		}
		if !enabled {
			t.Fatal("Enabled() = false, want true")
		}
	}
	if store.calls != 1 {
		t.Fatalf("store calls = %d, want 1", store.calls)
	}
}

func TestServiceEnabledValidatesInput(t *testing.T) {
	service := NewService(&countingStore{})
	_, err := service.Enabled(context.Background(), "", "checkout")
	if !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("Enabled() error = %v, want ErrInvalidArgument", err)
	}
}

type countingStore struct {
	value bool
	calls int
}

func (s *countingStore) Get(context.Context, string, string) (bool, error) {
	s.calls++
	return s.value, nil
}
