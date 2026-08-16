package httpcancel

import (
	"context"
	"net/http"
)

type contextKey string

const requestMarker contextKey = "request-marker"

// MarkRequest adds middleware-owned request metadata without replacing the
// cancellation and deadline signals carried by the incoming request.
func MarkRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(context.Background(), requestMarker, "edge")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Marker(ctx context.Context) string {
	value, _ := ctx.Value(requestMarker).(string)
	return value
}
