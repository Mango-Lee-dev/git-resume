package service

import (
	"context"
	"time"

	"github.com/wootaiklee/git-resume/internal/db"
	"github.com/wootaiklee/git-resume/pkg/models"
)

// QueryOptions for filtering results
type QueryOptions struct {
	Project  string
	Category string
	FromDate time.Time
	ToDate   time.Time
	Page     int
	PageSize int
}

// SetDefaults sets default pagination values
func (o *QueryOptions) SetDefaults() {
	if o.Page <= 0 {
		o.Page = 1
	}
	if o.PageSize <= 0 || o.PageSize > 100 {
		o.PageSize = 20
	}
}

// ResultsService handles result queries
type ResultsService struct {
	cache *db.Cache
}

// NewResultsService creates a new results service
func NewResultsService(cache *db.Cache) *ResultsService {
	return &ResultsService{cache: cache}
}

// Query returns filtered and paginated results
func (s *ResultsService) Query(ctx context.Context, opts QueryOptions) ([]models.AnalysisResult, int, error) {
	opts.SetDefaults()

	// Get all results and filter
	allResults, err := s.cache.GetAllResults()
	if err != nil {
		return nil, 0, err
	}

	// Apply filters
	filtered := make([]models.AnalysisResult, 0)
	for _, r := range allResults {
		if !matchesFilters(r, opts) {
			continue
		}
		filtered = append(filtered, r)
	}

	total := len(filtered)

	// Apply pagination
	start := (opts.Page - 1) * opts.PageSize
	if start >= total {
		return []models.AnalysisResult{}, total, nil
	}

	end := start + opts.PageSize
	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

// GetByID returns a single result by ID
func (s *ResultsService) GetByID(ctx context.Context, id int64) (*models.AnalysisResult, error) {
	results, err := s.cache.GetAllResults()
	if err != nil {
		return nil, err
	}

	for _, r := range results {
		if r.ID == id {
			return &r, nil
		}
	}

	return nil, nil
}

// GetAll returns all results
func (s *ResultsService) GetAll(ctx context.Context) ([]models.AnalysisResult, error) {
	return s.cache.GetAllResults()
}

// GetByFilters returns results matching filters (no pagination)
func (s *ResultsService) GetByFilters(ctx context.Context, opts QueryOptions) ([]models.AnalysisResult, error) {
	allResults, err := s.cache.GetAllResults()
	if err != nil {
		return nil, err
	}

	filtered := make([]models.AnalysisResult, 0)
	for _, r := range allResults {
		if !matchesFilters(r, opts) {
			continue
		}
		filtered = append(filtered, r)
	}

	return filtered, nil
}

// GetStats returns statistics about the results
func (s *ResultsService) GetStats(ctx context.Context) (*Stats, error) {
	results, err := s.cache.GetAllResults()
	if err != nil {
		return nil, err
	}

	inputTokens, outputTokens, totalCost, err := s.cache.GetTotalTokenUsage()
	if err != nil {
		return nil, err
	}

	stats := &Stats{
		TotalResults:      len(results),
		CategoryBreakdown: make(map[string]int),
		ProjectBreakdown:  make(map[string]int),
		InputTokens:       inputTokens,
		OutputTokens:      outputTokens,
		TotalCost:         totalCost,
	}

	for _, r := range results {
		stats.CategoryBreakdown[string(r.Category)]++
		stats.ProjectBreakdown[r.Project]++
	}

	return stats, nil
}

// Stats holds result statistics
type Stats struct {
	TotalResults      int
	CategoryBreakdown map[string]int
	ProjectBreakdown  map[string]int
	InputTokens       int
	OutputTokens      int
	TotalCost         float64
}

// matchesFilters checks if a result matches the query filters
func matchesFilters(r models.AnalysisResult, opts QueryOptions) bool {
	if opts.Project != "" && r.Project != opts.Project {
		return false
	}

	if opts.Category != "" && string(r.Category) != opts.Category {
		return false
	}

	if !opts.FromDate.IsZero() && r.Date.Before(opts.FromDate) {
		return false
	}

	if !opts.ToDate.IsZero() && r.Date.After(opts.ToDate) {
		return false
	}

	return true
}
