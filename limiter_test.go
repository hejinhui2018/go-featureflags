package httpratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSameIPv4AddressSharesLimitAcrossPorts(t *testing.T) {
	handler := New(1).Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	first := httptest.NewRequest(http.MethodGet, "/items", nil)
	first.RemoteAddr = "203.0.113.10:41001"
	firstResponse := httptest.NewRecorder()
	handler.ServeHTTP(firstResponse, first)
	if firstResponse.Code != http.StatusNoContent {
		t.Fatalf("first status = %d, want %d", firstResponse.Code, http.StatusNoContent)
	}

	second := httptest.NewRequest(http.MethodGet, "/items", nil)
	second.RemoteAddr = "203.0.113.10:41002"
	secondResponse := httptest.NewRecorder()
	handler.ServeHTTP(secondResponse, second)
	if secondResponse.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d, want %d", secondResponse.Code, http.StatusTooManyRequests)
	}
}

func TestSameIPv6AddressSharesLimitAcrossPorts(t *testing.T) {
	handler := New(1).Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	for i, remote := range []string{"[2001:db8::20]:51001", "[2001:db8::20]:51002"} {
		request := httptest.NewRequest(http.MethodGet, "/items", nil)
		request.RemoteAddr = remote
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		want := http.StatusNoContent
		if i == 1 {
			want = http.StatusTooManyRequests
		}
		if response.Code != want {
			t.Fatalf("request %d status = %d, want %d", i+1, response.Code, want)
		}
	}
}

func TestDifferentAddressesHaveIndependentLimits(t *testing.T) {
	handler := New(1).Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	for _, remote := range []string{"203.0.113.10:41001", "203.0.113.11:41001"} {
		request := httptest.NewRequest(http.MethodGet, "/items", nil)
		request.RemoteAddr = remote
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusNoContent {
			t.Fatalf("remote %s status = %d, want %d", remote, response.Code, http.StatusNoContent)
		}
	}
}
