package screens

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wootaiklee/git-resume/internal/tui/styles"
)

type WelcomeModel struct {
	width, height int
}

func NewWelcomeModel() WelcomeModel {
	return WelcomeModel{}
}

func (m WelcomeModel) Init() tea.Cmd {
	return nil
}

func (m WelcomeModel) Update(msg tea.Msg) (WelcomeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			return m, Navigate(NavToRepoSelect)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m WelcomeModel) View() string {
	var b strings.Builder

	// Logo
	logo := styles.LogoStyle().Render(styles.Logo)
	b.WriteString(logo)
	b.WriteString("\n\n")

	// Description
	desc := lipgloss.NewStyle().
		Foreground(styles.Subtle).
		Width(60).
		Align(lipgloss.Center).
		Render("Transform your Git commit history into impactful resume bullet points using Claude AI")

	b.WriteString(desc)
	b.WriteString("\n\n\n")

	// Features
	features := []string{
		"Select repositories to analyze",
		"Configure date range and template",
		"Generate STAR-format achievements",
		"Export to CSV, Markdown, or JSON",
	}

	for _, f := range features {
		b.WriteString(lipgloss.NewStyle().
			Foreground(styles.Secondary).
			Render("  ✓ "))
		b.WriteString(f)
		b.WriteString("\n")
	}

	b.WriteString("\n\n")

	// Call to action
	cta := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		Render("Press Enter to start")

	b.WriteString(cta)
	b.WriteString("\n\n")

	// Footer
	footer := styles.HelpStyle.Render("q: quit • ?: help")
	b.WriteString(footer)

	// Center the content
	content := b.String()
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m WelcomeModel) SetSize(width, height int) WelcomeModel {
	m.width = width
	m.height = height
	return m
}
