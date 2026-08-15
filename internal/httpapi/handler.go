package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/hejinhui2018/recordstore-validation-003/internal/store"
)

type Handler struct {
	store *store.Memory
}

func NewHandler(memory *store.Memory) http.Handler {
	return &Handler{store: memory}
}

type response struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
	Error string `json:"error,omitempty"`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, response{Error: "method not allowed"})
		return
	}
	key := strings.TrimPrefix(r.URL.Path, "/v1/records/")
	if strings.TrimSpace(key) == "" {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid key"})
		return
	}

	value, err := h.store.Get(r.Context(), key)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeJSON(w, http.StatusNotFound, response{Error: "not found"})
		case errors.Is(err, context.Canceled):
			writeJSON(w, http.StatusRequestTimeout, response{Error: "request canceled"})
		default:
			writeJSON(w, http.StatusInternalServerError, response{Error: "internal error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, response{Key: strings.TrimSpace(key), Value: value})
}

func writeJSON(w http.ResponseWriter, status int, value response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
