package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hejinhui2018/go-featureflags/internal/flags"
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

func TestHandlerTrimsTenantHeader(t *testing.T) {
	stub := &capturingStub{enabled: true}
	handler := NewHandler(stub)
	request := httptest.NewRequest(http.MethodGet, "/v1/flags/checkout", nil)
	request.Header.Set("X-Tenant-ID", "  tenant-a  ")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if stub.gotTenant != "tenant-a" {
		t.Fatalf("service received tenant = %q, want %q", stub.gotTenant, "tenant-a")
	}

	var body struct {
		Tenant string `json:"tenant"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body.Tenant != "tenant-a" {
		t.Fatalf("response tenant = %q, want %q", body.Tenant, "tenant-a")
	}
}

func TestHandlerBlankTenantAfterTrimReturns400(t *testing.T) {
	store := flags.NewMemoryStore(map[string]map[string]bool{
		"tenant-a": {"checkout": true},
	})
	handler := NewHandler(flags.NewService(store))
	request := httptest.NewRequest(http.MethodGet, "/v1/flags/checkout", nil)
	request.Header.Set("X-Tenant-ID", "   ")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

type capturingStub struct {
	enabled   bool
	gotTenant string
}

func (s *capturingStub) Enabled(_ context.Context, tenant, _ string) (bool, error) {
	s.gotTenant = tenant
	return s.enabled, nil
}
