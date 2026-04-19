package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wootaiklee/git-resume/internal/tui/state"
	"github.com/wootaiklee/git-resume/internal/tui/styles"
)

type DoneModel struct {
	state         *state.AppState
	width, height int
}

func NewDoneModel(s *state.AppState) DoneModel {
	return DoneModel{
		state: s,
	}
}

func (m DoneModel) Init() tea.Cmd {
	return nil
}

func (m DoneModel) Update(msg tea.Msg) (DoneModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "q", "esc":
			return m, tea.Quit
		case "r":
			// Restart - go back to repo selection
			m.state.Reset()
			return m, Navigate(NavToRepoSelect)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m DoneModel) View() string {
	var b strings.Builder

	// Success icon
	successBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Secondary).
		Padding(1, 3).
		Align(lipgloss.Center)

	b.WriteString(successBox.Render(styles.SuccessStyle.Render("✓ Export Complete!")))
	b.WriteString("\n\n")

	// File path
	b.WriteString(styles.LabelStyle.Render("File saved to:"))
	b.WriteString("\n")
	b.WriteString(styles.BoldStyle.Render(m.state.OutputPath))
	b.WriteString("\n\n")

	// Summary
	b.WriteString(styles.LabelStyle.Render("Summary:"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  • Bullet points: %d\n", len(m.state.Results)))
	b.WriteString(fmt.Sprintf("  • Tokens used: %d\n", m.state.TokensUsed))
	b.WriteString(fmt.Sprintf("  • Repositories: %d\n", len(m.state.RepoPaths)))
	b.WriteString(fmt.Sprintf("  • Template: %s\n", m.state.TemplateName))
	b.WriteString(fmt.Sprintf("  • Format: %s\n", m.state.ExportFormat))
	b.WriteString("\n")

	// Category breakdown
	categories := make(map[string]int)
	for _, r := range m.state.Results {
		categories[string(r.Category)]++
	}

	if len(categories) > 0 {
		b.WriteString(styles.LabelStyle.Render("By category:"))
		b.WriteString("\n")
		for cat, count := range categories {
			b.WriteString(fmt.Sprintf("  • %s: %d\n", cat, count))
		}
		b.WriteString("\n")
	}

	// Next steps
	nextStepsStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.Subtle).
		Padding(1, 2)

	nextSteps := `Next steps:
1. Open the file in your preferred editor
2. Review and customize the bullet points
3. Copy relevant achievements to your resume`

	b.WriteString(nextStepsStyle.Render(nextSteps))
	b.WriteString("\n\n")

	// Footer
	help := "Enter/q: quit • r: start over"
	b.WriteString(styles.HelpStyle.Render(help))

	// Center content
	content := b.String()
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m DoneModel) SetSize(width, height int) DoneModel {
	m.width = width
	m.height = height
	return m
}
