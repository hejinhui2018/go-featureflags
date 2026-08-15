package idemqueue

import (
	"errors"
	"io"
	"net/http"
)

// Handler returns an HTTP handler backed by processor and execute.
func Handler(processor *Processor, execute Executor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		result, err := processor.Process(r.Context(), r.Header.Get("Idempotency-Key"), payload, execute)
		switch {
		case errors.Is(err, ErrMissingKey):
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, ErrKeyConflict):
			http.Error(w, err.Error(), http.StatusConflict)
			return
		case err != nil:
			http.Error(w, "request failed", http.StatusInternalServerError)
			return
		}

		status := result.Status
		if status == 0 {
			status = http.StatusOK
		}
		w.WriteHeader(status)
		_, _ = w.Write(result.Body)
	})
}
