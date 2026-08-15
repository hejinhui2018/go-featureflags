package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hejinhui2018/recordstore-validation-003/internal/store"
)

func TestGetRecord(t *testing.T) {
	handler := NewHandler(store.NewMemory(map[string]string{"alpha": "one"}))
	req := httptest.NewRequest(http.MethodGet, "/v1/records/alpha", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Body.String(); got != "{\"key\":\"alpha\",\"value\":\"one\"}\n" {
		t.Fatalf("body = %q, want found response", got)
	}
}

func TestGetMissingRecord(t *testing.T) {
	handler := NewHandler(store.NewMemory(nil))
	req := httptest.NewRequest(http.MethodGet, "/v1/records/missing", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	if got := rec.Body.String(); got != "{\"error\":\"not found\"}\n" {
		t.Fatalf("body = %q, want not-found response", got)
	}
}
