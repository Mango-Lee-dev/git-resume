package state

import (
	"time"

	"github.com/wootaiklee/git-resume/pkg/models"
)

// DateRangeMode indicates how date range was specified
type DateRangeMode int

const (
	DateRangeMonth DateRangeMode = iota
	DateRangeCustom
)

// AppState holds all configuration gathered from TUI screens
type AppState struct {
	// Repository selection
	RepoPaths []string

	// Date range
	DateMode DateRangeMode
	Month    int // 1-12
	Year     int
	FromDate time.Time
	ToDate   time.Time

	// Template
	TemplateName string

	// Options
	BatchSize int
	DryRun    bool
	Notify    bool

	// Export
	ExportFormat models.ExportFormat
	OutputPath   string

	// Results (populated during analysis)
	Commits    []models.Commit
	Results    []models.AnalysisResult
	TokensUsed int

	// Progress tracking
	TotalBatches  int
	CurrentBatch  int
	AnalysisError error
}

// NewAppState creates a new app state with sensible defaults
func NewAppState() *AppState {
	now := time.Now()
	return &AppState{
		RepoPaths:    []string{},
		DateMode:     DateRangeMonth,
		Month:        int(now.Month()),
		Year:         now.Year(),
		TemplateName: "default",
		BatchSize:    5,
		DryRun:       false,
		Notify:       false,
		ExportFormat: models.FormatCSV,
	}
}

// GetDateRange returns the effective date range based on the mode
func (s *AppState) GetDateRange() (from, to time.Time) {
	if s.DateMode == DateRangeMonth {
		from = time.Date(s.Year, time.Month(s.Month), 1, 0, 0, 0, 0, time.Local)
		to = from.AddDate(0, 1, 0).Add(-time.Second)
	} else {
		from = s.FromDate
		// Ensure 'to' includes the entire day
		to = time.Date(s.ToDate.Year(), s.ToDate.Month(), s.ToDate.Day(), 23, 59, 59, 0, time.Local)
	}
	return
}

// HasResults returns true if analysis has completed with results
func (s *AppState) HasResults() bool {
	return len(s.Results) > 0
}

// Reset clears the results and analysis state
func (s *AppState) Reset() {
	s.Commits = nil
	s.Results = nil
	s.TokensUsed = 0
	s.TotalBatches = 0
	s.CurrentBatch = 0
	s.AnalysisError = nil
	s.OutputPath = ""
}
