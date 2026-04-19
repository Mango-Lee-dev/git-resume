package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/wootaiklee/git-resume/internal/tui"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch interactive Terminal UI",
	Long: `Launch an interactive terminal user interface for git-resume.

The TUI provides a guided workflow for:
- Selecting repositories to analyze
- Configuring date ranges
- Choosing prompt templates
- Running analysis with real-time progress
- Previewing and exporting results

Example:
  git-resume tui`,
	RunE: runTUI,
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

func runTUI(cmd *cobra.Command, args []string) error {
	p := tea.NewProgram(
		tui.New(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	model, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	// Show final output path if available
	if m, ok := model.(tui.Model); ok {
		state := m.State()
		if state.OutputPath != "" {
			fmt.Printf("\nResults saved to: %s\n", state.OutputPath)
		}
	}

	return nil
}
