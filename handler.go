package leasecache

import (
	"net/http"
	"strings"
)

// Handler serves values at GET /leases/{key}.
func Handler(store *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		key := strings.TrimPrefix(r.URL.Path, "/leases/")
		if key == "" || key == r.URL.Path {
			http.Error(w, ErrEmptyKey.Error(), http.StatusBadRequest)
			return
		}

		value, ok := store.Get(key)
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(value)
	})
}
