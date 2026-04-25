package handlers

import (
	"net/http"

	"github.com/wootaiklee/git-resume/internal/api/dto"
	"github.com/wootaiklee/git-resume/internal/service"
)

// StatsHandler handles statistics endpoints
type StatsHandler struct {
	results *service.ResultsService
}

// NewStatsHandler creates a new stats handler
func NewStatsHandler(results *service.ResultsService) *StatsHandler {
	return &StatsHandler{results: results}
}

// GetStats handles GET /api/stats
func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.results.GetStats(r.Context())
	if err != nil {
		respondError(w, r, dto.NewInternalError("failed to get statistics: "+err.Error()))
		return
	}

	respondOK(w, dto.StatsResponse{
		TotalResults: stats.TotalResults,
		TokensUsed: dto.TokenUsageStats{
			InputTokens:  stats.InputTokens,
			OutputTokens: stats.OutputTokens,
			TotalCost:    stats.TotalCost,
		},
		CategoryBreakdown: stats.CategoryBreakdown,
		ProjectBreakdown:  stats.ProjectBreakdown,
	})
}
