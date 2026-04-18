package export

import (
	"fmt"

	"github.com/wootaiklee/git-resume/pkg/models"
)

// Exporter interface for different export formats
type Exporter interface {
	Export(results []models.AnalysisResult, filename string) (string, error)
}

// NewExporter creates an exporter for the specified format
func NewExporter(format models.ExportFormat, outputDir string) (Exporter, error) {
	switch format {
	case models.FormatCSV:
		return NewCSVExporter(outputDir), nil
	case models.FormatJSON:
		return NewJSONExporter(outputDir), nil
	case models.FormatMarkdown:
		return NewMarkdownExporter(outputDir), nil
	default:
		return nil, fmt.Errorf("unknown format: %s", format)
	}
}
