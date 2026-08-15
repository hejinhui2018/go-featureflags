package httplimitercancel

import "net/http"

// LimitConcurrent allows at most limit requests to execute the wrapped handler.
func LimitConcurrent(next http.Handler, limit int) http.Handler {
	if limit < 1 {
		panic("limit must be positive")
	}
	sem := make(chan struct{}, limit)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sem <- struct{}{}
		defer func() { <-sem }()
		if r.Context().Err() != nil {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		next.ServeHTTP(w, r)
	})
}
