package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wootaiklee/git-resume/internal/db"
	"github.com/wootaiklee/git-resume/internal/export"
	"github.com/wootaiklee/git-resume/pkg/models"
)

var (
	exportFormat string
	outputFile   string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export cached analysis results to various formats",
	Long: `Export previously analyzed results from the cache database
to various formats including CSV, Markdown, and JSON.

Example:
  git-resume export --format=markdown --output=resume.md
  git-resume export --format=json
  git-resume export --format=csv`,
	RunE: runExport,
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "csv", "output format (csv, markdown, json)")
	exportCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output filename (auto-generated if not specified)")
}

func runExport(cmd *cobra.Command, args []string) error {
	// Initialize database
	dbPath := viper.GetString("db")
	if dbPath == "" {
		dbPath = viper.GetString("DB_PATH")
	}
	if dbPath == "" {
		dbPath = "./data/cache.db"
	}

	database, err := db.New(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer database.Close()

	cache := db.NewCache(database)

	// Get all results from cache
	results, err := cache.GetAllResults()
	if err != nil {
		return fmt.Errorf("failed to get results: %w", err)
	}

	if len(results) == 0 {
		fmt.Println("No results found in cache. Run 'git-resume analyze' first.")
		return nil
	}

	fmt.Printf("Found %d cached results\n", len(results))

	// Determine output directory
	outputDir := viper.GetString("output")
	if outputDir == "" {
		outputDir = viper.GetString("OUTPUT_DIR")
	}
	if outputDir == "" {
		outputDir = "./output"
	}

	// Parse format
	format := parseFormat(exportFormat)

	// Create exporter
	exporter, err := export.NewExporter(format, outputDir)
	if err != nil {
		return fmt.Errorf("failed to create exporter: %w", err)
	}

	// Determine filename
	var outputPath string
	if outputFile != "" {
		outputPath, err = exporter.Export(results, outputFile)
	} else {
		// Use timestamp-based filename
		switch exp := exporter.(type) {
		case *export.CSVExporter:
			outputPath, err = exp.ExportWithTimestamp(results, "resume")
		case *export.MarkdownExporter:
			outputPath, err = exp.ExportWithTimestamp(results, "resume")
		case *export.JSONExporter:
			outputPath, err = exp.ExportWithTimestamp(results, "resume")
		default:
			outputPath, err = exporter.Export(results, "resume"+getExtension(format))
		}
	}

	if err != nil {
		return fmt.Errorf("failed to export: %w", err)
	}

	fmt.Printf("Exported %d results to: %s\n", len(results), outputPath)
	return nil
}

func parseFormat(s string) models.ExportFormat {
	switch strings.ToLower(s) {
	case "json":
		return models.FormatJSON
	case "markdown", "md":
		return models.FormatMarkdown
	default:
		return models.FormatCSV
	}
}

func getExtension(format models.ExportFormat) string {
	switch format {
	case models.FormatJSON:
		return ".json"
	case models.FormatMarkdown:
		return ".md"
	default:
		return ".csv"
	}
}
