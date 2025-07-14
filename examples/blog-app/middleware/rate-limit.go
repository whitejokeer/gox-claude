package middleware

import (
    "net/http"
    "log"
)

// RateLimit is a middleware that rate-limit
func RateLimit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Pre-processing
        log.Printf("[RateLimit] Request: %s %s", r.Method, r.URL.Path)
        
        // TODO: Add your middleware logic here
        
        // Call the next handler
        next.ServeHTTP(w, r)
        
        // Post-processing (if needed)
    })
}

// RateLimitFunc is a middleware function that can be used with gorilla/mux
func RateLimitFunc(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Pre-processing
        log.Printf("[RateLimit] Request: %s %s", r.Method, r.URL.Path)
        
        // TODO: Add your middleware logic here
        
        // Call the next handler
        next(w, r)
        
        // Post-processing (if needed)
    }
}