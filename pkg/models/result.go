package models

import "time"

// AnalysisResult represents the LLM-generated resume bullet point
type AnalysisResult struct {
	ID            int64     `json:"id"`
	CommitHash    string    `json:"commit_hash"`
	Date          time.Time `json:"date"`
	Project       string    `json:"project"`
	Category      Category  `json:"category"`
	ImpactSummary string    `json:"impact_summary"`
	CreatedAt     time.Time `json:"created_at"`
}

// BatchResult represents results from a batch of commits
type BatchResult struct {
	Results     []AnalysisResult `json:"results"`
	TokensUsed  int              `json:"tokens_used"`
	ProcessedAt time.Time        `json:"processed_at"`
}

// ExportFormat represents the output format type
type ExportFormat string

const (
	FormatCSV      ExportFormat = "csv"
	FormatJSON     ExportFormat = "json"
	FormatMarkdown ExportFormat = "markdown"
)
