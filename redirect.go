package httpredirect

import "net/http"

// RedirectToCanonical sends requests to the canonical endpoint.
//
// A 308 Permanent Redirect is used so that the original request method
// and body are preserved across the redirect. This is important for POST
// requests whose body must reach the canonical endpoint intact.
func RedirectToCanonical(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/canonical", http.StatusPermanentRedirect)
}
