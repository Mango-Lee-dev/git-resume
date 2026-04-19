package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wootaiklee/git-resume/internal/tui/screens"
	"github.com/wootaiklee/git-resume/internal/tui/state"
	"github.com/wootaiklee/git-resume/internal/tui/styles"
)

// Screen represents the current active screen
type Screen int

const (
	ScreenWelcome Screen = iota
	ScreenRepoSelect
	ScreenDateRange
	ScreenTemplate
	ScreenOptions
	ScreenAnalysis
	ScreenResults
	ScreenExport
	ScreenDone
)

// ScreenTransitionMsg is sent when we need to change screens
type ScreenTransitionMsg struct {
	NextScreen Screen
}

// Model is the main Bubble Tea application model
type Model struct {
	screen Screen
	state  *state.AppState

	// Screen models
	welcome    screens.WelcomeModel
	repoSelect screens.RepoSelectModel
	dateRange  screens.DateRangeModel
	template   screens.TemplateModel
	options    screens.OptionsModel
	analysis   screens.AnalysisModel
	results    screens.ResultsModel
	export     screens.ExportModel
	done       screens.DoneModel

	width, height int
	quitting      bool
}

// KeyMap defines key bindings
type KeyMap struct {
	Quit key.Binding
	Help key.Binding
	Back key.Binding
}

var keys = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
}

// New creates a new TUI application
func New() Model {
	appState := state.NewAppState()
	return Model{
		screen:     ScreenWelcome,
		state:      appState,
		welcome:    screens.NewWelcomeModel(),
		repoSelect: screens.NewRepoSelectModel(appState),
		dateRange:  screens.NewDateRangeModel(appState),
		template:   screens.NewTemplateModel(appState),
		options:    screens.NewOptionsModel(appState),
		analysis:   screens.NewAnalysisModel(appState),
		results:    screens.NewResultsModel(appState),
		export:     screens.NewExportModel(appState),
		done:       screens.NewDoneModel(appState),
	}
}

// State returns the application state
func (m Model) State() *state.AppState {
	return m.state
}

// Init initializes the application
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.welcome.Init(),
	)
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global key handling
		switch {
		case key.Matches(msg, keys.Quit):
			// Allow quit only when not in analysis
			if m.screen != ScreenAnalysis {
				m.quitting = true
				return m, tea.Quit
			}
		case key.Matches(msg, keys.Back):
			if m.screen != ScreenWelcome && m.screen != ScreenAnalysis && m.screen != ScreenDone {
				return m.navigateBack()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Propagate to all screens
		m.repoSelect = m.repoSelect.SetSize(msg.Width, msg.Height)
		m.dateRange = m.dateRange.SetSize(msg.Width, msg.Height)
		m.template = m.template.SetSize(msg.Width, msg.Height)
		m.options = m.options.SetSize(msg.Width, msg.Height)
		m.analysis = m.analysis.SetSize(msg.Width, msg.Height)
		m.results = m.results.SetSize(msg.Width, msg.Height)
		m.export = m.export.SetSize(msg.Width, msg.Height)

	case ScreenTransitionMsg:
		return m.transitionTo(msg.NextScreen)

	case screens.NavigateMsg:
		return m.handleNavigation(msg)
	}

	// Delegate to current screen
	return m.updateCurrentScreen(msg)
}

// View renders the current view
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	content := m.viewCurrentScreen()

	// Add padding
	return lipgloss.NewStyle().Padding(1, 2).Render(content)
}

func (m Model) updateCurrentScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.screen {
	case ScreenWelcome:
		m.welcome, cmd = m.welcome.Update(msg)
	case ScreenRepoSelect:
		m.repoSelect, cmd = m.repoSelect.Update(msg)
	case ScreenDateRange:
		m.dateRange, cmd = m.dateRange.Update(msg)
	case ScreenTemplate:
		m.template, cmd = m.template.Update(msg)
	case ScreenOptions:
		m.options, cmd = m.options.Update(msg)
	case ScreenAnalysis:
		m.analysis, cmd = m.analysis.Update(msg)
	case ScreenResults:
		m.results, cmd = m.results.Update(msg)
	case ScreenExport:
		m.export, cmd = m.export.Update(msg)
	case ScreenDone:
		m.done, cmd = m.done.Update(msg)
	}

	return m, cmd
}

func (m Model) viewCurrentScreen() string {
	switch m.screen {
	case ScreenWelcome:
		return m.welcome.View()
	case ScreenRepoSelect:
		return m.repoSelect.View()
	case ScreenDateRange:
		return m.dateRange.View()
	case ScreenTemplate:
		return m.template.View()
	case ScreenOptions:
		return m.options.View()
	case ScreenAnalysis:
		return m.analysis.View()
	case ScreenResults:
		return m.results.View()
	case ScreenExport:
		return m.export.View()
	case ScreenDone:
		return m.done.View()
	default:
		return "Unknown screen"
	}
}

func (m Model) transitionTo(next Screen) (tea.Model, tea.Cmd) {
	m.screen = next

	// Initialize the new screen if needed
	var cmd tea.Cmd
	switch next {
	case ScreenRepoSelect:
		m.repoSelect = screens.NewRepoSelectModel(m.state)
		cmd = m.repoSelect.Init()
	case ScreenDateRange:
		m.dateRange = screens.NewDateRangeModel(m.state)
		cmd = m.dateRange.Init()
	case ScreenTemplate:
		m.template = screens.NewTemplateModel(m.state)
		cmd = m.template.Init()
	case ScreenOptions:
		m.options = screens.NewOptionsModel(m.state)
		cmd = m.options.Init()
	case ScreenAnalysis:
		m.analysis = screens.NewAnalysisModel(m.state)
		cmd = m.analysis.Init()
	case ScreenResults:
		m.results = screens.NewResultsModel(m.state)
		cmd = m.results.Init()
	case ScreenExport:
		m.export = screens.NewExportModel(m.state)
		cmd = m.export.Init()
	case ScreenDone:
		m.done = screens.NewDoneModel(m.state)
		cmd = m.done.Init()
	}

	return m, cmd
}

func (m Model) handleNavigation(nav screens.NavigateMsg) (tea.Model, tea.Cmd) {
	var screen Screen
	switch nav {
	case screens.NavToWelcome:
		screen = ScreenWelcome
	case screens.NavToRepoSelect:
		screen = ScreenRepoSelect
	case screens.NavToDateRange:
		screen = ScreenDateRange
	case screens.NavToTemplate:
		screen = ScreenTemplate
	case screens.NavToOptions:
		screen = ScreenOptions
	case screens.NavToAnalysis:
		screen = ScreenAnalysis
	case screens.NavToResults:
		screen = ScreenResults
	case screens.NavToExport:
		screen = ScreenExport
	case screens.NavToDone:
		screen = ScreenDone
	default:
		return m, nil
	}
	return m.transitionTo(screen)
}

func (m Model) navigateBack() (tea.Model, tea.Cmd) {
	var prevScreen Screen

	switch m.screen {
	case ScreenRepoSelect:
		prevScreen = ScreenWelcome
	case ScreenDateRange:
		prevScreen = ScreenRepoSelect
	case ScreenTemplate:
		prevScreen = ScreenDateRange
	case ScreenOptions:
		prevScreen = ScreenTemplate
	case ScreenResults:
		prevScreen = ScreenOptions
	case ScreenExport:
		prevScreen = ScreenResults
	default:
		return m, nil
	}

	return m.transitionTo(prevScreen)
}

// Helper to create transition command
func TransitionTo(screen Screen) tea.Cmd {
	return func() tea.Msg {
		return ScreenTransitionMsg{NextScreen: screen}
	}
}

// Screen constants for external use
const (
	ToWelcome    = ScreenWelcome
	ToRepoSelect = ScreenRepoSelect
	ToDateRange  = ScreenDateRange
	ToTemplate   = ScreenTemplate
	ToOptions    = ScreenOptions
	ToAnalysis   = ScreenAnalysis
	ToResults    = ScreenResults
	ToExport     = ScreenExport
	ToDone       = ScreenDone
)

// Help text for footer
func (m Model) helpView() string {
	helpText := "esc: back • q: quit"
	if m.screen == ScreenAnalysis {
		helpText = "ctrl+c: cancel"
	}
	return styles.HelpStyle.Render(helpText)
}
