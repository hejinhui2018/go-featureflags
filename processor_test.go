package idemqueue

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProcessorReplaysSamePayload(t *testing.T) {
	processor := NewProcessor()
	calls := 0
	execute := func(context.Context, []byte) (Result, error) {
		calls++
		return Result{Status: http.StatusCreated, Body: []byte("job-1")}, nil
	}

	first, err := processor.Process(context.Background(), "request-1", []byte("alpha"), execute)
	if err != nil {
		t.Fatalf("first Process() error = %v", err)
	}
	second, err := processor.Process(context.Background(), "request-1", []byte("alpha"), execute)
	if err != nil {
		t.Fatalf("second Process() error = %v", err)
	}
	if calls != 1 {
		t.Fatalf("executor calls = %d, want 1", calls)
	}
	if string(first.Body) != "job-1" || string(second.Body) != "job-1" {
		t.Fatalf("replayed bodies = %q and %q", first.Body, second.Body)
	}
}

func TestProcessorRejectsKeyReuseWithDifferentPayload(t *testing.T) {
	processor := NewProcessor()
	calls := 0
	execute := func(context.Context, []byte) (Result, error) {
		calls++
		return Result{Status: http.StatusCreated, Body: []byte("job-1")}, nil
	}

	if _, err := processor.Process(context.Background(), "request-1", []byte("alpha"), execute); err != nil {
		t.Fatalf("first Process() error = %v", err)
	}
	if _, err := processor.Process(context.Background(), "request-1", []byte("beta"), execute); !errors.Is(err, ErrKeyConflict) {
		t.Fatalf("second Process() error = %v, want %v", err, ErrKeyConflict)
	}
	if calls != 1 {
		t.Fatalf("executor calls = %d, want 1", calls)
	}
}

func TestHandlerReturnsConflictForChangedPayload(t *testing.T) {
	processor := NewProcessor()
	calls := 0
	execute := func(context.Context, []byte) (Result, error) {
		calls++
		return Result{Status: http.StatusCreated, Body: []byte("job-1")}, nil
	}
	handler := Handler(processor, execute)

	first := httptest.NewRequest(http.MethodPost, "/jobs", strings.NewReader("alpha"))
	first.Header.Set("Idempotency-Key", "request-1")
	firstRecorder := httptest.NewRecorder()
	handler.ServeHTTP(firstRecorder, first)
	if firstRecorder.Code != http.StatusCreated {
		t.Fatalf("first status = %d, want %d", firstRecorder.Code, http.StatusCreated)
	}

	second := httptest.NewRequest(http.MethodPost, "/jobs", strings.NewReader("beta"))
	second.Header.Set("Idempotency-Key", "request-1")
	secondRecorder := httptest.NewRecorder()
	handler.ServeHTTP(secondRecorder, second)
	if secondRecorder.Code != http.StatusConflict {
		t.Fatalf("second status = %d, want %d", secondRecorder.Code, http.StatusConflict)
	}
	if calls != 1 {
		t.Fatalf("executor calls = %d, want 1", calls)
	}
}

func TestHandlerRequiresIdempotencyKey(t *testing.T) {
	handler := Handler(NewProcessor(), func(context.Context, []byte) (Result, error) {
		return Result{Status: http.StatusCreated}, nil
	})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/jobs", strings.NewReader("alpha")))

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}
