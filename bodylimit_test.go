package httpbodyclose

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type recordingBody struct {
	io.Reader
	closed bool
}

func (b *recordingBody) Close() error {
	b.closed = true
	return nil
}

func TestLimitMiddlewarePassesBodyWithinLimit(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", &recordingBody{Reader: strings.NewReader("hello")})
	rec := httptest.NewRecorder()
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		data, _ := io.ReadAll(r.Body)
		if string(data) != "hello" {
			t.Fatalf("body = %q", data)
		}
		w.WriteHeader(http.StatusCreated)
	})
	LimitMiddleware{Limit: 16, Next: next}.ServeHTTP(rec, req)
	if !called || rec.Code != http.StatusCreated {
		t.Fatalf("called=%v status=%d", called, rec.Code)
	}
}

func TestLimitMiddlewareClosesBodyWhenTooLarge(t *testing.T) {
	body := &recordingBody{Reader: strings.NewReader("0123456789")}
	req := httptest.NewRequest(http.MethodPost, "/", body)
	rec := httptest.NewRecorder()
	called := false
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called = true })
	LimitMiddleware{Limit: 4, Next: next}.ServeHTTP(rec, req)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", rec.Code)
	}
	if called {
		t.Fatal("next handler called for oversized body")
	}
	if !body.closed {
		t.Fatal("request body was not closed")
	}
}
