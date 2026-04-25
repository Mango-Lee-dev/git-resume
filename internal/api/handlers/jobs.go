package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wootaiklee/git-resume/internal/api/dto"
	"github.com/wootaiklee/git-resume/internal/api/jobs"
)

// JobsHandler handles job-related endpoints
type JobsHandler struct {
	jobManager *jobs.Manager
}

// NewJobsHandler creates a new jobs handler
func NewJobsHandler(jobManager *jobs.Manager) *JobsHandler {
	return &JobsHandler{jobManager: jobManager}
}

// Get handles GET /api/jobs/:id
func (h *JobsHandler) Get(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")

	job, found := h.jobManager.Get(jobID)
	if !found {
		respondError(w, r, dto.NewJobNotFoundError(jobID))
		return
	}

	respondOK(w, job.ToResponse())
}

// Cancel handles DELETE /api/jobs/:id
func (h *JobsHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")

	job, found := h.jobManager.Get(jobID)
	if !found {
		respondError(w, r, dto.NewJobNotFoundError(jobID))
		return
	}

	if err := h.jobManager.Cancel(jobID); err != nil {
		respondError(w, r, dto.NewBadRequestError(err.Error()))
		return
	}

	// Return the updated job status
	job, _ = h.jobManager.Get(jobID)
	respondOK(w, job.ToResponse())
}

// List handles GET /api/jobs
func (h *JobsHandler) List(w http.ResponseWriter, r *http.Request) {
	allJobs := h.jobManager.List()

	responses := make([]dto.JobResponse, len(allJobs))
	for i, job := range allJobs {
		responses[i] = job.ToResponse()
	}

	respondOK(w, map[string]interface{}{
		"jobs":  responses,
		"total": len(responses),
	})
}
