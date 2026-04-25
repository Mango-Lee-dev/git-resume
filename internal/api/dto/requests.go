package dto

// AnalyzeRequest represents the POST /api/analyze request body
type AnalyzeRequest struct {
	Repos     []string `json:"repos"`
	FromDate  string   `json:"from_date,omitempty"`
	ToDate    string   `json:"to_date,omitempty"`
	Month     int      `json:"month,omitempty"`
	Year      int      `json:"year,omitempty"`
	Template  string   `json:"template,omitempty"`
	BatchSize int      `json:"batch_size,omitempty"`
	DryRun    bool     `json:"dry_run,omitempty"`
}

// Validate validates the analyze request
func (r *AnalyzeRequest) Validate() map[string]string {
	errors := make(map[string]string)

	if len(r.Repos) == 0 {
		errors["repos"] = "at least one repository path is required"
	}

	if r.Month != 0 && (r.Month < 1 || r.Month > 12) {
		errors["month"] = "month must be between 1 and 12"
	}

	if r.BatchSize < 0 || r.BatchSize > 20 {
		errors["batch_size"] = "batch_size must be between 1 and 20"
	}

	// Check date format if provided
	if r.FromDate != "" && !isValidDateFormat(r.FromDate) {
		errors["from_date"] = "invalid date format, use YYYY-MM-DD"
	}

	if r.ToDate != "" && !isValidDateFormat(r.ToDate) {
		errors["to_date"] = "invalid date format, use YYYY-MM-DD"
	}

	return errors
}

// SetDefaults sets default values for optional fields
func (r *AnalyzeRequest) SetDefaults() {
	if r.Template == "" {
		r.Template = "default"
	}
	if r.BatchSize == 0 {
		r.BatchSize = 5
	}
}

// isValidDateFormat checks if date string is in YYYY-MM-DD format
func isValidDateFormat(date string) bool {
	if len(date) != 10 {
		return false
	}
	// Basic format check: YYYY-MM-DD
	if date[4] != '-' || date[7] != '-' {
		return false
	}
	return true
}

// ResultsQueryParams represents GET /api/results query parameters
type ResultsQueryParams struct {
	Project  string `query:"project"`
	Category string `query:"category"`
	FromDate string `query:"from"`
	ToDate   string `query:"to"`
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
}

// SetDefaults sets default pagination values
func (p *ResultsQueryParams) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 || p.PageSize > 100 {
		p.PageSize = 20
	}
}

// ExportQueryParams represents GET /api/export query parameters
type ExportQueryParams struct {
	Format   string `query:"format"`
	Project  string `query:"project"`
	FromDate string `query:"from"`
	ToDate   string `query:"to"`
}

// SetDefaults sets default export format
func (p *ExportQueryParams) SetDefaults() {
	if p.Format == "" {
		p.Format = "json"
	}
}
