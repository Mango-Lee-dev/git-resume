package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
	"github.com/wootaiklee/git-resume/internal/export"
	"github.com/wootaiklee/git-resume/internal/tui/state"
	"github.com/wootaiklee/git-resume/internal/tui/styles"
	"github.com/wootaiklee/git-resume/pkg/models"
)

type exportFormat struct {
	format      models.ExportFormat
	name        string
	ext         string
	description string
}

var exportFormats = []exportFormat{
	{models.FormatCSV, "CSV", "csv", "Spreadsheet format, good for analysis"},
	{models.FormatMarkdown, "Markdown", "md", "Ready to paste into resume"},
	{models.FormatJSON, "JSON", "json", "Structured data for other tools"},
}

type exportCompleteMsg struct {
	path string
	err  error
}

type ExportModel struct {
	state *state.AppState

	cursor    int
	exporting bool
	errorMsg  string

	width, height int
}

func NewExportModel(s *state.AppState) ExportModel {
	// Find current selection
	cursor := 0
	for i, f := range exportFormats {
		if f.format == s.ExportFormat {
			cursor = i
			break
		}
	}

	return ExportModel{
		state:  s,
		cursor: cursor,
	}
}

func (m ExportModel) Init() tea.Cmd {
	return nil
}

func (m ExportModel) Update(msg tea.Msg) (ExportModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.exporting {
			return m, nil
		}

		m.errorMsg = ""

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case "down", "j":
			if m.cursor < len(exportFormats)-1 {
				m.cursor++
			}
			return m, nil

		case "enter":
			m.state.ExportFormat = exportFormats[m.cursor].format
			m.exporting = true
			return m, m.doExport()

		case "esc":
			return m, Navigate(NavToResults)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case exportCompleteMsg:
		m.exporting = false
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			return m, nil
		}
		m.state.OutputPath = msg.path
		return m, Navigate(NavToDone)
	}

	return m, nil
}

func (m ExportModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(styles.TitleStyle.Render("Export Results"))
	b.WriteString("\n\n")

	// Summary
	b.WriteString(fmt.Sprintf("Exporting %d bullet points\n\n", len(m.state.Results)))

	// Format selection
	b.WriteString(styles.LabelStyle.Render("Select format:"))
	b.WriteString("\n\n")

	for i, f := range exportFormats {
		cursor := "  "
		if i == m.cursor {
			cursor = "▸ "
		}

		itemStyle := styles.NormalItemStyle
		if i == m.cursor {
			itemStyle = styles.SelectedItemStyle
		}

		// Format name with extension
		b.WriteString(cursor)
		b.WriteString(itemStyle.Render(f.name))
		b.WriteString(styles.HelpStyle.Render(fmt.Sprintf(" (.%s)", f.ext)))
		b.WriteString("\n")

		// Description
		b.WriteString("    ")
		b.WriteString(styles.HelpStyle.Render(f.description))
		b.WriteString("\n\n")
	}

	// Preview of filename
	selectedFormat := exportFormats[m.cursor]
	previewStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Subtle).
		Padding(0, 1)

	preview := fmt.Sprintf("Output: resume_YYYY-MM-DD_to_YYYY-MM-DD.%s", selectedFormat.ext)
	b.WriteString(previewStyle.Render(preview))
	b.WriteString("\n")

	// Exporting indicator
	if m.exporting {
		b.WriteString("\n")
		b.WriteString(styles.ActiveStyle.Render("Exporting..."))
		b.WriteString("\n")
	}

	// Error message
	if m.errorMsg != "" {
		b.WriteString("\n")
		b.WriteString(styles.ErrorStyle.Render("! " + m.errorMsg))
		b.WriteString("\n")
	}

	// Footer
	b.WriteString("\n")
	help := "↑/↓: select format • Enter: export • Esc: back"
	b.WriteString(styles.HelpStyle.Render(help))

	return b.String()
}

func (m ExportModel) doExport() tea.Cmd {
	return func() tea.Msg {
		outputDir := viper.GetString("output")
		if outputDir == "" {
			outputDir = viper.GetString("OUTPUT_DIR")
		}
		if outputDir == "" {
			outputDir = "./output"
		}

		exporter, err := export.NewExporter(m.state.ExportFormat, outputDir)
		if err != nil {
			return exportCompleteMsg{err: err}
		}

		var path string

		switch exp := exporter.(type) {
		case *export.CSVExporter:
			path, err = exp.ExportWithTimestamp(m.state.Results, "resume")
		case *export.MarkdownExporter:
			path, err = exp.ExportWithTimestamp(m.state.Results, "resume")
		case *export.JSONExporter:
			path, err = exp.ExportWithTimestamp(m.state.Results, "resume")
		default:
			// Fallback to Export method
			path, err = exporter.Export(m.state.Results, "resume")
		}

		if err != nil {
			return exportCompleteMsg{err: err}
		}

		return exportCompleteMsg{path: path}
	}
}

func (m ExportModel) SetSize(width, height int) ExportModel {
	m.width = width
	m.height = height
	return m
}
