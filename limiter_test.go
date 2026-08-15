package httplimitercancel

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestCanceledRequestReleasesSlot(t *testing.T) {
	var handled atomic.Int32
	middleware := LimitConcurrent(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handled.Add(1)
		w.WriteHeader(http.StatusNoContent)
	}), 1)

	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	first := httptest.NewRequest(http.MethodGet, "/work", nil).WithContext(canceled)
	middleware.ServeHTTP(httptest.NewRecorder(), first)

	done := make(chan struct{})
	go func() {
		second := httptest.NewRequest(http.MethodGet, "/work", nil)
		middleware.ServeHTTP(httptest.NewRecorder(), second)
		close(done)
	}()

	select {
	case <-done:
		if got := handled.Load(); got != 1 {
			t.Fatalf("handled requests = %d, want 1", got)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("request remained blocked after the canceled request returned")
	}
}

func TestCompletedRequestReleasesSlot(t *testing.T) {
	var handled atomic.Int32
	middleware := LimitConcurrent(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handled.Add(1)
		w.WriteHeader(http.StatusNoContent)
	}), 1)

	for i := 0; i < 2; i++ {
		request := httptest.NewRequest(http.MethodGet, "/work", nil)
		response := httptest.NewRecorder()
		middleware.ServeHTTP(response, request)
		if response.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
		}
	}
	if got := handled.Load(); got != 2 {
		t.Fatalf("handled requests = %d, want 2", got)
	}
}
