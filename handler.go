package appendledger

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

// Handler accepts POST /streams/{stream}/records requests.
func Handler(store *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		stream := strings.TrimPrefix(r.URL.Path, "/streams/")
		stream = strings.TrimSuffix(stream, "/records")
		if stream == "" || stream == r.URL.Path || strings.Contains(stream, "/") {
			http.Error(w, ErrEmptyStream.Error(), http.StatusBadRequest)
			return
		}

		var record Record
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&record); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		created, err := store.Append(stream, record.Sequence, record.Value)
		switch {
		case errors.Is(err, ErrSequenceGap), errors.Is(err, ErrRecordConflict):
			http.Error(w, err.Error(), http.StatusConflict)
		case err != nil:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case created:
			w.WriteHeader(http.StatusCreated)
		default:
			w.WriteHeader(http.StatusOK)
		}
	})
}
