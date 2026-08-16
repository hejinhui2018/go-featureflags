package appendledger

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func postRecord(t *testing.T, handler http.Handler, stream, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, "/streams/"+stream+"/records", strings.NewReader(body))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func TestHandlerRejectsSequenceGapWithoutAdvancing(t *testing.T) {
	store := NewStore()
	handler := Handler(store)
	if response := postRecord(t, handler, "orders", `{"sequence":1,"value":"created"}`); response.Code != http.StatusCreated {
		t.Fatalf("first status = %d, want %d", response.Code, http.StatusCreated)
	}
	response := postRecord(t, handler, "orders", `{"sequence":3,"value":"shipped"}`)
	if response.Code != http.StatusConflict {
		t.Fatalf("gap status = %d, want %d", response.Code, http.StatusConflict)
	}
	if next := store.NextSequence("orders"); next != 2 {
		t.Fatalf("next sequence after gap = %d, want 2", next)
	}
}

func TestHandlerAppendsContiguousRecords(t *testing.T) {
	store := NewStore()
	handler := Handler(store)
	postRecord(t, handler, "orders", `{"sequence":1,"value":"created"}`)
	response := postRecord(t, handler, "orders", `{"sequence":2,"value":"paid"}`)
	if response.Code != http.StatusCreated {
		t.Fatalf("contiguous status = %d, want %d", response.Code, http.StatusCreated)
	}
	if next := store.NextSequence("orders"); next != 3 {
		t.Fatalf("next sequence = %d, want 3", next)
	}
}

func TestHandlerKeepsIdempotencyAndConflictSemantics(t *testing.T) {
	store := NewStore()
	handler := Handler(store)
	postRecord(t, handler, "orders", `{"sequence":1,"value":"created"}`)
	if response := postRecord(t, handler, "orders", `{"sequence":1,"value":"created"}`); response.Code != http.StatusOK {
		t.Fatalf("repeat status = %d, want %d", response.Code, http.StatusOK)
	}
	if response := postRecord(t, handler, "orders", `{"sequence":1,"value":"different"}`); response.Code != http.StatusConflict {
		t.Fatalf("conflict status = %d, want %d", response.Code, http.StatusConflict)
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
		code   int
	}{
		{name: "method", method: http.MethodGet, path: "/streams/orders/records", code: http.StatusMethodNotAllowed},
		{name: "missing stream", method: http.MethodPost, path: "/streams//records", body: `{"sequence":1,"value":"x"}`, code: http.StatusBadRequest},
		{name: "invalid sequence", method: http.MethodPost, path: "/streams/orders/records", body: `{"sequence":0,"value":"x"}`, code: http.StatusBadRequest},
		{name: "empty value", method: http.MethodPost, path: "/streams/orders/records", body: `{"sequence":1,"value":""}`, code: http.StatusBadRequest},
		{name: "unknown field", method: http.MethodPost, path: "/streams/orders/records", body: `{"sequence":1,"value":"x","extra":true}`, code: http.StatusBadRequest},
		{name: "trailing json", method: http.MethodPost, path: "/streams/orders/records", body: `{"sequence":1,"value":"x"}{}`, code: http.StatusBadRequest},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != test.code {
				t.Fatalf("status = %d, want %d", response.Code, test.code)
			}
		})
	}
}

func TestStoreValidation(t *testing.T) {
	store := NewStore()
	if _, err := store.Append("", 1, "x"); !errors.Is(err, ErrEmptyStream) {
		t.Fatalf("empty stream error = %v, want %v", err, ErrEmptyStream)
	}
	if _, err := store.Append("orders", 0, "x"); !errors.Is(err, ErrInvalidSeq) {
		t.Fatalf("invalid sequence error = %v, want %v", err, ErrInvalidSeq)
	}
	if _, err := store.Append("orders", 1, ""); !errors.Is(err, ErrEmptyValue) {
		t.Fatalf("empty value error = %v, want %v", err, ErrEmptyValue)
	}
}
