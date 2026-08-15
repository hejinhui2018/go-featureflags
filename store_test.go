package quotareservoir

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testStore(t *testing.T, capacity int64) *Store {
	t.Helper()
	store, err := NewStore(capacity)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return store
}

func postBatch(t *testing.T, handler http.Handler, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, "/reservations/batch", strings.NewReader(body))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func TestHandlerRejectsBatchWithoutPartialReservations(t *testing.T) {
	store := testStore(t, 10)
	response := postBatch(t, Handler(store), `{"reservations":[{"name":"alpha","units":4},{"name":"beta","units":7}]}`)
	if response.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusConflict)
	}
	if used := store.Used(); used != 0 {
		t.Fatalf("used after rejected batch = %d, want 0", used)
	}
	if _, ok := store.Reservation("alpha"); ok {
		t.Fatal("alpha was partially reserved")
	}
}

func TestHandlerAcceptsExactCapacity(t *testing.T) {
	store := testStore(t, 10)
	response := postBatch(t, Handler(store), `{"reservations":[{"name":"alpha","units":4},{"name":"beta","units":6}]}`)
	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusCreated)
	}
	if used := store.Used(); used != 10 {
		t.Fatalf("used = %d, want 10", used)
	}
}

func TestStorePreservesExistingReservationsOnFailure(t *testing.T) {
	store := testStore(t, 10)
	if err := store.ReserveBatch([]Reservation{{Name: "alpha", Units: 3}}); err != nil {
		t.Fatalf("initial reservation: %v", err)
	}
	if err := store.ReserveBatch([]Reservation{{Name: "beta", Units: 4}, {Name: "gamma", Units: 4}}); !errors.Is(err, ErrCapacityExceeded) {
		t.Fatalf("over-capacity error = %v, want %v", err, ErrCapacityExceeded)
	}
	if used := store.Used(); used != 3 {
		t.Fatalf("used after rejected batch = %d, want 3", used)
	}
	if _, ok := store.Reservation("beta"); ok {
		t.Fatal("beta was partially reserved")
	}
}

func TestHandlerValidatesBatch(t *testing.T) {
	store := testStore(t, 10)
	handler := Handler(store)
	tests := []struct {
		name string
		path string
		body string
		code int
	}{
		{name: "wrong path", path: "/reservations", body: `{"reservations":[]}`, code: http.StatusMethodNotAllowed},
		{name: "empty batch", path: "/reservations/batch", body: `{"reservations":[]}`, code: http.StatusBadRequest},
		{name: "empty name", path: "/reservations/batch", body: `{"reservations":[{"name":"","units":1}]}`, code: http.StatusBadRequest},
		{name: "invalid units", path: "/reservations/batch", body: `{"reservations":[{"name":"alpha","units":0}]}`, code: http.StatusBadRequest},
		{name: "duplicate name", path: "/reservations/batch", body: `{"reservations":[{"name":"alpha","units":1},{"name":"alpha","units":1}]}`, code: http.StatusBadRequest},
		{name: "unknown field", path: "/reservations/batch", body: `{"reservations":[{"name":"alpha","units":1,"extra":true}]}`, code: http.StatusBadRequest},
		{name: "trailing json", path: "/reservations/batch", body: `{"reservations":[{"name":"alpha","units":1}]}{}`, code: http.StatusBadRequest},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.path, strings.NewReader(test.body))
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != test.code {
				t.Fatalf("status = %d, want %d", response.Code, test.code)
			}
		})
	}
}

func TestNewStoreRejectsInvalidCapacity(t *testing.T) {
	if _, err := NewStore(0); !errors.Is(err, ErrInvalidCapacity) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidCapacity)
	}
}
