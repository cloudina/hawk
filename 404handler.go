package main

import (
	"net/http"
)

// Handle404 ...
func Handle404(helmet SimpleHelmet) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, val := range helmet.headers {
			w.Header().Set(key,val)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
	})
}