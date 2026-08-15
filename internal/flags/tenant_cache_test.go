package flags

import (
	"context"
	"testing"
)

func TestServiceEnabledSeparatesTenantCaches(t *testing.T) {
	store := NewMemoryStore(map[string]map[string]bool{
		"tenant-a": {"checkout": true},
		"tenant-b": {"checkout": false},
	})
	service := NewService(store)

	enabled, err := service.Enabled(context.Background(), "tenant-a", "checkout")
	if err != nil {
		t.Fatalf("Enabled(tenant-a) error = %v", err)
	}
	if !enabled {
		t.Fatal("Enabled(tenant-a, checkout) = false, want true")
	}

	enabled, err = service.Enabled(context.Background(), "tenant-b", "checkout")
	if err != nil {
		t.Fatalf("Enabled(tenant-b) error = %v", err)
	}
	if enabled {
		t.Fatal("Enabled(tenant-b, checkout) = true, want false")
	}
}
