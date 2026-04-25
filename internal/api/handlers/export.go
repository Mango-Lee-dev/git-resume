package handlers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/wootaiklee/git-resume/internal/api/dto"
	"github.com/wootaiklee/git-resume/internal/service"
	"github.com/wootaiklee/git-resume/pkg/models"
)

// ExportHandler handles export endpoints
type ExportHandler struct {
	results *service.ResultsService
}

// NewExportHandler creates a new export handler
func NewExportHandler(results *service.ResultsService) *ExportHandler {
	return &ExportHandler{results: results}
}

// Export handles GET /api/export
func (h *ExportHandler) Export(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	opts := parseExportQuery(r)
	results, err := h.results.GetByFilters(r.Context(), opts)
	if err != nil {
		respondError(w, r, dto.NewInternalError("failed to get results"))
		return
	}

	switch format {
	case "json":
		h.exportJSON(w, results, opts)
	case "csv":
		h.exportCSV(w, results)
	case "markdown", "md":
		h.exportMarkdown(w, results)
	default:
		respondError(w, r, dto.NewBadRequestError("unsupported format: "+format))
	}
}

func (h *ExportHandler) exportJSON(w http.ResponseWriter, results []models.AnalysisResult, opts service.QueryOptions) {
	// Build response with metadata
	type jsonExport struct {
		Metadata     dto.ExportMetadata   `json:"metadata"`
		Achievements []dto.ResultResponse `json:"achievements"`
		Summary      struct {
			ByCategory map[string]int `json:"by_category"`
			ByProject  map[string]int `json:"by_project"`
		} `json:"summary"`
	}

	export := jsonExport{
		Metadata: dto.ExportMetadata{
			GeneratedAt: time.Now(),
			TotalCount:  len(results),
			Format:      "json",
		},
		Achievements: make([]dto.ResultResponse, len(results)),
	}

	if !opts.FromDate.IsZero() {
		export.Metadata.FromDate = opts.FromDate.Format("2006-01-02")
	}
	if !opts.ToDate.IsZero() {
		export.Metadata.ToDate = opts.ToDate.Format("2006-01-02")
	}

	export.Summary.ByCategory = make(map[string]int)
	export.Summary.ByProject = make(map[string]int)

	for i, result := range results {
		export.Achievements[i] = dto.FromAnalysisResult(result)
		export.Summary.ByCategory[string(result.Category)]++
		export.Summary.ByProject[result.Project]++
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=resume-export.json")
	json.NewEncoder(w).Encode(export)
}

func (h *ExportHandler) exportCSV(w http.ResponseWriter, results []models.AnalysisResult) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=resume-export.csv")

	// Write UTF-8 BOM for Excel compatibility
	w.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Header
	writer.Write([]string{"Date", "Project", "Category", "Impact Summary", "Commit Hash"})

	// Data rows
	for _, r := range results {
		writer.Write([]string{
			r.Date.Format("2006-01-02"),
			r.Project,
			string(r.Category),
			r.ImpactSummary,
			r.CommitHash[:7],
		})
	}
}

func (h *ExportHandler) exportMarkdown(w http.ResponseWriter, results []models.AnalysisResult) {
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=resume-export.md")

	var buf bytes.Buffer

	buf.WriteString("# Resume Achievements\n\n")
	buf.WriteString(fmt.Sprintf("*Generated on %s*\n\n", time.Now().Format("January 2, 2006")))

	// Group by category
	categoryEmoji := map[string]string{
		"Feature":  "✨",
		"Fix":      "🐛",
		"Refactor": "♻️",
		"Test":     "🧪",
		"Docs":     "📚",
		"Chore":    "🔧",
	}

	byCategory := make(map[string][]models.AnalysisResult)
	for _, r := range results {
		cat := string(r.Category)
		byCategory[cat] = append(byCategory[cat], r)
	}

	for cat, items := range byCategory {
		emoji := categoryEmoji[cat]
		if emoji == "" {
			emoji = "📌"
		}
		buf.WriteString(fmt.Sprintf("## %s %s\n\n", emoji, cat))

		for _, r := range items {
			buf.WriteString(fmt.Sprintf("- %s (%s, %s) `%s`\n",
				r.ImpactSummary,
				r.Project,
				r.Date.Format("2006-01-02"),
				r.CommitHash[:7],
			))
		}
		buf.WriteString("\n")
	}

	w.Write(buf.Bytes())
}

func parseExportQuery(r *http.Request) service.QueryOptions {
	opts := service.QueryOptions{
		Project: r.URL.Query().Get("project"),
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

	return opts
}
