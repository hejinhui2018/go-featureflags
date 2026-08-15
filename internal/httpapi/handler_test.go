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

func TestPutMissingValue(t *testing.T) {
	h := NewHandler(store.NewMemory())

	// Pre-populate a record so we can verify it is NOT overwritten.
	put1 := httptest.NewRecorder()
	h.ServeHTTP(put1, httptest.NewRequest(http.MethodPut, "/v1/records/theme", strings.NewReader(`{"value":"dark"}`)))
	if put1.Code != http.StatusOK {
		t.Fatalf("seed PUT status = %d, want 200", put1.Code)
	}

	// PUT with {} (missing value field) on an existing key must 400 and must not overwrite.
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/v1/records/theme", strings.NewReader(`{}`)))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("PUT {} status = %d, want 400", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"invalid value"`) {
		t.Fatalf("error body = %s, want it to contain \"invalid value\"", rec.Body.String())
	}

	get := httptest.NewRecorder()
	h.ServeHTTP(get, httptest.NewRequest(http.MethodGet, "/v1/records/theme", nil))
	if get.Code != http.StatusOK || !strings.Contains(get.Body.String(), `"value":"dark"`) {
		t.Fatalf("GET after failed overwrite = %d %s, want 200 with value dark", get.Code, get.Body.String())
	}

	// PUT with {} on a brand-new key must 400 and must not create a record.
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, httptest.NewRequest(http.MethodPut, "/v1/records/newkey", strings.NewReader(`{}`)))
	if rec2.Code != http.StatusBadRequest {
		t.Fatalf("PUT {} new key status = %d, want 400", rec2.Code)
	}

	get2 := httptest.NewRecorder()
	h.ServeHTTP(get2, httptest.NewRequest(http.MethodGet, "/v1/records/newkey", nil))
	if get2.Code != http.StatusNotFound {
		t.Fatalf("GET after failed create = %d, want 404", get2.Code)
	}
}
