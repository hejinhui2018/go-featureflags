package httpbodyclose

import (
	"bytes"
	"io"
	"net/http"
)

// LimitMiddleware buffers request bodies up to Limit before invoking Next.
type LimitMiddleware struct {
	Limit int64
	Next  http.Handler
}

func (m LimitMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.Next == nil {
		m.Next = http.NotFoundHandler()
	}
	limited := http.MaxBytesReader(w, r.Body, m.Limit)
	data, err := io.ReadAll(limited)
	if err != nil {
		r.Body.Close()
		http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
		return
	}
	if err := r.Body.Close(); err != nil {
		http.Error(w, "request body close failed", http.StatusBadRequest)
		return
	}
	r.Body = io.NopCloser(bytes.NewReader(data))
	m.Next.ServeHTTP(w, r)
}
