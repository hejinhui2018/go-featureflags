package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/hejinhui2018/configstore-validation-004/internal/store"
)

type Handler struct {
	store *store.Memory
}

func NewHandler(memory *store.Memory) http.Handler {
	return &Handler{store: memory}
}

type writeRequest struct {
	Value string `json:"value"`
}

type response struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
	Error string `json:"error,omitempty"`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/v1/records/")
	if strings.TrimSpace(key) == "" {
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid key"})
		return
	}

	switch r.Method {
	case http.MethodPut:
		var req writeRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, response{Error: "invalid JSON"})
			return
		}
		if req.Value == "" {
			writeJSON(w, http.StatusBadRequest, response{Error: "invalid value"})
			return
		}
		if err := h.store.Put(r.Context(), key, req.Value); err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, response{Key: strings.TrimSpace(key), Value: req.Value})
	case http.MethodGet:
		value, err := h.store.Get(r.Context(), key)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, response{Key: strings.TrimSpace(key), Value: value})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, response{Error: "method not allowed"})
	}
}

func writeStoreError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrInvalidKey):
		writeJSON(w, http.StatusBadRequest, response{Error: "invalid key"})
	case errors.Is(err, store.ErrNotFound):
		writeJSON(w, http.StatusNotFound, response{Error: "not found"})
	case errors.Is(err, context.Canceled):
		writeJSON(w, http.StatusRequestTimeout, response{Error: "request canceled"})
	default:
		writeJSON(w, http.StatusInternalServerError, response{Error: "internal error"})
	}
}

func writeJSON(w http.ResponseWriter, status int, value response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
