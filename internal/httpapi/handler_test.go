package httpapi

import (
	"context"
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

type stubService struct {
	enabled bool
}

func (s stubService) Enabled(context.Context, string, string) (bool, error) {
	return s.enabled, nil
}
