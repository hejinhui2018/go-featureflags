package checkpointstore

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func postCheckpoint(t *testing.T, handler http.Handler, stream, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, "/checkpoints/"+stream, strings.NewReader(body))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func TestHandlerAcceptsRepeatedCheckpoint(t *testing.T) {
	store := NewStore()
	handler := Handler(store)

	first := postCheckpoint(t, handler, "orders", `{"offset":42}`)
	if first.Code != http.StatusCreated {
		t.Fatalf("first checkpoint status = %d, want %d", first.Code, http.StatusCreated)
	}

	repeated := postCheckpoint(t, handler, "orders", `{"offset":42}`)
	if repeated.Code != http.StatusOK {
		t.Fatalf("repeated checkpoint status = %d, want %d", repeated.Code, http.StatusOK)
	}
	if offset, ok := store.Current("orders"); !ok || offset != 42 {
		t.Fatalf("stored checkpoint = (%d, %t), want (42, true)", offset, ok)
	}
}

func TestHandlerRejectsOlderCheckpoint(t *testing.T) {
	store := NewStore()
	handler := Handler(store)
	postCheckpoint(t, handler, "orders", `{"offset":42}`)

	response := postCheckpoint(t, handler, "orders", `{"offset":41}`)
	if response.Code != http.StatusConflict {
		t.Fatalf("older checkpoint status = %d, want %d", response.Code, http.StatusConflict)
	}
	if offset, _ := store.Current("orders"); offset != 42 {
		t.Fatalf("stored checkpoint = %d after stale request, want 42", offset)
	}
}

func TestHandlerAdvancesNewerCheckpoint(t *testing.T) {
	store := NewStore()
	handler := Handler(store)
	postCheckpoint(t, handler, "orders", `{"offset":42}`)

	response := postCheckpoint(t, handler, "orders", `{"offset":43}`)
	if response.Code != http.StatusCreated {
		t.Fatalf("newer checkpoint status = %d, want %d", response.Code, http.StatusCreated)
	}
	if offset, _ := store.Current("orders"); offset != 43 {
		t.Fatalf("stored checkpoint = %d, want 43", offset)
	}
}

func TestStoreValidationAndIsolation(t *testing.T) {
	store := NewStore()
	if _, err := store.Advance("", 1); !errors.Is(err, ErrEmptyStream) {
		t.Fatalf("empty stream error = %v, want %v", err, ErrEmptyStream)
	}
	if _, err := store.Advance("orders", -1); !errors.Is(err, ErrNegativeOffset) {
		t.Fatalf("negative offset error = %v, want %v", err, ErrNegativeOffset)
	}
	if _, err := store.Advance("orders", 3); err != nil {
		t.Fatalf("advance orders: %v", err)
	}
	if _, err := store.Advance("payments", 1); err != nil {
		t.Fatalf("advance payments: %v", err)
	}
	if offset, _ := store.Current("orders"); offset != 3 {
		t.Fatalf("orders checkpoint = %d, want 3", offset)
	}
}

func TestHandlerValidatesRequests(t *testing.T) {
	store := NewStore()
	handler := Handler(store)

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		status int
	}{
		{name: "method", method: http.MethodGet, path: "/checkpoints/orders", status: http.StatusMethodNotAllowed},
		{name: "missing stream", method: http.MethodPost, path: "/checkpoints/", body: `{"offset":1}`, status: http.StatusBadRequest},
		{name: "nested stream", method: http.MethodPost, path: "/checkpoints/a/b", body: `{"offset":1}`, status: http.StatusBadRequest},
		{name: "malformed json", method: http.MethodPost, path: "/checkpoints/orders", body: `{`, status: http.StatusBadRequest},
		{name: "unknown field", method: http.MethodPost, path: "/checkpoints/orders", body: `{"offset":1,"force":true}`, status: http.StatusBadRequest},
		{name: "trailing json", method: http.MethodPost, path: "/checkpoints/orders", body: `{"offset":1}{}`, status: http.StatusBadRequest},
		{name: "negative offset", method: http.MethodPost, path: "/checkpoints/orders", body: `{"offset":-1}`, status: http.StatusBadRequest},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != test.status {
				t.Fatalf("status = %d, want %d", response.Code, test.status)
			}
		})
	}
}
