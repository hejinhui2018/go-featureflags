package httpapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hejinhui2018/configstore-validation-004/internal/store"
)

func TestPutAndGetRecord(t *testing.T) {
	h := NewHandler(store.NewMemory())
	put := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/v1/records/theme", strings.NewReader(`{"value":"dark"}`))
	h.ServeHTTP(put, req)
	if put.Code != http.StatusOK {
		t.Fatalf("PUT status = %d, want 200", put.Code)
	}

	get := httptest.NewRecorder()
	h.ServeHTTP(get, httptest.NewRequest(http.MethodGet, "/v1/records/theme", nil))
	if get.Code != http.StatusOK || !strings.Contains(get.Body.String(), `"value":"dark"`) {
		t.Fatalf("GET response = %d %s", get.Code, get.Body.String())
	}
}

func TestPutMalformedJSON(t *testing.T) {
	h := NewHandler(store.NewMemory())
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/v1/records/theme", strings.NewReader("{")))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestGetMissingRecord(t *testing.T) {
	h := NewHandler(store.NewMemory())
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/records/missing", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	if _, err := io.ReadAll(rec.Body); err != nil {
		t.Fatal(err)
	}
}
