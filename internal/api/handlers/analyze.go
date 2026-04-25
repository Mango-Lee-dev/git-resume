package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/wootaiklee/git-resume/internal/api/dto"
	"github.com/wootaiklee/git-resume/internal/api/jobs"
	"github.com/wootaiklee/git-resume/internal/service"
)

// AnalyzeHandler handles analysis endpoints
type AnalyzeHandler struct {
	jobManager *jobs.Manager
}

// NewAnalyzeHandler creates a new analyze handler
func NewAnalyzeHandler(jobManager *jobs.Manager) *AnalyzeHandler {
	return &AnalyzeHandler{jobManager: jobManager}
}

// Start handles POST /api/analyze
func (h *AnalyzeHandler) Start(w http.ResponseWriter, r *http.Request) {
	var req dto.AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, r, dto.NewBadRequestError("invalid JSON body"))
		return
	}

	// Validate request
	if errs := req.Validate(); len(errs) > 0 {
		respondError(w, r, dto.NewValidationError(errs))
		return
	}

	// Set defaults
	req.SetDefaults()

	// Convert to service config
	cfg, err := toAnalyzeConfig(req)
	if err != nil {
		respondError(w, r, dto.NewBadRequestError(err.Error()))
		return
	}

	// Submit job
	job, err := h.jobManager.Submit(cfg)
	if err != nil {
		respondError(w, r, dto.NewInternalError("failed to submit job"))
		return
	}

	respondAccepted(w, dto.AnalyzeStartResponse{
		JobID:   job.ID,
		Message: "Analysis job started",
	})
}

// toAnalyzeConfig converts request DTO to service config
func toAnalyzeConfig(req dto.AnalyzeRequest) (service.AnalyzeConfig, error) {
	cfg := service.AnalyzeConfig{
		Repos:     req.Repos,
		Template:  req.Template,
		BatchSize: req.BatchSize,
		DryRun:    req.DryRun,
	}

	now := time.Now()

	// Parse dates
	if req.FromDate != "" && req.ToDate != "" {
		from, err := time.Parse("2006-01-02", req.FromDate)
		if err != nil {
			return cfg, err
		}
		to, err := time.Parse("2006-01-02", req.ToDate)
		if err != nil {
			return cfg, err
		}
		cfg.FromDate = from
		cfg.ToDate = to.Add(24*time.Hour - time.Second)
	} else if req.Month > 0 {
		year := req.Year
		if year == 0 {
			year = now.Year()
		}
		cfg.FromDate = time.Date(year, time.Month(req.Month), 1, 0, 0, 0, 0, time.Local)
		cfg.ToDate = cfg.FromDate.AddDate(0, 1, 0).Add(-time.Second)
	} else {
		// Default: current month
		cfg.FromDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
		cfg.ToDate = now
	}

	return cfg, nil
}
