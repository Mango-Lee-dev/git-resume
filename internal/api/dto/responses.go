package dto

import (
	"time"

	"github.com/wootaiklee/git-resume/pkg/models"
)

// JobStatus represents the status of an async job
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
)

// AnalyzeStartResponse returned by POST /api/analyze
type AnalyzeStartResponse struct {
	JobID   string `json:"job_id"`
	Message string `json:"message"`
}

// JobResponse represents an async job status
type JobResponse struct {
	ID          string     `json:"id"`
	Status      JobStatus  `json:"status"`
	Progress    int        `json:"progress"`
	Phase       string     `json:"phase,omitempty"`
	Message     string     `json:"message"`
	CreatedAt   time.Time  `json:"created_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	ResultCount int        `json:"result_count,omitempty"`
	Error       string     `json:"error,omitempty"`
}

// ResultResponse wraps a single AnalysisResult
type ResultResponse struct {
	ID            int64     `json:"id"`
	CommitHash    string    `json:"commit_hash"`
	Date          time.Time `json:"date"`
	Project       string    `json:"project"`
	Category      string    `json:"category"`
	ImpactSummary string    `json:"impact_summary"`
	CreatedAt     time.Time `json:"created_at"`
}

// FromAnalysisResult converts a model to response DTO
func FromAnalysisResult(r models.AnalysisResult) ResultResponse {
	return ResultResponse{
		ID:            r.ID,
		CommitHash:    r.CommitHash,
		Date:          r.Date,
		Project:       r.Project,
		Category:      string(r.Category),
		ImpactSummary: r.ImpactSummary,
		CreatedAt:     r.CreatedAt,
	}
}

// ResultsListResponse wraps paginated results
type ResultsListResponse struct {
	Results    []ResultResponse `json:"results"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

// TemplateResponse represents a template in API responses
type TemplateResponse struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ToneStyle   string   `json:"tone_style"`
	Focus       []string `json:"focus"`
}

// TemplatesListResponse wraps the templates list
type TemplatesListResponse struct {
	Templates []TemplateResponse `json:"templates"`
}

// StatsResponse for GET /api/stats
type StatsResponse struct {
	TotalResults      int                 `json:"total_results"`
	TotalCommits      int                 `json:"total_commits"`
	TokensUsed        TokenUsageStats     `json:"tokens_used"`
	CategoryBreakdown map[string]int      `json:"category_breakdown"`
	ProjectBreakdown  map[string]int      `json:"project_breakdown"`
	RecentActivity    []DailyActivityStat `json:"recent_activity,omitempty"`
}

// TokenUsageStats represents token usage statistics
type TokenUsageStats struct {
	InputTokens  int     `json:"input_tokens"`
	OutputTokens int     `json:"output_tokens"`
	TotalCost    float64 `json:"total_cost"`
}

// DailyActivityStat represents daily activity count
type DailyActivityStat struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// HealthResponse for GET /health
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// ExportMetadata for export response
type ExportMetadata struct {
	GeneratedAt time.Time `json:"generated_at"`
	FromDate    string    `json:"from_date,omitempty"`
	ToDate      string    `json:"to_date,omitempty"`
	TotalCount  int       `json:"total_count"`
	Format      string    `json:"format"`
}
