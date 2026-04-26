package middleware

import (
	"context"
	"net/http"

	"github.com/wootaiklee/git-resume/internal/api/session"
)

type sessionContextKey string

const (
	SessionIDKey sessionContextKey = "session_id"
	APIKeyKey    sessionContextKey = "api_key"
)

// SessionAuth middleware validates the session and injects the API key into context
func SessionAuth(manager *session.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionID := r.Header.Get("X-Session-ID")
			if sessionID == "" {
				http.Error(w, `{"error":"Session ID required","code":"SESSION_REQUIRED"}`, http.StatusUnauthorized)
				return
			}

			apiKey, err := manager.GetAPIKey(sessionID)
			if err != nil {
				if err == session.ErrSessionNotFound {
					http.Error(w, `{"error":"Session not found","code":"SESSION_NOT_FOUND"}`, http.StatusUnauthorized)
					return
				}
				if err == session.ErrSessionExpired {
					http.Error(w, `{"error":"Session expired","code":"SESSION_EXPIRED"}`, http.StatusUnauthorized)
					return
				}
				http.Error(w, `{"error":"Invalid session","code":"SESSION_INVALID"}`, http.StatusUnauthorized)
				return
			}

			// Inject session ID and API key into context
			ctx := context.WithValue(r.Context(), SessionIDKey, sessionID)
			ctx = context.WithValue(ctx, APIKeyKey, apiKey)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetSessionID retrieves the session ID from context
func GetSessionID(ctx context.Context) string {
	if id, ok := ctx.Value(SessionIDKey).(string); ok {
		return id
	}
	return ""
}

// GetAPIKey retrieves the API key from context
func GetAPIKey(ctx context.Context) string {
	if key, ok := ctx.Value(APIKeyKey).(string); ok {
		return key
	}
	return ""
}
