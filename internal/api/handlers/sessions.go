package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wootaiklee/git-resume/internal/api/dto"
	"github.com/wootaiklee/git-resume/internal/api/session"
)

type SessionsHandler struct {
	manager *session.Manager
}

func NewSessionsHandler(manager *session.Manager) *SessionsHandler {
	return &SessionsHandler{manager: manager}
}

type CreateSessionRequest struct {
	APIKey string `json:"api_key"`
}

type CreateSessionResponse struct {
	SessionID string `json:"session_id"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
	Message   string `json:"message"`
}

type SessionResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
	LastUsed  string `json:"last_used"`
}

// Create creates a new session with the provided API key
func (h *SessionsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, r, dto.NewBadRequestError("Invalid request body"))
		return
	}

	if req.APIKey == "" {
		respondError(w, r, dto.NewBadRequestError("API key is required"))
		return
	}

	sess, err := h.manager.CreateSession(req.APIKey)
	if err != nil {
		if errors.Is(err, session.ErrInvalidAPIKey) {
			respondError(w, r, dto.NewUnauthorizedError("Invalid API key"))
			return
		}
		respondError(w, r, dto.NewInternalError("Failed to create session"))
		return
	}

	respondJSON(w, http.StatusCreated, CreateSessionResponse{
		SessionID: sess.ID,
		CreatedAt: sess.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		ExpiresAt: sess.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
		Message:   "Session created successfully",
	})
}

// Get returns session information
func (h *SessionsHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	if sessionID == "" {
		respondError(w, r, dto.NewBadRequestError("Session ID is required"))
		return
	}

	sess, err := h.manager.GetSession(sessionID)
	if err != nil {
		if errors.Is(err, session.ErrSessionNotFound) {
			respondError(w, r, dto.NewNotFoundError("Session"))
			return
		}
		if errors.Is(err, session.ErrSessionExpired) {
			respondError(w, r, dto.NewUnauthorizedError("Session expired"))
			return
		}
		respondError(w, r, dto.NewInternalError("Failed to get session"))
		return
	}

	respondJSON(w, http.StatusOK, SessionResponse{
		ID:        sess.ID,
		CreatedAt: sess.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		ExpiresAt: sess.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
		LastUsed:  sess.LastUsed.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// Delete removes a session
func (h *SessionsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	if sessionID == "" {
		respondError(w, r, dto.NewBadRequestError("Session ID is required"))
		return
	}

	h.manager.DeleteSession(sessionID)
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Session deleted successfully",
	})
}

// Extend extends the session expiration
func (h *SessionsHandler) Extend(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	if sessionID == "" {
		respondError(w, r, dto.NewBadRequestError("Session ID is required"))
		return
	}

	if err := h.manager.ExtendSession(sessionID); err != nil {
		if errors.Is(err, session.ErrSessionNotFound) {
			respondError(w, r, dto.NewNotFoundError("Session"))
			return
		}
		respondError(w, r, dto.NewInternalError("Failed to extend session"))
		return
	}

	sess, _ := h.manager.GetSession(sessionID)
	respondJSON(w, http.StatusOK, SessionResponse{
		ID:        sess.ID,
		CreatedAt: sess.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		ExpiresAt: sess.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
		LastUsed:  sess.LastUsed.Format("2006-01-02T15:04:05Z07:00"),
	})
}
