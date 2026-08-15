package httptimeout

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestTimeoutReturnsCompletedResponse(t *testing.T) {
	handler := Timeout(200*time.Millisecond, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-State", "complete")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("created"))
	}))

	record := httptest.NewRecorder()
	handler.ServeHTTP(record, httptest.NewRequest(http.MethodGet, "/", nil))
	if record.Code != http.StatusCreated || record.Body.String() != "created" || record.Header().Get("X-Request-State") != "complete" {
		t.Fatalf("unexpected completed response: code=%d body=%q header=%q", record.Code, record.Body.String(), record.Header().Get("X-Request-State"))
	}
}

func TestTimeoutDiscardsLateHandlerWrite(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	finished := make(chan struct{})
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-release
		w.Header().Set("X-Late-Write", "present")
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("late handler response"))
		close(finished)
	})
	handler := Timeout(20*time.Millisecond, next)
	record := httptest.NewRecorder()
	done := make(chan struct{})
	go func() {
		handler.ServeHTTP(record, httptest.NewRequest(http.MethodGet, "/slow", nil))
		close(done)
	}()
	<-started
	<-done
	close(release)
	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("slow handler did not finish")
	}

	if record.Code != http.StatusGatewayTimeout {
		t.Fatalf("expected gateway timeout, got %d", record.Code)
	}
	if got := record.Body.String(); got != "request timed out\n" {
		t.Fatalf("expected timeout body only, got %q", got)
	}
	if record.Header().Get("X-Late-Write") != "" {
		t.Fatalf("late handler header leaked into timeout response: %q", record.Header().Get("X-Late-Write"))
	}
	if strings.Contains(record.Body.String(), "late handler") {
		t.Fatal("late handler body leaked into timeout response")
	}
}
