package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wootaiklee/git-resume/internal/api/dto"
	"github.com/wootaiklee/git-resume/internal/service"
)

// ResultsHandler handles result-related endpoints
type ResultsHandler struct {
	results *service.ResultsService
}

// NewResultsHandler creates a new results handler
func NewResultsHandler(results *service.ResultsService) *ResultsHandler {
	return &ResultsHandler{results: results}
}

// List handles GET /api/results
func (h *ResultsHandler) List(w http.ResponseWriter, r *http.Request) {
	opts := parseResultsQuery(r)

	results, total, err := h.results.Query(r.Context(), opts)
	if err != nil {
		respondError(w, r, dto.NewInternalError("failed to query results: "+err.Error()))
		return
	}

	// Convert to response DTOs
	items := make([]dto.ResultResponse, len(results))
	for i, result := range results {
		items[i] = dto.FromAnalysisResult(result)
	}

	totalPages := (total + opts.PageSize - 1) / opts.PageSize
	if totalPages == 0 {
		totalPages = 1
	}

	respondOK(w, dto.ResultsListResponse{
		Results:    items,
		Total:      total,
		Page:       opts.Page,
		PageSize:   opts.PageSize,
		TotalPages: totalPages,
	})
}

// Get handles GET /api/results/:id
func (h *ResultsHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, r, dto.NewBadRequestError("invalid result ID"))
		return
	}

	result, err := h.results.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, r, dto.NewInternalError("failed to get result"))
		return
	}

	if result == nil {
		respondError(w, r, dto.NewNotFoundError("result"))
		return
	}

	respondOK(w, dto.FromAnalysisResult(*result))
}

// parseResultsQuery parses query parameters for results list
func parseResultsQuery(r *http.Request) service.QueryOptions {
	opts := service.QueryOptions{
		Project:  r.URL.Query().Get("project"),
		Category: r.URL.Query().Get("category"),
	}

	if fromStr := r.URL.Query().Get("from"); fromStr != "" {
		if t, err := time.Parse("2006-01-02", fromStr); err == nil {
			opts.FromDate = t
		}
	}

	if toStr := r.URL.Query().Get("to"); toStr != "" {
		if t, err := time.Parse("2006-01-02", toStr); err == nil {
			opts.ToDate = t.Add(24*time.Hour - time.Second)
		}
	}

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			opts.Page = page
		}
	}

	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			opts.PageSize = pageSize
		}
	}

	opts.SetDefaults()
	return opts
}
