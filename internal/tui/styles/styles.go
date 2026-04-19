package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	Primary   = lipgloss.Color("#7C3AED") // Purple
	Secondary = lipgloss.Color("#10B981") // Green
	Subtle    = lipgloss.Color("#6B7280") // Gray
	Error     = lipgloss.Color("#EF4444") // Red
	Warning   = lipgloss.Color("#F59E0B") // Amber

	// Base styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Subtle).
			MarginBottom(1)

	BoldStyle = lipgloss.NewStyle().Bold(true)

	LabelStyle = lipgloss.NewStyle().
			Foreground(Subtle).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(Subtle).
			Italic(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Bold(true)

	// Component styles
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	DetailBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(Subtle).
			Padding(1, 2)

	// Input styles
	FocusedStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	BlurredStyle = lipgloss.NewStyle().
			Foreground(Subtle)

	// List item styles
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Primary).
				Bold(true).
				PaddingLeft(2)

	NormalItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			PaddingLeft(2)

	// Status styles
	ActiveStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Bold(true)

	InactiveStyle = lipgloss.NewStyle().
			Foreground(Subtle)

	// Footer
	FooterStyle = lipgloss.NewStyle().
			Foreground(Subtle).
			MarginTop(1).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(Subtle).
			PaddingTop(1)
)

// Logo for welcome screen
const Logo = `
   _____ _ _     ____
  / ____(_) |   |  _ \
 | |  __ _| |_  | |_) | ___  ___ _   _ _ __ ___  ___
 | | |_ | | __| |  _ < / _ \/ __| | | | '_ ` + "`" + ` _ \/ _ \
 | |__| | | |_  | |_) |  __/\__ \ |_| | | | | | |  __/
  \_____|_|\__| |____/ \___||___/\__,_|_| |_| |_|\___|
`

func LogoStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(Primary).Bold(true)
}
