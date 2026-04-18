// Package ui provides terminal user interface components.
package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mattn/go-isatty"
)

// Progress displays a progress bar with step information, percentage,
// elapsed time, and current task description in the terminal.
type Progress struct {
	total       int
	current     int
	description string
	startTime   time.Time
	isTTY       bool
	width       int
	writer      io.Writer
	mu          sync.Mutex
}

// NewProgress creates a new Progress instance with the given total steps.
// It automatically detects whether stdout is a TTY and adjusts behavior accordingly.
func NewProgress(total int) *Progress {
	isTTY := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
	width := 40 // default progress bar width

	return &Progress{
		total:  total,
		isTTY:  isTTY,
		width:  width,
		writer: os.Stdout,
	}
}

// Start begins the progress tracking with an initial description.
// This should be called once before any Update calls.
func (p *Progress) Start(description string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.startTime = time.Now()
	p.current = 0
	p.description = description

	p.render()
}

// Update updates the progress to the given step with a new description.
// The current value should be between 0 and total (inclusive).
func (p *Progress) Update(current int, description string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.current = current
	p.description = description

	p.render()
}

// Complete marks the progress as finished and displays the completion message.
func (p *Progress) Complete() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.current = p.total
	p.description = "Complete"

	p.render()

	// Print newline to move past the progress bar
	if p.isTTY {
		fmt.Fprintln(p.writer)
	}
}

// render draws the progress bar to the terminal.
// For TTY: displays an updating progress bar on a single line.
// For non-TTY: prints simple status messages.
func (p *Progress) render() {
	elapsed := time.Since(p.startTime)
	percentage := 0
	if p.total > 0 {
		percentage = (p.current * 100) / p.total
	}

	if p.isTTY {
		p.renderTTY(elapsed, percentage)
	} else {
		p.renderNonTTY(elapsed, percentage)
	}
}

// renderTTY renders an interactive progress bar for terminal output.
func (p *Progress) renderTTY(elapsed time.Duration, percentage int) {
	// Calculate filled width of the progress bar
	filledWidth := 0
	if p.total > 0 {
		filledWidth = (p.current * p.width) / p.total
	}

	// Build the progress bar
	bar := strings.Repeat("=", filledWidth)
	if filledWidth < p.width {
		bar += ">"
		bar += strings.Repeat(" ", p.width-filledWidth-1)
	}

	// Format elapsed time
	elapsedStr := formatDuration(elapsed)

	// Truncate description if too long
	desc := p.description
	maxDescLen := 30
	if len(desc) > maxDescLen {
		desc = desc[:maxDescLen-3] + "..."
	}

	// Print the progress bar with carriage return to overwrite
	fmt.Fprintf(p.writer, "\r[%s] %3d%% (%d/%d) %s | %s",
		bar,
		percentage,
		p.current,
		p.total,
		elapsedStr,
		desc,
	)

	// Clear any remaining characters from previous longer lines
	fmt.Fprintf(p.writer, "    ")
}

// renderNonTTY renders simple progress messages for non-terminal output.
func (p *Progress) renderNonTTY(elapsed time.Duration, percentage int) {
	elapsedStr := formatDuration(elapsed)

	fmt.Fprintf(p.writer, "[%d/%d] %3d%% | %s | %s\n",
		p.current,
		p.total,
		percentage,
		elapsedStr,
		p.description,
	)
}

// formatDuration formats a duration in a human-readable format (MM:SS or HH:MM:SS).
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// IsTTY returns whether the progress bar is running in a TTY environment.
func (p *Progress) IsTTY() bool {
	return p.isTTY
}

// SetWriter sets a custom writer for the progress output.
// This is useful for testing or redirecting output.
func (p *Progress) SetWriter(w io.Writer) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.writer = w
	// When using a custom writer, assume non-TTY mode
	p.isTTY = false
}

// SetWidth sets the width of the progress bar (number of characters).
func (p *Progress) SetWidth(width int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if width > 0 {
		p.width = width
	}
}
