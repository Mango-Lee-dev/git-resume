package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wootaiklee/git-resume/pkg/models"
)

// CSVExporter exports analysis results to CSV format
type CSVExporter struct {
	outputDir string
}

// NewCSVExporter creates a new CSV exporter
func NewCSVExporter(outputDir string) *CSVExporter {
	return &CSVExporter{outputDir: outputDir}
}

// Export writes results to a CSV file
func (e *CSVExporter) Export(results []models.AnalysisResult, filename string) (string, error) {
	if err := os.MkdirAll(e.outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	outputPath := filepath.Join(e.outputDir, filename)
	file, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write UTF-8 BOM for Excel compatibility
	file.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Date", "Project", "Category", "Impact_Summary", "Commit_Hash"}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write header: %w", err)
	}

	// Write data rows
	for _, result := range results {
		row := []string{
			result.Date.Format("2006-01-02"),
			result.Project,
			string(result.Category),
			result.ImpactSummary,
			result.CommitHash[:7],
		}
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write row: %w", err)
		}
	}

	return outputPath, nil
}

// ExportWithTimestamp exports with automatic timestamp in filename
func (e *CSVExporter) ExportWithTimestamp(results []models.AnalysisResult, prefix string) (string, error) {
	if len(results) == 0 {
		return "", fmt.Errorf("no results to export")
	}

	minDate, maxDate := getDateRange(results)
	filename := fmt.Sprintf("%s_%s_to_%s.csv",
		prefix,
		minDate.Format("2006-01-02"),
		maxDate.Format("2006-01-02"),
	)

	return e.Export(results, filename)
}
