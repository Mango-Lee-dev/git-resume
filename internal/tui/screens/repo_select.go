package screens

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wootaiklee/git-resume/internal/tui/state"
	"github.com/wootaiklee/git-resume/internal/tui/styles"
)

type RepoSelectModel struct {
	state         *state.AppState
	textInput     textinput.Model
	selectedRepos []string
	currentDir    string
	dirEntries    []os.DirEntry
	dirCursor     int
	mode          repoMode
	errorMsg      string
	width, height int
}

type repoMode int

const (
	modeBrowse repoMode = iota
	modeManual
)

func NewRepoSelectModel(s *state.AppState) RepoSelectModel {
	ti := textinput.New()
	ti.Placeholder = "Enter repository path..."
	ti.CharLimit = 256
	ti.Width = 50

	cwd, _ := os.Getwd()

	m := RepoSelectModel{
		state:         s,
		textInput:     ti,
		selectedRepos: append([]string{}, s.RepoPaths...),
		currentDir:    cwd,
		mode:          modeBrowse,
	}
	m.loadDirEntries()

	return m
}

func (m RepoSelectModel) Init() tea.Cmd {
	return nil
}

func (m RepoSelectModel) Update(msg tea.Msg) (RepoSelectModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.errorMsg = ""

		switch msg.String() {
		case "tab":
			// Toggle between browse and manual mode
			if m.mode == modeBrowse {
				m.mode = modeManual
				m.textInput.Focus()
			} else {
				m.mode = modeBrowse
				m.textInput.Blur()
			}
			return m, nil

		case "enter":
			if m.mode == modeManual {
				// Add path from text input
				path := strings.TrimSpace(m.textInput.Value())
				if path != "" {
					if m.addRepo(path) {
						m.textInput.SetValue("")
					}
				}
			} else {
				// In browse mode
				if len(m.dirEntries) > 0 && m.dirCursor < len(m.dirEntries) {
					entry := m.dirEntries[m.dirCursor]
					if entry.IsDir() {
						newPath := filepath.Join(m.currentDir, entry.Name())
						if isGitRepo(newPath) {
							m.addRepo(newPath)
						} else {
							m.currentDir = newPath
							m.loadDirEntries()
							m.dirCursor = 0
						}
					}
				}
			}
			return m, nil

		case "backspace":
			if m.mode == modeBrowse && m.textInput.Value() == "" {
				// Go up one directory
				parent := filepath.Dir(m.currentDir)
				if parent != m.currentDir {
					m.currentDir = parent
					m.loadDirEntries()
					m.dirCursor = 0
				}
				return m, nil
			}

		case "up", "k":
			if m.mode == modeBrowse && m.dirCursor > 0 {
				m.dirCursor--
			}
			return m, nil

		case "down", "j":
			if m.mode == modeBrowse && m.dirCursor < len(m.dirEntries)-1 {
				m.dirCursor++
			}
			return m, nil

		case "d":
			// Delete last selected repo
			if len(m.selectedRepos) > 0 {
				m.selectedRepos = m.selectedRepos[:len(m.selectedRepos)-1]
			}
			return m, nil

		case "ctrl+n", "ctrl+right":
			// Proceed to next screen
			if len(m.selectedRepos) > 0 {
				m.state.RepoPaths = m.selectedRepos
				return m, Navigate(NavToDateRange)
			}
			m.errorMsg = "Please select at least one repository"
			return m, nil

		case ".":
			// Add current directory if it's a git repo
			if isGitRepo(m.currentDir) {
				m.addRepo(m.currentDir)
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Update text input
	if m.mode == modeManual {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m RepoSelectModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(styles.TitleStyle.Render("Select Repositories"))
	b.WriteString("\n\n")

	// Current path
	pathStyle := lipgloss.NewStyle().Foreground(styles.Subtle)
	b.WriteString(pathStyle.Render("Current: " + m.currentDir))
	b.WriteString("\n\n")

	// Mode indicator
	browseIndicator := "● Browse"
	manualIndicator := "○ Manual"
	if m.mode == modeManual {
		browseIndicator = "○ Browse"
		manualIndicator = "● Manual"
	}
	b.WriteString(styles.LabelStyle.Render(browseIndicator + "  " + manualIndicator + "  (Tab to switch)"))
	b.WriteString("\n\n")

	if m.mode == modeManual {
		// Manual input mode
		b.WriteString("Enter path: ")
		b.WriteString(m.textInput.View())
		b.WriteString("\n")
	} else {
		// Directory browser
		maxItems := 10
		start := 0
		if m.dirCursor >= maxItems {
			start = m.dirCursor - maxItems + 1
		}

		for i := start; i < len(m.dirEntries) && i < start+maxItems; i++ {
			entry := m.dirEntries[i]
			name := entry.Name()

			// Style based on selection and type
			cursor := "  "
			if i == m.dirCursor {
				cursor = "▸ "
			}

			itemStyle := styles.NormalItemStyle
			if i == m.dirCursor {
				itemStyle = styles.SelectedItemStyle
			}

			// Add git indicator
			fullPath := filepath.Join(m.currentDir, name)
			gitIndicator := ""
			if entry.IsDir() && isGitRepo(fullPath) {
				gitIndicator = " " + styles.SuccessStyle.Render("[git]")
			}

			dirIndicator := ""
			if entry.IsDir() {
				dirIndicator = "/"
			}

			b.WriteString(cursor)
			b.WriteString(itemStyle.Render(name + dirIndicator))
			b.WriteString(gitIndicator)
			b.WriteString("\n")
		}

		if len(m.dirEntries) == 0 {
			b.WriteString(styles.HelpStyle.Render("  (empty directory)"))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")

	// Selected repos
	b.WriteString(styles.LabelStyle.Render("Selected repositories:"))
	b.WriteString("\n")
	if len(m.selectedRepos) == 0 {
		b.WriteString(styles.HelpStyle.Render("  (none selected)"))
		b.WriteString("\n")
	} else {
		for i, repo := range m.selectedRepos {
			b.WriteString(styles.SuccessStyle.Render("  ✓ "))
			b.WriteString(truncate(repo, 50))
			if i < len(m.selectedRepos)-1 {
				b.WriteString("\n")
			}
		}
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
	help := "↑/↓: navigate • Enter: select/open • .: add current • d: remove last • Ctrl+N: next"
	b.WriteString(styles.HelpStyle.Render(help))

	return b.String()
}

func (m *RepoSelectModel) loadDirEntries() {
	entries, err := os.ReadDir(m.currentDir)
	if err != nil {
		m.dirEntries = nil
		return
	}

	// Filter to show only directories
	var dirs []os.DirEntry
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			dirs = append(dirs, entry)
		}
	}
	m.dirEntries = dirs
}

func (m *RepoSelectModel) addRepo(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		m.errorMsg = "Invalid path"
		return false
	}

	if !isGitRepo(absPath) {
		m.errorMsg = "Not a git repository"
		return false
	}

	// Check for duplicates
	for _, r := range m.selectedRepos {
		if r == absPath {
			m.errorMsg = "Repository already selected"
			return false
		}
	}

	m.selectedRepos = append(m.selectedRepos, absPath)
	return true
}

func isGitRepo(path string) bool {
	gitPath := filepath.Join(path, ".git")
	info, err := os.Stat(gitPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (m RepoSelectModel) SetSize(width, height int) RepoSelectModel {
	m.width = width
	m.height = height
	m.textInput.Width = width - 20
	return m
}
