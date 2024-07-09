package main

import (
	"net/http"
)

// Define our struct
type SimpleHelmet struct {
	headers map[string]string
}

// Initialize it somewhere
func (helmet *SimpleHelmet) Default() {
	helmet.headers = make(map[string]string)
	helmet.headers["Cache-Control"] = "no-cache, no-store, max-age=0"
	helmet.headers["Content-Security-Policy"] = "frame-ancestors 'none'; default-src 'none'"
	helmet.headers["Pragma"] = "no-cache"
	helmet.headers["Referrer-Policy"] = "no-referrer"
	helmet.headers["Strict-Transport-Security"] = "max-age=31536000; includeSubDomains"
	helmet.headers["X-Content-Type-Options"] = "nosniff"
	helmet.headers["X-Dns-Prefetch-Control"] = "on"
	helmet.headers["X-Download-Options"] = "noopen"
	helmet.headers["X-Frame-Options"] = "DENY"
	helmet.headers["X-Permitted-Cross-Domain-Policies"] = "none"
	helmet.headers["X-Xss-Protection"] = "mode=block"
}

// Secure function, which will be called for each request
func (helmet *SimpleHelmet) Secure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, val := range helmet.headers {
			w.Header().Set(key, val)
		}
		next.ServeHTTP(w, r)
	})
}
