package httpproxyheaders

import (
	"net/http"
	"strings"
)

var hopByHopHeaders = []string{
	"Connection",
	"Proxy-Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

// PrepareBackendRequest clones an inbound request for forwarding to a backend.
func PrepareBackendRequest(in *http.Request) *http.Request {
	out := in.Clone(in.Context())
	out.Header = in.Header.Clone()
	out.RequestURI = ""

	// RFC 7230 §6.1: the Connection header may nominate additional
	// header fields that are specific to this connection and must not
	// be forwarded.  Collect them before deleting the Connection
	// header itself.
	if conn := out.Header.Get("Connection"); conn != "" {
		for _, name := range strings.Split(conn, ",") {
			name = strings.TrimSpace(name)
			if name != "" {
				out.Header.Del(name)
			}
		}
	}

	for _, name := range hopByHopHeaders {
		out.Header.Del(name)
	}
	return out
}
