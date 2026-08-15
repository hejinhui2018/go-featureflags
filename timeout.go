package httptimeout

import (
	"bytes"
	"net/http"
	"sync"
	"time"
)

// Timeout limits a handler's response time and keeps work that finishes later
// from changing the response already sent to the client.
func Timeout(limit time.Duration, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := &bufferedWriter{header: make(http.Header), status: http.StatusOK}
		finished := make(chan struct{})
		go func() {
			next.ServeHTTP(buf, r)
			close(finished)
		}()

		timer := time.NewTimer(limit)
		defer timer.Stop()
		select {
		case <-finished:
			bufferedResponse{
				header: buf.header,
				status: buf.status,
				body:   append([]byte(nil), buf.body.Bytes()...),
			}.writeTo(w)
		case <-timer.C:
			writeTimeout(w)
		}
	})
}

type bufferedResponse struct {
	header http.Header
	status int
	body   []byte
}

func runBuffered(next http.Handler, r *http.Request) bufferedResponse {
	recorder := &bufferedWriter{header: make(http.Header), status: http.StatusOK}
	next.ServeHTTP(recorder, r)
	return bufferedResponse{header: recorder.header, status: recorder.status, body: append([]byte(nil), recorder.body.Bytes()...)}
}

func (response bufferedResponse) writeTo(w http.ResponseWriter) {
	for key, values := range response.header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(response.status)
	_, _ = w.Write(response.body)
}

func writeTimeout(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusGatewayTimeout)
	_, _ = w.Write([]byte("request timed out\n"))
}

type bufferedWriter struct {
	header http.Header
	status int
	body   bytes.Buffer
	mu     sync.Mutex
}

func (w *bufferedWriter) Header() http.Header { return w.header }

func (w *bufferedWriter) WriteHeader(status int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.status == http.StatusOK && w.body.Len() == 0 {
		w.status = status
	}
}

func (w *bufferedWriter) Write(data []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.body.Write(data)
}
