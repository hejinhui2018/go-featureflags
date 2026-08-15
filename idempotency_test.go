package idempotency

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestSequentialRequestsReuseResponse(t *testing.T) {
	var calls atomic.Int32
	handler := New().Wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := calls.Add(1)
		w.Header().Set("X-Execution", fmt.Sprint(n))
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(w, "result-%d", n)
	}))

	first := performRequest(handler, "order-1")
	second := performRequest(handler, "order-1")

	if calls.Load() != 1 {
		t.Fatalf("business calls = %d, want 1", calls.Load())
	}
	assertSameResponse(t, first, second)
}

func TestDifferentKeysRemainIndependent(t *testing.T) {
	var calls atomic.Int32
	handler := New().Wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := calls.Add(1)
		w.Header().Set("X-Execution", fmt.Sprint(n))
		_, _ = fmt.Fprintf(w, "result-%d", n)
	}))

	first := performRequest(handler, "order-1")
	second := performRequest(handler, "order-2")

	if calls.Load() != 2 {
		t.Fatalf("business calls = %d, want 2", calls.Load())
	}
	if first.Body.String() == second.Body.String() {
		t.Fatalf("different keys returned the same body %q", first.Body.String())
	}
}

func TestConcurrentRequestsReuseInFlightResponse(t *testing.T) {
	var calls atomic.Int32
	entered := make(chan struct{}, 2)
	release := make(chan struct{})
	handler := New().Wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := calls.Add(1)
		entered <- struct{}{}
		<-release
		w.Header().Set("X-Execution", fmt.Sprint(n))
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, "result-%d", n)
	}))

	responses := make(chan *httptest.ResponseRecorder, 2)
	go func() { responses <- performRequest(handler, "order-1") }()
	waitForEntry(t, entered)
	go func() { responses <- performRequest(handler, "order-1") }()

	select {
	case <-entered:
	case <-time.After(200 * time.Millisecond):
	}
	close(release)
	first := <-responses
	second := <-responses

	if calls.Load() != 1 {
		t.Fatalf("business calls = %d, want 1", calls.Load())
	}
	assertSameResponse(t, first, second)
}

func TestConcurrentDifferentKeysDoNotBlockEachOther(t *testing.T) {
	entered := make(chan struct{}, 2)
	release := make(chan struct{})
	handler := New().Wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		entered <- struct{}{}
		<-release
		w.WriteHeader(http.StatusNoContent)
	}))

	done := make(chan struct{}, 2)
	go func() {
		performRequest(handler, "order-1")
		done <- struct{}{}
	}()
	go func() {
		performRequest(handler, "order-2")
		done <- struct{}{}
	}()

	waitForEntry(t, entered)
	waitForEntry(t, entered)
	close(release)
	<-done
	<-done
}

func waitForEntry(t *testing.T, entered <-chan struct{}) {
	t.Helper()
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("request did not reach the business handler")
	}
}

func performRequest(handler http.Handler, key string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodPost, "/orders", nil)
	request.Header.Set("Idempotency-Key", key)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func assertSameResponse(t *testing.T, first, second *httptest.ResponseRecorder) {
	t.Helper()
	if first.Code != second.Code {
		t.Fatalf("status codes differ: %d and %d", first.Code, second.Code)
	}
	if first.Header().Get("X-Execution") != second.Header().Get("X-Execution") {
		t.Fatalf("execution headers differ: %q and %q", first.Header().Get("X-Execution"), second.Header().Get("X-Execution"))
	}
	if first.Body.String() != second.Body.String() {
		t.Fatalf("response bodies differ: %q and %q", first.Body.String(), second.Body.String())
	}
}
