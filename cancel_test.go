package httpcancel

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMarkRequestPreservesCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
	cancel()
	rec := httptest.NewRecorder()
	observed := make(chan bool, 1)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			observed <- r.Context().Err() == context.Canceled && Marker(r.Context()) == "edge"
		case <-time.After(100 * time.Millisecond):
			observed <- false
		}
		w.WriteHeader(http.StatusNoContent)
	})
	MarkRequest(next).ServeHTTP(rec, req)
	if !<-observed {
		t.Fatal("request cancellation was not preserved")
	}
}

func TestMarkRequestPassesLiveRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if Marker(r.Context()) != "edge" {
			t.Fatalf("marker = %q", Marker(r.Context()))
		}
		w.WriteHeader(http.StatusAccepted)
	})
	MarkRequest(next).ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d", rec.Code)
	}
}
