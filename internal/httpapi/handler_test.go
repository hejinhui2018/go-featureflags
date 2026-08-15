package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerReturnsFlag(t *testing.T) {
	handler := NewHandler(stubService{enabled: true})
	request := httptest.NewRequest(http.MethodGet, "/v1/flags/checkout", nil)
	request.Header.Set("X-Tenant-ID", "tenant-a")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if !strings.Contains(response.Body.String(), `"enabled":true`) {
		t.Fatalf("body = %q, want enabled=true", response.Body.String())
	}
}

func TestHandlerTrimsWhitespaceInResponse(t *testing.T) {
	handler := NewHandler(stubService{enabled: true})
	request := httptest.NewRequest(http.MethodGet, "/v1/flags/%20checkout%20", nil)
	request.Header.Set("X-Tenant-ID", "  tenant-a  ")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}

	var body map[string]any
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if got, want := body["tenant"], "tenant-a"; got != want {
		t.Fatalf("tenant = %q, want %q", got, want)
	}
	if got, want := body["flag"], "checkout"; got != want {
		t.Fatalf("flag = %q, want %q", got, want)
	}
	if got := body["enabled"]; got != true {
		t.Fatalf("enabled = %v, want true", got)
	}
}

type stubService struct {
	enabled bool
}

func (s stubService) Enabled(context.Context, string, string) (bool, error) {
	return s.enabled, nil
}
