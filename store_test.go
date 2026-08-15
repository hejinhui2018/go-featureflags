package leasecache

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandlerReturnsNotFoundAtExactExpiry(t *testing.T) {
	current := time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)
	store := NewStore(func() time.Time { return current })
	if err := store.Put("session-1", []byte("active"), 5*time.Minute); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	current = current.Add(5 * time.Minute)
	recorder := httptest.NewRecorder()
	Handler(store).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/leases/session-1", nil))
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body = %q", recorder.Code, http.StatusNotFound, recorder.Body.String())
	}
}

func TestStoreReturnsValueBeforeExpiry(t *testing.T) {
	current := time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)
	store := NewStore(func() time.Time { return current })
	if err := store.Put("session-1", []byte("active"), 5*time.Minute); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	current = current.Add(5*time.Minute - time.Nanosecond)
	value, ok := store.Get("session-1")
	if !ok || string(value) != "active" {
		t.Fatalf("Get() = %q, %v, want active, true", value, ok)
	}
}

func TestStoreCopiesInputAndOutput(t *testing.T) {
	store := NewStore(func() time.Time { return time.Unix(0, 0) })
	input := []byte("active")
	if err := store.Put("session-1", input, time.Minute); err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	input[0] = 'X'

	first, ok := store.Get("session-1")
	if !ok || string(first) != "active" {
		t.Fatalf("first Get() = %q, %v", first, ok)
	}
	first[0] = 'Y'
	second, _ := store.Get("session-1")
	if string(second) != "active" {
		t.Fatalf("second Get() = %q, want active", second)
	}
}

func TestPutValidatesArguments(t *testing.T) {
	store := NewStore(nil)
	if err := store.Put("", nil, time.Minute); !errors.Is(err, ErrEmptyKey) {
		t.Fatalf("empty key error = %v, want %v", err, ErrEmptyKey)
	}
	if err := store.Put("session-1", nil, 0); !errors.Is(err, ErrInvalidTTL) {
		t.Fatalf("zero ttl error = %v, want %v", err, ErrInvalidTTL)
	}
}
