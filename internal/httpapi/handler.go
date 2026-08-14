package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/hejinhui2018/go-featureflags/internal/flags"
)

type FlagService interface {
	Enabled(rctx context.Context, tenant, name string) (bool, error)
}

type Handler struct {
	service FlagService
}

func NewHandler(service FlagService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet || !strings.HasPrefix(r.URL.Path, "/v1/flags/") {
		http.NotFound(w, r)
		return
	}

	name := strings.TrimPrefix(r.URL.Path, "/v1/flags/")
	tenant := r.Header.Get("X-Tenant-ID")
	enabled, err := h.service.Enabled(r.Context(), tenant, name)
	if err != nil {
		switch {
		case errors.Is(err, flags.ErrInvalidArgument):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, flags.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"tenant":  tenant,
		"flag":    name,
		"enabled": enabled,
	})
}
