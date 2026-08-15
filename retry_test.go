package httpretrybody

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestRetryResendsPostBody(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read request: %v", err)
		}
		if calls.Add(1) == 1 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		if string(body) != `{"message":"hello"}` {
			t.Errorf("retry body = %q, want original payload", body)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("accepted"))
	}))
	defer server.Close()

	oneShotBody := io.Reader(io.NopCloser(strings.NewReader(`{"message":"hello"}`)))
	req, err := http.NewRequest(http.MethodPost, server.URL, oneShotBody)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := DoWithRetry(server.Client(), req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusCreated)
	}
}

func TestSuccessfulRequestIsSentOnce(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL, strings.NewReader("data"))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := DoWithRetry(server.Client(), req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if got := calls.Load(); got != 1 {
		t.Fatalf("request count = %d, want 1", got)
	}
}
