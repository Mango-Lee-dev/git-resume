package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/wootaiklee/git-resume/internal/api/dto"
	"github.com/wootaiklee/git-resume/internal/api/middleware"
)

// respondJSON sends a JSON response with the given status code
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// Log error but can't do much else at this point
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

// respondError sends an error response
func respondError(w http.ResponseWriter, r *http.Request, err *dto.APIError) {
	// Add request ID to error if available
	if requestID := middleware.GetRequestID(r.Context()); requestID != "" {
		err = err.WithRequestID(requestID)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)

	json.NewEncoder(w).Encode(err)
}

// respondOK sends a 200 OK response with data
func respondOK(w http.ResponseWriter, data interface{}) {
	respondJSON(w, http.StatusOK, data)
}

// respondCreated sends a 201 Created response with data
func respondCreated(w http.ResponseWriter, data interface{}) {
	respondJSON(w, http.StatusCreated, data)
}

// respondAccepted sends a 202 Accepted response with data
func respondAccepted(w http.ResponseWriter, data interface{}) {
	respondJSON(w, http.StatusAccepted, data)
}

// respondNoContent sends a 204 No Content response
func respondNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
