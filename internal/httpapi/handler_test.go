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

func TestPutMissingValue(t *testing.T) {
	h := NewHandler(store.NewMemory())
	seed := httptest.NewRecorder()
	h.ServeHTTP(seed, httptest.NewRequest(http.MethodPut, "/v1/records/theme", strings.NewReader(`{"value":"dark"}`)))
	if seed.Code != http.StatusOK {
		t.Fatalf("seed status = %d, want 200", seed.Code)
	}

	missing := httptest.NewRecorder()
	h.ServeHTTP(missing, httptest.NewRequest(http.MethodPut, "/v1/records/theme", strings.NewReader(`{}`)))
	if missing.Code != http.StatusBadRequest || !strings.Contains(missing.Body.String(), `"invalid value"`) {
		t.Fatalf("missing value response = %d %s, want 400 invalid value", missing.Code, missing.Body.String())
	}

	get := httptest.NewRecorder()
	h.ServeHTTP(get, httptest.NewRequest(http.MethodGet, "/v1/records/theme", nil))
	if get.Code != http.StatusOK || !strings.Contains(get.Body.String(), `"value":"dark"`) {
		t.Fatalf("GET after rejected update = %d %s, want original value", get.Code, get.Body.String())
	}

	newRecord := httptest.NewRecorder()
	h.ServeHTTP(newRecord, httptest.NewRequest(http.MethodPut, "/v1/records/newkey", strings.NewReader(`{}`)))
	if newRecord.Code != http.StatusBadRequest {
		t.Fatalf("missing value for new key status = %d, want 400", newRecord.Code)
	}
	newGet := httptest.NewRecorder()
	h.ServeHTTP(newGet, httptest.NewRequest(http.MethodGet, "/v1/records/newkey", nil))
	if newGet.Code != http.StatusNotFound {
		t.Fatalf("GET after rejected create = %d, want 404", newGet.Code)
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
