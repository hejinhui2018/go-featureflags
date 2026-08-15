package httpredirect

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostRedirectPreservesMethodAndBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/submit" {
			RedirectToCanonical(w, r)
			return
		}
		body, _ := io.ReadAll(r.Body)
		w.Header().Set("X-Observed-Method", r.Method)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(body)
	}))
	defer server.Close()

	response, err := server.Client().Post(server.URL+"/submit", "text/plain", strings.NewReader("order=42"))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusCreated)
	}
	if got := response.Header.Get("X-Observed-Method"); got != http.MethodPost {
		t.Fatalf("final method = %q, want POST", got)
	}
	if string(body) != "order=42" {
		t.Fatalf("final body = %q, want original body", body)
	}
}

func TestGetRedirectStillReachesCanonicalEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/submit" {
			RedirectToCanonical(w, r)
			return
		}
		w.Header().Set("X-Observed-Method", r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	response, err := server.Client().Get(server.URL + "/submit")
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusOK)
	}
	if got := response.Header.Get("X-Observed-Method"); got != http.MethodGet {
		t.Fatalf("final method = %q, want GET", got)
	}
}
