package ttlstore

import (
	"errors"
	"testing"
	"time"
)

type fakeClock struct{ now time.Time }

func (c *fakeClock) Now() time.Time { return c.now }
func (c *fakeClock) Advance(d time.Duration) {
	c.now = c.now.Add(d)
}

func TestPermanentValue(t *testing.T) {
	clock := &fakeClock{now: time.Unix(1000, 0)}
	s := New(clock.Now)
	if err := s.Put("theme", "dark", 0); err != nil {
		t.Fatal(err)
	}
	clock.Advance(24 * time.Hour)
	got, err := s.Get("theme")
	if err != nil || got != "dark" {
		t.Fatalf("got value=%q err=%v", got, err)
	}
}

func TestValueExpires(t *testing.T) {
	clock := &fakeClock{now: time.Unix(1000, 0)}
	s := New(clock.Now)
	if err := s.Put("token", "abc", 10*time.Second); err != nil {
		t.Fatal(err)
	}
	clock.Advance(9 * time.Second)
	if got, err := s.Get("token"); err != nil || got != "abc" {
		t.Fatalf("before expiry got value=%q err=%v", got, err)
	}
	clock.Advance(time.Second)
	if _, err := s.Get("token"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("at expiry expected ErrNotFound, got %v", err)
	}
	if got := s.Len(); got != 0 {
		t.Fatalf("expected empty store, got len=%d", got)
	}
}

func TestPutUpdatesValueAndTTL(t *testing.T) {
	clock := &fakeClock{now: time.Unix(1000, 0)}
	s := New(clock.Now)
	if err := s.Put("theme", "dark", time.Second); err != nil {
		t.Fatal(err)
	}
	if err := s.Put("theme", "light", 0); err != nil {
		t.Fatal(err)
	}
	clock.Advance(time.Hour)
	got, err := s.Get("theme")
	if err != nil || got != "light" {
		t.Fatalf("got value=%q err=%v", got, err)
	}
}

func TestPutRejectsNegativeTTLWithoutChangingStore(t *testing.T) {
	clock := &fakeClock{now: time.Unix(1000, 0)}
	s := New(clock.Now)
	if err := s.Put("theme", "dark", 0); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		key   string
		value string
	}{
		{key: "new", value: "value"},
		{key: "theme", value: "light"},
	} {
		if err := s.Put(tc.key, tc.value, -time.Nanosecond); !errors.Is(err, ErrInvalidTTL) {
			t.Fatalf("Put(%q) expected ErrInvalidTTL, got %v", tc.key, err)
		}
	}
	if got := s.Len(); got != 1 {
		t.Fatalf("rejected writes changed len to %d", got)
	}
	got, err := s.Get("theme")
	if err != nil || got != "dark" {
		t.Fatalf("rejected update changed value: value=%q err=%v", got, err)
	}
	if _, err := s.Get("new"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("rejected insert created a record: %v", err)
	}
}

func TestEmptyKeyAndDelete(t *testing.T) {
	s := New(nil)
	if err := s.Put("", "value", time.Second); !errors.Is(err, ErrEmptyKey) {
		t.Fatalf("expected ErrEmptyKey, got %v", err)
	}
	if s.Delete("missing") {
		t.Fatal("delete reported a missing key")
	}
	if err := s.Put("key", "value", 0); err != nil {
		t.Fatal(err)
	}
	if !s.Delete("key") {
		t.Fatal("delete did not remove existing key")
	}
	if _, err := s.Get("key"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
