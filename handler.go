package checkpointstore

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type advanceRequest struct {
	Offset int64 `json:"offset"`
}

// Handler accepts POST /checkpoints/{stream} requests.
func Handler(store *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		stream := strings.TrimPrefix(r.URL.Path, "/checkpoints/")
		if stream == "" || stream == r.URL.Path || strings.Contains(stream, "/") {
			http.Error(w, ErrEmptyStream.Error(), http.StatusBadRequest)
			return
		}

		var request advanceRequest
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

		changed, err := store.Advance(stream, request.Offset)
		switch {
		case errors.Is(err, ErrStaleOffset):
			http.Error(w, err.Error(), http.StatusConflict)
		case err != nil:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case changed:
			w.WriteHeader(http.StatusCreated)
		default:
			w.WriteHeader(http.StatusOK)
		}
	})
}
