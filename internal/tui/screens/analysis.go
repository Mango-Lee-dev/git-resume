package screens

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
	"github.com/wootaiklee/git-resume/internal/db"
	"github.com/wootaiklee/git-resume/internal/git"
	"github.com/wootaiklee/git-resume/internal/llm"
	"github.com/wootaiklee/git-resume/internal/tui/state"
	"github.com/wootaiklee/git-resume/internal/tui/styles"
	"github.com/wootaiklee/git-resume/pkg/models"
)

// Analysis phase
type analysisPhase int

const (
	phaseScanning analysisPhase = iota
	phaseFiltering
	phaseProcessing
	phaseSaving
	phaseComplete
	phaseError
)

// Messages for async operations
type commitsFetchedMsg struct {
	commits      []models.Commit
	repoProjects map[string]string
	err          error
}

type filterCompleteMsg struct {
	unprocessed []models.Commit
	err         error
}

type batchProcessedMsg struct {
	batchNum   int
	results    []models.AnalysisResult
	tokensUsed int
	err        error
}

type analysisCompleteMsg struct {
	err error
}

type AnalysisModel struct {
	state *state.AppState

	spinner  spinner.Model
	progress progress.Model

	phase      analysisPhase
	statusText string

	commits      []models.Commit
	repoProjects map[string]string
	unprocessed  []models.Commit
	results      []models.AnalysisResult

	currentBatch int
	totalBatches int
	tokensUsed   int

	err    error
	ctx    context.Context
	cancel context.CancelFunc

	width, height int
}

func NewAnalysisModel(s *state.AppState) AnalysisModel {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = styles.ActiveStyle

	p := progress.New(progress.WithDefaultGradient())

	ctx, cancel := context.WithCancel(context.Background())

	return AnalysisModel{
		state:    s,
		spinner:  sp,
		progress: p,
		phase:    phaseScanning,
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (m AnalysisModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.scanRepositories(),
	)
}

func (m AnalysisModel) Update(msg tea.Msg) (AnalysisModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.cancel()
			m.phase = phaseError
			m.err = fmt.Errorf("cancelled by user")
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = msg.Width - 10

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case commitsFetchedMsg:
		if msg.err != nil {
			m.phase = phaseError
			m.err = msg.err
			return m, nil
		}

		m.commits = msg.commits
		m.repoProjects = msg.repoProjects
		m.state.Commits = msg.commits
		m.statusText = fmt.Sprintf("Found %d commits", len(msg.commits))

		if len(msg.commits) == 0 {
			m.phase = phaseComplete
			m.statusText = "No commits found in the selected date range"
			return m, Navigate(NavToResults)
		}

		m.phase = phaseFiltering
		return m, m.filterCommits()

	case filterCompleteMsg:
		if msg.err != nil {
			m.phase = phaseError
			m.err = msg.err
			return m, nil
		}

		m.unprocessed = msg.unprocessed
		m.statusText = fmt.Sprintf("%d new commits to process", len(msg.unprocessed))

		if len(msg.unprocessed) == 0 {
			m.phase = phaseComplete
			m.statusText = "All commits have already been processed"
			return m, Navigate(NavToResults)
		}

		// Check for dry run
		if m.state.DryRun {
			m.phase = phaseComplete
			m.statusText = fmt.Sprintf("Dry run: would process %d commits", len(m.unprocessed))
			// Estimate tokens
			tokens := 0
			for _, c := range m.unprocessed {
				tokens += llm.EstimateTokens(models.CommitBatch{Commits: []models.Commit{c}})
			}
			m.statusText += fmt.Sprintf(" (~%d tokens)", tokens)
			return m, Navigate(NavToResults)
		}

		m.phase = phaseProcessing
		m.totalBatches = (len(m.unprocessed) + m.state.BatchSize - 1) / m.state.BatchSize
		return m, m.processBatches()

	case batchProcessedMsg:
		if msg.err != nil {
			m.statusText = fmt.Sprintf("Batch %d failed: %v", msg.batchNum, msg.err)
			// Continue with next batch
		} else {
			m.results = append(m.results, msg.results...)
			m.tokensUsed += msg.tokensUsed
		}

		m.currentBatch = msg.batchNum
		m.statusText = fmt.Sprintf("Processed batch %d/%d", m.currentBatch, m.totalBatches)

		// Update progress
		percent := float64(m.currentBatch) / float64(m.totalBatches)
		return m, m.progress.SetPercent(percent)

	case analysisCompleteMsg:
		if msg.err != nil {
			m.phase = phaseError
			m.err = msg.err
			return m, nil
		}

		m.state.Results = m.results
		m.state.TokensUsed = m.tokensUsed
		m.phase = phaseComplete

		// Transition to results screen
		return m, Navigate(NavToResults)
	}

	return m, nil
}

func (m AnalysisModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(styles.TitleStyle.Render("Analyzing Commits"))
	b.WriteString("\n\n")

	// Phase indicator
	switch m.phase {
	case phaseScanning:
		b.WriteString(m.spinner.View())
		b.WriteString(" Scanning repositories...\n")

	case phaseFiltering:
		b.WriteString(m.spinner.View())
		b.WriteString(" Filtering commits...\n")

	case phaseProcessing:
		b.WriteString(fmt.Sprintf("Processing batch %d/%d\n", m.currentBatch, m.totalBatches))
		b.WriteString(m.progress.View())
		b.WriteString("\n")

	case phaseSaving:
		b.WriteString(m.spinner.View())
		b.WriteString(" Saving results...\n")

	case phaseComplete:
		b.WriteString(styles.SuccessStyle.Render("✓ Complete"))
		b.WriteString("\n")

	case phaseError:
		b.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("✗ Error: %v", m.err)))
		b.WriteString("\n")
	}

	// Status text
	if m.statusText != "" {
		b.WriteString("\n")
		b.WriteString(m.statusText)
		b.WriteString("\n")
	}

	// Statistics
	if len(m.results) > 0 {
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("Generated: %d bullet points\n", len(m.results)))
		b.WriteString(fmt.Sprintf("Tokens used: %d\n", m.tokensUsed))
	}

	// Footer
	b.WriteString("\n")
	if m.phase == phaseProcessing || m.phase == phaseScanning || m.phase == phaseFiltering {
		b.WriteString(styles.HelpStyle.Render("Ctrl+C: cancel"))
	} else if m.phase == phaseError {
		b.WriteString(styles.HelpStyle.Render("Esc: go back"))
	}

	return b.String()
}

// scanRepositories scans all selected repositories for commits
func (m AnalysisModel) scanRepositories() tea.Cmd {
	return func() tea.Msg {
		from, to := m.state.GetDateRange()

		var allCommits []models.Commit
		repoProjects := make(map[string]string)

		for _, repoPath := range m.state.RepoPaths {
			projectName := filepath.Base(repoPath)

			parser, err := git.NewParser(repoPath)
			if err != nil {
				continue
			}

			commits, err := parser.GetCommits(from, to)
			if err != nil {
				continue
			}

			for i := range commits {
				repoProjects[commits[i].Hash] = projectName
			}

			allCommits = append(allCommits, commits...)
		}

		return commitsFetchedMsg{
			commits:      allCommits,
			repoProjects: repoProjects,
		}
	}
}

// filterCommits filters out already processed commits
func (m AnalysisModel) filterCommits() tea.Cmd {
	return func() tea.Msg {
		dbPath := viper.GetString("db")
		if dbPath == "" {
			dbPath = viper.GetString("DB_PATH")
		}
		if dbPath == "" {
			dbPath = "./data/cache.db"
		}

		database, err := db.New(dbPath)
		if err != nil {
			// If we can't connect to DB, assume all commits are unprocessed
			return filterCompleteMsg{unprocessed: m.commits}
		}
		defer database.Close()

		cache := db.NewCache(database)
		unprocessed, err := cache.FilterUnprocessed(m.commits)
		if err != nil {
			return filterCompleteMsg{err: err}
		}

		// Sort by score (highest first)
		sortByScore(unprocessed)

		return filterCompleteMsg{unprocessed: unprocessed}
	}
}

// processBatches processes commits in batches
func (m AnalysisModel) processBatches() tea.Cmd {
	return func() tea.Msg {
		// Get API key
		apiKey := viper.GetString("CLAUDE_API_KEY")
		if apiKey == "" {
			return analysisCompleteMsg{err: fmt.Errorf("CLAUDE_API_KEY not set")}
		}

		// Initialize database
		dbPath := viper.GetString("db")
		if dbPath == "" {
			dbPath = viper.GetString("DB_PATH")
		}
		if dbPath == "" {
			dbPath = "./data/cache.db"
		}

		database, err := db.New(dbPath)
		if err != nil {
			return analysisCompleteMsg{err: fmt.Errorf("failed to initialize database: %w", err)}
		}
		defer database.Close()

		cache := db.NewCache(database)

		// Initialize LLM client
		templateMgr := llm.NewTemplateManager()
		templateMgr.SetTemplate(m.state.TemplateName)
		client := llm.NewClientWithTemplate(apiKey, templateMgr)

		// Create batches
		batches := createBatches(m.unprocessed, m.state.BatchSize, m.repoProjects)

		var allResults []models.AnalysisResult
		var totalTokens int

		for i, batch := range batches {
			select {
			case <-m.ctx.Done():
				return analysisCompleteMsg{err: m.ctx.Err()}
			default:
			}

			result, err := client.AnalyzeCommits(batch)
			if err != nil {
				continue
			}

			// Save to cache
			cache.SaveResults(result.Results)
			allResults = append(allResults, result.Results...)
			totalTokens += result.TokensUsed

			// Note: In a real implementation, we'd send progress updates via channels
			// For simplicity, we process all batches and return final result
			_ = i // batch number
		}

		// Update state
		m.state.Results = allResults
		m.state.TokensUsed = totalTokens

		return analysisCompleteMsg{}
	}
}

// createBatches creates commit batches for processing
func createBatches(commits []models.Commit, size int, repoProjects map[string]string) []models.CommitBatch {
	var batches []models.CommitBatch

	for i := 0; i < len(commits); i += size {
		end := i + size
		if end > len(commits) {
			end = len(commits)
		}

		batchCommits := commits[i:end]

		// Determine project name
		projectName := "multi-repo"
		if len(batchCommits) > 0 {
			if proj, ok := repoProjects[batchCommits[0].Hash]; ok {
				projectName = proj
			}
		}

		batch := models.CommitBatch{
			Commits: batchCommits,
			Project: projectName,
		}

		if len(batch.Commits) > 0 {
			batch.FromDate = batch.Commits[len(batch.Commits)-1].Date
			batch.ToDate = batch.Commits[0].Date
		}

		batches = append(batches, batch)
	}

	return batches
}

// sortByScore sorts commits by importance score (highest first)
func sortByScore(commits []models.Commit) {
	for i := 0; i < len(commits); i++ {
		for j := i + 1; j < len(commits); j++ {
			if commits[j].Score > commits[i].Score {
				commits[i], commits[j] = commits[j], commits[i]
			}
		}
	}
}

func (m AnalysisModel) SetSize(width, height int) AnalysisModel {
	m.width = width
	m.height = height
	m.progress.Width = width - 10
	return m
}
