package screens

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wootaiklee/git-resume/internal/llm"
	"github.com/wootaiklee/git-resume/internal/tui/state"
	"github.com/wootaiklee/git-resume/internal/tui/styles"
)

type templateItem struct {
	key    string
	config llm.TemplateConfig
}

type TemplateModel struct {
	state         *state.AppState
	templates     []templateItem
	cursor        int
	width, height int
}

func NewTemplateModel(s *state.AppState) TemplateModel {
	// Get all templates
	var templates []templateItem
	for key, config := range llm.BuiltinTemplates {
		templates = append(templates, templateItem{key: key, config: config})
	}

	// Sort by name for consistent ordering
	sort.Slice(templates, func(i, j int) bool {
		// Put "default" first
		if templates[i].key == "default" {
			return true
		}
		if templates[j].key == "default" {
			return false
		}
		return templates[i].key < templates[j].key
	})

	// Find current selection
	cursor := 0
	for i, t := range templates {
		if t.key == s.TemplateName {
			cursor = i
			break
		}
	}

	return TemplateModel{
		state:     s,
		templates: templates,
		cursor:    cursor,
	}
}

func (m TemplateModel) Init() tea.Cmd {
	return nil
}

func (m TemplateModel) Update(msg tea.Msg) (TemplateModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case "down", "j":
			if m.cursor < len(m.templates)-1 {
				m.cursor++
			}
			return m, nil

		case "enter", "ctrl+n":
			m.state.TemplateName = m.templates[m.cursor].key
			return m, Navigate(NavToOptions)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m TemplateModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(styles.TitleStyle.Render("Select Template"))
	b.WriteString("\n")
	b.WriteString(styles.SubtitleStyle.Render("Choose a template that matches your target role"))
	b.WriteString("\n\n")

	// Template list on left, details on right
	listWidth := 25

	// Build list
	var listBuilder strings.Builder
	for i, t := range m.templates {
		cursor := "  "
		if i == m.cursor {
			cursor = "▸ "
		}

		itemStyle := styles.NormalItemStyle
		if i == m.cursor {
			itemStyle = styles.SelectedItemStyle
		}

		name := t.config.Name
		if t.key == "default" {
			name += " ★"
		}

		listBuilder.WriteString(cursor)
		listBuilder.WriteString(itemStyle.Render(name))
		listBuilder.WriteString("\n")
	}

	// Build detail panel
	selected := m.templates[m.cursor]
	var detailBuilder strings.Builder

	detailBuilder.WriteString(styles.BoldStyle.Render(selected.config.Name))
	detailBuilder.WriteString("\n\n")

	detailBuilder.WriteString(styles.LabelStyle.Render("Description:"))
	detailBuilder.WriteString("\n")
	detailBuilder.WriteString(selected.config.Description)
	detailBuilder.WriteString("\n\n")

	detailBuilder.WriteString(styles.LabelStyle.Render("Tone:"))
	detailBuilder.WriteString("\n")
	detailBuilder.WriteString(selected.config.ToneStyle)
	detailBuilder.WriteString("\n\n")

	detailBuilder.WriteString(styles.LabelStyle.Render("Focus Areas:"))
	detailBuilder.WriteString("\n")
	for _, focus := range selected.config.Focus {
		detailBuilder.WriteString("• " + focus + "\n")
	}

	// Style the panels
	listPanel := lipgloss.NewStyle().
		Width(listWidth).
		Render(listBuilder.String())

	detailPanel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Primary).
		Padding(1, 2).
		Width(50).
		Render(detailBuilder.String())

	// Join horizontally
	content := lipgloss.JoinHorizontal(lipgloss.Top, listPanel, "  ", detailPanel)
	b.WriteString(content)

	// Footer
	b.WriteString("\n\n")
	help := "↑/↓: navigate • Enter: select"
	b.WriteString(styles.HelpStyle.Render(help))

	return b.String()
}

func (m TemplateModel) SetSize(width, height int) TemplateModel {
	m.width = width
	m.height = height
	return m
}
