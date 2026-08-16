package ttlcache

import (
	"testing"
	"time"
)

func TestPositiveTTLStoresUntilDeadline(t *testing.T) {
	now := time.Date(2026, 8, 16, 12, 0, 0, 0, time.UTC)
	cache := New(func() time.Time { return now })

	cache.Set("session", "open", 10*time.Second)
	value, ok := cache.Get("session")
	if !ok || value != "open" {
		t.Fatalf("positive ttl lookup = %q, %v; want open, true", value, ok)
	}

	now = now.Add(11 * time.Second)
	value, ok = cache.Get("session")
	if ok || value != "" {
		t.Fatalf("expired positive ttl lookup = %q, %v; want empty, false", value, ok)
	}
}

func TestZeroTTLIsNotStored(t *testing.T) {
	now := time.Date(2026, 8, 16, 12, 0, 0, 0, time.UTC)
	cache := New(func() time.Time { return now })

	cache.Set("boundary", "stale", 0)
	value, ok := cache.Get("boundary")
	if ok || value != "" {
		t.Fatalf("zero ttl lookup = %q, %v; want empty, false", value, ok)
	}
}
