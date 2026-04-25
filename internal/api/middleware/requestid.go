package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

// RequestIDKey is the context key for request ID
const RequestIDKey contextKey = "request_id"

// RequestIDHeader is the header name for request ID
const RequestIDHeader = "X-Request-ID"

// RequestID returns a middleware that injects a request ID into the context
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set response header
		w.Header().Set(RequestIDHeader, requestID)

		// Add to context
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx context.Context) string {
	if rid := ctx.Value(RequestIDKey); rid != nil {
		return rid.(string)
	}
	return ""
}
