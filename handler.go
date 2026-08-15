package quotareservoir

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type batchRequest struct {
	Reservations []Reservation `json:"reservations"`
}

// Handler accepts POST /reservations/batch requests.
func Handler(store *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/reservations/batch" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var request batchRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&request); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := store.ReserveBatch(request.Reservations); err != nil {
			if errors.Is(err, ErrCapacityExceeded) {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})
}
