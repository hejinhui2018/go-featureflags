package idempotency

import (
	"bytes"
	"net/http"
	"sync"
)

type cachedResponse struct {
	status int
	header http.Header
	body   []byte
}

// inflightEntry tracks a request that is currently being processed.
// The done channel is closed once the response is ready; any concurrent
// request for the same key waits on it and then reuses the stored response.
type inflightEntry struct {
	done     chan struct{}
	response cachedResponse
}

// Middleware reuses a completed response for repeated idempotency keys.
type Middleware struct {
	mu        sync.Mutex
	responses map[string]*inflightEntry
}

func New() *Middleware {
	return &Middleware{responses: make(map[string]*inflightEntry)}
}

func (m *Middleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			next.ServeHTTP(w, r)
			return
		}

		m.mu.Lock()
		entry, ok := m.responses[key]
		if !ok {
			// First request for this key: register an in-flight entry
			// so concurrent requests for the same key can wait on it.
			entry = &inflightEntry{done: make(chan struct{})}
			m.responses[key] = entry
			m.mu.Unlock()

			// Execute the business logic exactly once.
			buffered := newBufferedResponse()
			next.ServeHTTP(buffered, r)
			response := buffered.snapshot()

			entry.response = response
			close(entry.done)

			writeResponse(w, response)
			return
		}
		m.mu.Unlock()

		// Another request for this key is in-flight (or already completed).
		// Wait for the result, then reuse the same status, headers, and body.
		<-entry.done
		writeResponse(w, entry.response)
	})
}

type bufferedResponse struct {
	header http.Header
	status int
	body   bytes.Buffer
}

func newBufferedResponse() *bufferedResponse {
	return &bufferedResponse{header: make(http.Header)}
}

func (w *bufferedResponse) Header() http.Header {
	return w.header
}

func (w *bufferedResponse) WriteHeader(status int) {
	if w.status == 0 {
		w.status = status
	}
}

func (w *bufferedResponse) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.body.Write(p)
}

func (w *bufferedResponse) snapshot() cachedResponse {
	status := w.status
	if status == 0 {
		status = http.StatusOK
	}
	return cachedResponse{
		status: status,
		header: w.header.Clone(),
		body:   append([]byte(nil), w.body.Bytes()...),
	}
}

func writeResponse(w http.ResponseWriter, response cachedResponse) {
	for name, values := range response.header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	w.WriteHeader(response.status)
	_, _ = w.Write(response.body)
}
