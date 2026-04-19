package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wootaiklee/git-resume/internal/tui/state"
	"github.com/wootaiklee/git-resume/internal/tui/styles"
)

type ResultsModel struct {
	state *state.AppState

	cursor        int
	offset        int
	maxVisible    int
	width, height int
}

func NewResultsModel(s *state.AppState) ResultsModel {
	return ResultsModel{
		state:      s,
		cursor:     0,
		offset:     0,
		maxVisible: 8,
	}
}

func (m ResultsModel) Init() tea.Cmd {
	return nil
}

func (m ResultsModel) Update(msg tea.Msg) (ResultsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.offset {
					m.offset = m.cursor
				}
			}
			return m, nil

		case "down", "j":
			if m.cursor < len(m.state.Results)-1 {
				m.cursor++
				if m.cursor >= m.offset+m.maxVisible {
					m.offset = m.cursor - m.maxVisible + 1
				}
			}
			return m, nil

		case "enter", "ctrl+n":
			return m, Navigate(NavToExport)

		case "esc":
			// Go back to options (but note: re-running analysis would require resetting state)
			return m, Navigate(NavToOptions)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.maxVisible = (msg.Height - 15) / 3
		if m.maxVisible < 3 {
			m.maxVisible = 3
		}
	}

	return m, nil
}

func (m ResultsModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(styles.TitleStyle.Render("Analysis Results"))
	b.WriteString("\n\n")

	// Summary
	summaryStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Secondary).
		Padding(0, 2)

	summary := fmt.Sprintf(
		"Generated %s bullet points | Tokens used: %s",
		styles.BoldStyle.Render(fmt.Sprintf("%d", len(m.state.Results))),
		styles.BoldStyle.Render(fmt.Sprintf("%d", m.state.TokensUsed)),
	)
	b.WriteString(summaryStyle.Render(summary))
	b.WriteString("\n\n")

	if len(m.state.Results) == 0 {
		b.WriteString(styles.HelpStyle.Render("No results to display"))
		b.WriteString("\n\n")
		b.WriteString(styles.HelpStyle.Render("Esc: go back"))
		return b.String()
	}

	// Results list
	b.WriteString(styles.LabelStyle.Render("Results:"))
	b.WriteString("\n\n")

	endIdx := m.offset + m.maxVisible
	if endIdx > len(m.state.Results) {
		endIdx = len(m.state.Results)
	}

	for i := m.offset; i < endIdx; i++ {
		r := m.state.Results[i]
		isSelected := i == m.cursor

		// Category badge
		categoryStyle := lipgloss.NewStyle().
			Padding(0, 1).
			Background(getCategoryColor(string(r.Category))).
			Foreground(lipgloss.Color("#FFFFFF"))

		// Cursor
		cursor := "  "
		if isSelected {
			cursor = "▸ "
		}

		// Item style
		itemStyle := lipgloss.NewStyle()
		if isSelected {
			itemStyle = itemStyle.Bold(true)
		}

		// Build item
		b.WriteString(cursor)
		b.WriteString(categoryStyle.Render(string(r.Category)))
		b.WriteString(" ")
		b.WriteString(styles.HelpStyle.Render(r.CommitHash[:7]))
		b.WriteString("\n")

		// Summary (truncate if too long)
		summaryText := r.ImpactSummary
		maxLen := m.width - 10
		if maxLen > 80 {
			maxLen = 80
		}
		if len(summaryText) > maxLen {
			summaryText = summaryText[:maxLen-3] + "..."
		}

		b.WriteString("   ")
		b.WriteString(itemStyle.Render(summaryText))
		b.WriteString("\n")

		// Project and date
		b.WriteString("   ")
		b.WriteString(styles.HelpStyle.Render(fmt.Sprintf("%s • %s", r.Project, r.Date.Format("Jan 2, 2006"))))
		b.WriteString("\n\n")
	}

	// Scroll indicator
	if len(m.state.Results) > m.maxVisible {
		scrollInfo := fmt.Sprintf("Showing %d-%d of %d", m.offset+1, endIdx, len(m.state.Results))
		b.WriteString(styles.HelpStyle.Render(scrollInfo))
		b.WriteString("\n")
	}

	// Footer
	b.WriteString("\n")
	help := "↑/↓: navigate • Enter: export • Esc: back"
	b.WriteString(styles.HelpStyle.Render(help))

	return b.String()
}

func getCategoryColor(category string) lipgloss.Color {
	switch category {
	case "Feature":
		return lipgloss.Color("#10B981") // Green
	case "Fix":
		return lipgloss.Color("#EF4444") // Red
	case "Refactor":
		return lipgloss.Color("#8B5CF6") // Purple
	case "Test":
		return lipgloss.Color("#F59E0B") // Amber
	case "Docs":
		return lipgloss.Color("#3B82F6") // Blue
	case "Chore":
		return lipgloss.Color("#6B7280") // Gray
	default:
		return lipgloss.Color("#6B7280")
	}
}

func (m ResultsModel) SetSize(width, height int) ResultsModel {
	m.width = width
	m.height = height
	m.maxVisible = (height - 15) / 3
	if m.maxVisible < 3 {
		m.maxVisible = 3
	}
	return m
}
