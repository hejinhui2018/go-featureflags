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

func TestServiceEnabledIsolatesTenants(t *testing.T) {
	store := NewMemoryStore(map[string]map[string]bool{
		"tenant-a": {"checkout": true},
		"tenant-b": {"checkout": false},
	})
	service := NewService(store)

	enabledA, err := service.Enabled(context.Background(), "tenant-a", "checkout")
	if err != nil {
		t.Fatalf("Enabled(tenant-a) error = %v", err)
	}
	if !enabledA {
		t.Fatal("Enabled(tenant-a, checkout) = false, want true")
	}

	enabledB, err := service.Enabled(context.Background(), "tenant-b", "checkout")
	if err != nil {
		t.Fatalf("Enabled(tenant-b) error = %v", err)
	}
	if enabledB {
		t.Fatal("Enabled(tenant-b, checkout) = true, want false")
	}

	// Re-query tenant-a to confirm it still returns the correct cached value.
	enabledA2, err := service.Enabled(context.Background(), "tenant-a", "checkout")
	if err != nil {
		t.Fatalf("Enabled(tenant-a) second call error = %v", err)
	}
	if !enabledA2 {
		t.Fatal("Enabled(tenant-a, checkout) second call = false, want true")
	}
}

func TestServiceEnabledCancelledContextReturnsError(t *testing.T) {
	store := &countingStore{value: true}
	service := NewService(store)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	enabled, err := service.Enabled(ctx, "tenant-a", "checkout")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Enabled() error = %v, want context.Canceled", err)
	}
	if enabled {
		t.Fatal("Enabled() = true, want false on cancelled context")
	}
	if store.calls != 0 {
		t.Fatalf("store calls = %d, want 0 on cancelled context", store.calls)
	}
}

func TestServiceEnabledCancelledContextDoesNotPolluteCache(t *testing.T) {
	store := &countingStore{value: true}
	service := NewService(store)

	enabled, err := service.Enabled(context.Background(), "tenant-a", "checkout")
	if err != nil || !enabled {
		t.Fatalf("Enabled() = (%v, %v), want (true, nil)", enabled, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	enabled, err = service.Enabled(ctx, "tenant-a", "checkout")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Enabled() error = %v, want context.Canceled", err)
	}
	if enabled {
		t.Fatal("Enabled() = true, want false on cancelled context")
	}

	enabled, err = service.Enabled(context.Background(), "tenant-a", "checkout")
	if err != nil || !enabled {
		t.Fatalf("Enabled() = (%v, %v), want (true, nil)", enabled, err)
	}
	if store.calls != 1 {
		t.Fatalf("store calls = %d, want 1", store.calls)
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
