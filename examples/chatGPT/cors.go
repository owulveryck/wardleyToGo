package main

import (
	"net/http"
)

// corsHandler is a middleware that enables CORS
func corsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the request origin
		origin := r.Header.Get("Origin")

		// Check if the request origin is allowed
		if isAllowedOrigin(origin) {
			// Set the CORS headers
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, openai-conversation-id, openai-ephemeral-user-id, openai-*, sentry-trace, baggage")
		}
		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue with the next handler
		next.ServeHTTP(w, r)
	})
}

// isAllowedOrigin checks if the given origin is allowed
func isAllowedOrigin(origin string) bool {
	allowedOrigins := []string{
		"http://localhost:3333",
		"https://chat.openai.com",
		// Add more allowed origins here
	}

	// Check if the origin is in the allowed origins list
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}

	return false
}
