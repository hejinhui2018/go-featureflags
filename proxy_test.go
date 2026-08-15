package httpproxyheaders

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPrepareBackendRequestRemovesConnectionNominatedHeaders(t *testing.T) {
	in := httptest.NewRequest(http.MethodGet, "http://proxy.local/data", nil)
	in.Header.Set("Connection", "keep-alive, X-Internal-Hop, x-debug-hop")
	in.Header.Set("Keep-Alive", "timeout=5")
	in.Header.Set("X-Internal-Hop", "internal-value")
	in.Header.Set("X-Debug-Hop", "debug-value")

	out := PrepareBackendRequest(in)
	for _, name := range []string{"Connection", "Keep-Alive", "X-Internal-Hop", "X-Debug-Hop"} {
		if got := out.Header.Get(name); got != "" {
			t.Errorf("forwarded header %s = %q, want absent", name, got)
		}
	}
	if got := in.Header.Get("X-Internal-Hop"); got != "internal-value" {
		t.Fatalf("input request was modified: X-Internal-Hop = %q", got)
	}
}

func TestPrepareBackendRequestPreservesEndToEndHeaders(t *testing.T) {
	in := httptest.NewRequest(http.MethodPost, "http://proxy.local/items", nil)
	in.Header.Set("Accept", "application/json")
	in.Header.Set("Content-Type", "application/json")
	in.Header.Set("X-Request-ID", "request-123")

	out := PrepareBackendRequest(in)
	for name, want := range map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
		"X-Request-ID": "request-123",
	} {
		if got := out.Header.Get(name); got != want {
			t.Errorf("header %s = %q, want %q", name, got, want)
		}
	}
	if out.RequestURI != "" {
		t.Fatalf("RequestURI = %q, want empty", out.RequestURI)
	}
}
