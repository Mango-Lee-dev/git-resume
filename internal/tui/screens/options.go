package screens

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wootaiklee/git-resume/internal/tui/state"
	"github.com/wootaiklee/git-resume/internal/tui/styles"
)

type OptionsModel struct {
	state *state.AppState

	batchSizeInput textinput.Model
	dryRun         bool
	notify         bool

	focusIndex    int
	errorMsg      string
	width, height int
}

const (
	optBatchSize = iota
	optDryRun
	optNotify
	optCount
)

func NewOptionsModel(s *state.AppState) OptionsModel {
	batchInput := textinput.New()
	batchInput.Placeholder = "5"
	batchInput.SetValue(fmt.Sprintf("%d", s.BatchSize))
	batchInput.CharLimit = 2
	batchInput.Width = 5

	return OptionsModel{
		state:          s,
		batchSizeInput: batchInput,
		dryRun:         s.DryRun,
		notify:         s.Notify,
		focusIndex:     0,
	}
}

func (m OptionsModel) Init() tea.Cmd {
	m.batchSizeInput.Focus()
	return textinput.Blink
}

func (m OptionsModel) Update(msg tea.Msg) (OptionsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.errorMsg = ""

		switch msg.String() {
		case "tab", "down", "j":
			m.focusIndex = (m.focusIndex + 1) % optCount
			m.updateFocus()
			return m, nil

		case "shift+tab", "up", "k":
			m.focusIndex = (m.focusIndex - 1 + optCount) % optCount
			m.updateFocus()
			return m, nil

		case " ", "x":
			// Toggle for checkbox options
			switch m.focusIndex {
			case optDryRun:
				m.dryRun = !m.dryRun
			case optNotify:
				m.notify = !m.notify
			}
			return m, nil

		case "enter", "ctrl+n":
			if m.validate() {
				m.saveToState()
				return m, Navigate(NavToAnalysis)
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Update text input if focused
	if m.focusIndex == optBatchSize {
		var cmd tea.Cmd
		m.batchSizeInput, cmd = m.batchSizeInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m OptionsModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(styles.TitleStyle.Render("Additional Options"))
	b.WriteString("\n\n")

	// Summary of selections so far
	b.WriteString(styles.LabelStyle.Render("Summary:"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  Repositories: %d selected\n", len(m.state.RepoPaths)))

	from, to := m.state.GetDateRange()
	b.WriteString(fmt.Sprintf("  Date Range: %s to %s\n", from.Format("Jan 2, 2006"), to.Format("Jan 2, 2006")))
	b.WriteString(fmt.Sprintf("  Template: %s\n", m.state.TemplateName))
	b.WriteString("\n")

	// Options
	b.WriteString(styles.LabelStyle.Render("Options:"))
	b.WriteString("\n\n")

	// Batch size
	batchLabel := "  Batch Size: "
	if m.focusIndex == optBatchSize {
		batchLabel = "▸ Batch Size: "
		b.WriteString(styles.FocusedStyle.Render(batchLabel))
	} else {
		b.WriteString(batchLabel)
	}
	b.WriteString(m.batchSizeInput.View())
	b.WriteString(styles.HelpStyle.Render(" (commits per API call)"))
	b.WriteString("\n\n")

	// Dry run
	dryRunCheck := "[ ]"
	if m.dryRun {
		dryRunCheck = "[x]"
	}
	dryRunLabel := fmt.Sprintf("  %s Dry Run", dryRunCheck)
	if m.focusIndex == optDryRun {
		dryRunLabel = fmt.Sprintf("▸ %s Dry Run", dryRunCheck)
		b.WriteString(styles.FocusedStyle.Render(dryRunLabel))
	} else {
		b.WriteString(dryRunLabel)
	}
	b.WriteString(styles.HelpStyle.Render(" (preview without API calls)"))
	b.WriteString("\n\n")

	// Notify
	notifyCheck := "[ ]"
	if m.notify {
		notifyCheck = "[x]"
	}
	notifyLabel := fmt.Sprintf("  %s Slack Notify", notifyCheck)
	if m.focusIndex == optNotify {
		notifyLabel = fmt.Sprintf("▸ %s Slack Notify", notifyCheck)
		b.WriteString(styles.FocusedStyle.Render(notifyLabel))
	} else {
		b.WriteString(notifyLabel)
	}
	b.WriteString(styles.HelpStyle.Render(" (send notification when complete)"))
	b.WriteString("\n")

	// Error message
	if m.errorMsg != "" {
		b.WriteString("\n")
		b.WriteString(styles.ErrorStyle.Render("! " + m.errorMsg))
	}

	// Footer
	b.WriteString("\n\n")
	help := "Tab: next • Space: toggle • Enter: start analysis"
	b.WriteString(styles.HelpStyle.Render(help))

	return b.String()
}

func (m *OptionsModel) updateFocus() {
	m.batchSizeInput.Blur()
	if m.focusIndex == optBatchSize {
		m.batchSizeInput.Focus()
	}
}

func (m *OptionsModel) validate() bool {
	batchSize, err := strconv.Atoi(m.batchSizeInput.Value())
	if err != nil || batchSize < 1 || batchSize > 20 {
		m.errorMsg = "Batch size must be between 1 and 20"
		return false
	}
	return true
}

func (m *OptionsModel) saveToState() {
	m.state.BatchSize, _ = strconv.Atoi(m.batchSizeInput.Value())
	m.state.DryRun = m.dryRun
	m.state.Notify = m.notify
}

func (m OptionsModel) SetSize(width, height int) OptionsModel {
	m.width = width
	m.height = height
	return m
}
