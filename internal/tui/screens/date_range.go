package screens

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wootaiklee/git-resume/internal/tui/state"
	"github.com/wootaiklee/git-resume/internal/tui/styles"
)

type DateRangeModel struct {
	state *state.AppState

	mode state.DateRangeMode

	// Month/Year mode inputs
	monthInput textinput.Model
	yearInput  textinput.Model

	// Custom date mode inputs
	fromInput textinput.Model
	toInput   textinput.Model

	focusIndex    int
	errorMsg      string
	width, height int
}

func NewDateRangeModel(s *state.AppState) DateRangeModel {
	now := time.Now()

	monthInput := textinput.New()
	monthInput.Placeholder = "4"
	monthInput.SetValue(fmt.Sprintf("%d", s.Month))
	monthInput.CharLimit = 2
	monthInput.Width = 5

	yearInput := textinput.New()
	yearInput.Placeholder = "2024"
	yearInput.SetValue(fmt.Sprintf("%d", s.Year))
	yearInput.CharLimit = 4
	yearInput.Width = 8

	fromInput := textinput.New()
	fromInput.Placeholder = "YYYY-MM-DD"
	if !s.FromDate.IsZero() {
		fromInput.SetValue(s.FromDate.Format("2006-01-02"))
	} else {
		// Default to first day of current month
		fromInput.SetValue(time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"))
	}
	fromInput.CharLimit = 10
	fromInput.Width = 15

	toInput := textinput.New()
	toInput.Placeholder = "YYYY-MM-DD"
	if !s.ToDate.IsZero() {
		toInput.SetValue(s.ToDate.Format("2006-01-02"))
	} else {
		toInput.SetValue(now.Format("2006-01-02"))
	}
	toInput.CharLimit = 10
	toInput.Width = 15

	m := DateRangeModel{
		state:      s,
		mode:       s.DateMode,
		monthInput: monthInput,
		yearInput:  yearInput,
		fromInput:  fromInput,
		toInput:    toInput,
		focusIndex: 0,
	}

	m.updateFocus()
	return m
}

func (m DateRangeModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m DateRangeModel) Update(msg tea.Msg) (DateRangeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.errorMsg = ""

		switch msg.String() {
		case "1":
			m.mode = state.DateRangeMonth
			m.focusIndex = 0
			m.updateFocus()
			return m, nil

		case "2":
			m.mode = state.DateRangeCustom
			m.focusIndex = 0
			m.updateFocus()
			return m, nil

		case "tab", "down", "j":
			m.focusIndex++
			maxIdx := 1
			if m.mode == state.DateRangeCustom {
				maxIdx = 1
			}
			if m.focusIndex > maxIdx {
				m.focusIndex = 0
			}
			m.updateFocus()
			return m, nil

		case "shift+tab", "up", "k":
			m.focusIndex--
			maxIdx := 1
			if m.mode == state.DateRangeCustom {
				maxIdx = 1
			}
			if m.focusIndex < 0 {
				m.focusIndex = maxIdx
			}
			m.updateFocus()
			return m, nil

		case "enter", "ctrl+n":
			if m.validate() {
				m.saveToState()
				return m, Navigate(NavToTemplate)
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Update the focused input
	var cmd tea.Cmd
	if m.mode == state.DateRangeMonth {
		if m.focusIndex == 0 {
			m.monthInput, cmd = m.monthInput.Update(msg)
		} else {
			m.yearInput, cmd = m.yearInput.Update(msg)
		}
	} else {
		if m.focusIndex == 0 {
			m.fromInput, cmd = m.fromInput.Update(msg)
		} else {
			m.toInput, cmd = m.toInput.Update(msg)
		}
	}

	return m, cmd
}

func (m DateRangeModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(styles.TitleStyle.Render("Select Date Range"))
	b.WriteString("\n\n")

	// Mode selector
	monthStyle := styles.BlurredStyle
	customStyle := styles.BlurredStyle
	if m.mode == state.DateRangeMonth {
		monthStyle = styles.FocusedStyle
	} else {
		customStyle = styles.FocusedStyle
	}

	b.WriteString(monthStyle.Render("[1] Month/Year"))
	b.WriteString("    ")
	b.WriteString(customStyle.Render("[2] Custom Range"))
	b.WriteString("\n\n")

	if m.mode == state.DateRangeMonth {
		// Month/Year mode
		b.WriteString(m.renderField("Month", m.monthInput, m.focusIndex == 0))
		b.WriteString("  ")
		b.WriteString(m.renderField("Year", m.yearInput, m.focusIndex == 1))
		b.WriteString("\n\n")

		// Show preview
		month, _ := strconv.Atoi(m.monthInput.Value())
		year, _ := strconv.Atoi(m.yearInput.Value())
		if month >= 1 && month <= 12 && year > 2000 {
			from := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
			to := from.AddDate(0, 1, 0).Add(-time.Second)
			preview := fmt.Sprintf("Period: %s to %s", from.Format("Jan 2, 2006"), to.Format("Jan 2, 2006"))
			b.WriteString(styles.HelpStyle.Render(preview))
		}
	} else {
		// Custom date mode
		b.WriteString(m.renderField("From", m.fromInput, m.focusIndex == 0))
		b.WriteString("\n")
		b.WriteString(m.renderField("To", m.toInput, m.focusIndex == 1))
		b.WriteString("\n\n")

		// Show preview
		from, fromErr := time.Parse("2006-01-02", m.fromInput.Value())
		to, toErr := time.Parse("2006-01-02", m.toInput.Value())
		if fromErr == nil && toErr == nil {
			days := int(to.Sub(from).Hours() / 24)
			preview := fmt.Sprintf("Period: %d days", days)
			b.WriteString(styles.HelpStyle.Render(preview))
		}
	}

	// Error message
	if m.errorMsg != "" {
		b.WriteString("\n\n")
		b.WriteString(styles.ErrorStyle.Render("! " + m.errorMsg))
	}

	// Footer
	b.WriteString("\n\n")
	help := "Tab: next field • 1/2: switch mode • Enter: continue"
	b.WriteString(styles.HelpStyle.Render(help))

	return b.String()
}

func (m DateRangeModel) renderField(label string, input textinput.Model, focused bool) string {
	labelStyle := styles.LabelStyle
	if focused {
		labelStyle = styles.FocusedStyle
	}

	return labelStyle.Render(label+": ") + input.View()
}

func (m *DateRangeModel) updateFocus() {
	m.monthInput.Blur()
	m.yearInput.Blur()
	m.fromInput.Blur()
	m.toInput.Blur()

	if m.mode == state.DateRangeMonth {
		if m.focusIndex == 0 {
			m.monthInput.Focus()
		} else {
			m.yearInput.Focus()
		}
	} else {
		if m.focusIndex == 0 {
			m.fromInput.Focus()
		} else {
			m.toInput.Focus()
		}
	}
}

func (m *DateRangeModel) validate() bool {
	if m.mode == state.DateRangeMonth {
		month, err := strconv.Atoi(m.monthInput.Value())
		if err != nil || month < 1 || month > 12 {
			m.errorMsg = "Invalid month (1-12)"
			return false
		}

		year, err := strconv.Atoi(m.yearInput.Value())
		if err != nil || year < 2000 || year > 2100 {
			m.errorMsg = "Invalid year (2000-2100)"
			return false
		}
	} else {
		from, err := time.Parse("2006-01-02", m.fromInput.Value())
		if err != nil {
			m.errorMsg = "Invalid 'from' date (use YYYY-MM-DD)"
			return false
		}

		to, err := time.Parse("2006-01-02", m.toInput.Value())
		if err != nil {
			m.errorMsg = "Invalid 'to' date (use YYYY-MM-DD)"
			return false
		}

		if to.Before(from) {
			m.errorMsg = "'To' date must be after 'from' date"
			return false
		}
	}

	return true
}

func (m *DateRangeModel) saveToState() {
	m.state.DateMode = m.mode

	if m.mode == state.DateRangeMonth {
		m.state.Month, _ = strconv.Atoi(m.monthInput.Value())
		m.state.Year, _ = strconv.Atoi(m.yearInput.Value())
	} else {
		m.state.FromDate, _ = time.Parse("2006-01-02", m.fromInput.Value())
		m.state.ToDate, _ = time.Parse("2006-01-02", m.toInput.Value())
	}
}

func (m DateRangeModel) SetSize(width, height int) DateRangeModel {
	m.width = width
	m.height = height
	return m
}
