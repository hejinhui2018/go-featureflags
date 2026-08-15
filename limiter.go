package httpratelimit

import (
	"net"
	"net/http"
	"sync"
)

// Limiter applies a fixed request allowance per client address.
type Limiter struct {
	mu     sync.Mutex
	limit  int
	counts map[string]int
}

func New(limit int) *Limiter {
	if limit < 1 {
		panic("limit must be positive")
	}
	return &Limiter{limit: limit, counts: make(map[string]int)}
}

// clientIPFromAddr extracts the IP portion from an address that may
// include a port. It handles both IPv4 ("host:port") and IPv6
// ("[host]:port") formats. If the address cannot be parsed it is
// returned as-is so that rate-limiting still applies.
func clientIPFromAddr(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}

func (l *Limiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := clientIPFromAddr(r.RemoteAddr)
		l.mu.Lock()
		if l.counts[client] >= l.limit {
			l.mu.Unlock()
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		l.counts[client]++
		l.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
