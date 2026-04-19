package screens

import tea "github.com/charmbracelet/bubbletea"

// NavigateMsg is a simple string-based navigation message
type NavigateMsg string

const (
	NavToWelcome    NavigateMsg = "welcome"
	NavToRepoSelect NavigateMsg = "repo_select"
	NavToDateRange  NavigateMsg = "date_range"
	NavToTemplate   NavigateMsg = "template"
	NavToOptions    NavigateMsg = "options"
	NavToAnalysis   NavigateMsg = "analysis"
	NavToResults    NavigateMsg = "results"
	NavToExport     NavigateMsg = "export"
	NavToDone       NavigateMsg = "done"
)

// Navigate creates a command to navigate to a screen
func Navigate(dest NavigateMsg) tea.Cmd {
	return func() tea.Msg {
		return dest
	}
}

// Helper function to truncate strings
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
