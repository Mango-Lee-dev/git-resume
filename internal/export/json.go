package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/wootaiklee/git-resume/pkg/models"
)

// JSONExporter exports analysis results to JSON format
type JSONExporter struct {
	outputDir string
}

// NewJSONExporter creates a new JSON exporter
func NewJSONExporter(outputDir string) *JSONExporter {
	return &JSONExporter{outputDir: outputDir}
}

// JSONOutput represents the full JSON export structure
type JSONOutput struct {
	Metadata    JSONMetadata         `json:"metadata"`
	Achievements []JSONAchievement   `json:"achievements"`
	Summary     JSONSummary          `json:"summary"`
}

// JSONMetadata contains export metadata
type JSONMetadata struct {
	GeneratedAt string `json:"generated_at"`
	FromDate    string `json:"from_date"`
	ToDate      string `json:"to_date"`
	TotalCount  int    `json:"total_count"`
	Version     string `json:"version"`
}

// JSONAchievement represents a single achievement entry
type JSONAchievement struct {
	Date        string `json:"date"`
	Project     string `json:"project"`
	Category    string `json:"category"`
	Summary     string `json:"summary"`
	CommitHash  string `json:"commit_hash"`
}

// JSONSummary contains category counts
type JSONSummary struct {
	ByCategory map[string]int `json:"by_category"`
	ByProject  map[string]int `json:"by_project"`
}

// Export writes results to a JSON file
func (e *JSONExporter) Export(results []models.AnalysisResult, filename string) (string, error) {
	if err := os.MkdirAll(e.outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	outputPath := filepath.Join(e.outputDir, filename)
	file, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Build output structure
	output := e.buildOutput(results)

	// Encode with indentation for readability
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(output); err != nil {
		return "", fmt.Errorf("failed to encode JSON: %w", err)
	}

	return outputPath, nil
}

// ExportWithTimestamp exports with automatic timestamp in filename
func (e *JSONExporter) ExportWithTimestamp(results []models.AnalysisResult, prefix string) (string, error) {
	if len(results) == 0 {
		return "", fmt.Errorf("no results to export")
	}

	minDate, maxDate := getDateRange(results)
	filename := fmt.Sprintf("%s_%s_to_%s.json",
		prefix,
		minDate.Format("2006-01-02"),
		maxDate.Format("2006-01-02"),
	)

	return e.Export(results, filename)
}

func (e *JSONExporter) buildOutput(results []models.AnalysisResult) JSONOutput {
	output := JSONOutput{
		Metadata: JSONMetadata{
			GeneratedAt: time.Now().UTC().Format(time.RFC3339),
			TotalCount:  len(results),
			Version:     "1.0.0",
		},
		Achievements: make([]JSONAchievement, 0, len(results)),
		Summary: JSONSummary{
			ByCategory: make(map[string]int),
			ByProject:  make(map[string]int),
		},
	}

	if len(results) > 0 {
		minDate, maxDate := getDateRange(results)
		output.Metadata.FromDate = minDate.Format("2006-01-02")
		output.Metadata.ToDate = maxDate.Format("2006-01-02")
	}

	for _, r := range results {
		// Add achievement
		output.Achievements = append(output.Achievements, JSONAchievement{
			Date:       r.Date.Format("2006-01-02"),
			Project:    r.Project,
			Category:   string(r.Category),
			Summary:    r.ImpactSummary,
			CommitHash: r.CommitHash[:7],
		})

		// Update summary counts
		output.Summary.ByCategory[string(r.Category)]++
		if r.Project != "" {
			output.Summary.ByProject[r.Project]++
		}
	}

	return output
}
