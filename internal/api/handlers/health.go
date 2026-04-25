package handlers

import (
	"net/http"
	"time"

	"github.com/wootaiklee/git-resume/internal/api/dto"
	"github.com/wootaiklee/git-resume/internal/db"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db *db.DB
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *db.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Health handles GET /health
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	respondOK(w, dto.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// Ready handles GET /ready
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	// Check database connection
	if h.db == nil {
		respondJSON(w, http.StatusServiceUnavailable, dto.HealthResponse{
			Status:    "not ready",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	respondOK(w, dto.HealthResponse{
		Status:    "ready",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}
