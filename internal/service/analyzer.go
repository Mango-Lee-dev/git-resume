package service

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"github.com/wootaiklee/git-resume/internal/db"
	"github.com/wootaiklee/git-resume/internal/git"
	"github.com/wootaiklee/git-resume/internal/llm"
	"github.com/wootaiklee/git-resume/pkg/models"
)

// AnalyzePhase represents the current phase of analysis
type AnalyzePhase string

const (
	PhaseScanning   AnalyzePhase = "scanning"
	PhaseFiltering  AnalyzePhase = "filtering"
	PhaseAnalyzing  AnalyzePhase = "analyzing"
	PhaseSaving     AnalyzePhase = "saving"
	PhaseComplete   AnalyzePhase = "complete"
	PhaseError      AnalyzePhase = "error"
)

// AnalyzeConfig holds analysis parameters
type AnalyzeConfig struct {
	Repos      []string
	FromDate   time.Time
	ToDate     time.Time
	Template   string
	BatchSize  int
	DryRun     bool
}

// AnalyzeProgress reports analysis progress
type AnalyzeProgress struct {
	Phase          AnalyzePhase
	CurrentBatch   int
	TotalBatches   int
	CommitsFound   int
	CommitsToProcess int
	ResultsCreated int
	Progress       int // 0-100
	Message        string
	Error          error
}

// AnalyzeResult holds the final analysis outcome
type AnalyzeResult struct {
	Results         []models.AnalysisResult
	CommitsScanned  int
	CommitsSkipped  int
	CommitsProcessed int
	TokensUsed      int
}

// Analyzer orchestrates the analysis workflow
type Analyzer struct {
	cache       *db.Cache
	apiKey      string
}

// NewAnalyzer creates a new analyzer service
func NewAnalyzer(cache *db.Cache, apiKey string) *Analyzer {
	return &Analyzer{
		cache:  cache,
		apiKey: apiKey,
	}
}

// Analyze runs the full analysis workflow with progress reporting
func (a *Analyzer) Analyze(ctx context.Context, cfg AnalyzeConfig, progressCh chan<- AnalyzeProgress) (*AnalyzeResult, error) {
	defer close(progressCh)

	result := &AnalyzeResult{}

	// Phase 1: Scanning repositories
	sendProgress(progressCh, AnalyzeProgress{
		Phase:   PhaseScanning,
		Message: "Scanning repositories...",
	})

	var allCommits []models.Commit
	repoProjects := make(map[string]string) // hash -> project name

	for _, repoPath := range cfg.Repos {
		// Check for cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		projectName := filepath.Base(repoPath)

		parser, err := git.NewParser(repoPath)
		if err != nil {
			continue // Skip failed repos
		}

		commits, err := parser.GetCommits(cfg.FromDate, cfg.ToDate)
		if err != nil {
			continue
		}

		for i := range commits {
			repoProjects[commits[i].Hash] = projectName
		}

		allCommits = append(allCommits, commits...)
	}

	result.CommitsScanned = len(allCommits)

	if len(allCommits) == 0 {
		sendProgress(progressCh, AnalyzeProgress{
			Phase:    PhaseComplete,
			Progress: 100,
			Message:  "No commits found",
		})
		return result, nil
	}

	sendProgress(progressCh, AnalyzeProgress{
		Phase:        PhaseScanning,
		CommitsFound: len(allCommits),
		Progress:     10,
		Message:      fmt.Sprintf("Found %d commits", len(allCommits)),
	})

	// Phase 2: Filtering
	sendProgress(progressCh, AnalyzeProgress{
		Phase:        PhaseFiltering,
		CommitsFound: len(allCommits),
		Progress:     20,
		Message:      "Filtering processed commits...",
	})

	unprocessed, err := a.cache.FilterUnprocessed(allCommits)
	if err != nil {
		sendProgress(progressCh, AnalyzeProgress{
			Phase:   PhaseError,
			Message: "Failed to filter commits",
			Error:   err,
		})
		return nil, fmt.Errorf("failed to filter commits: %w", err)
	}

	result.CommitsSkipped = len(allCommits) - len(unprocessed)
	result.CommitsProcessed = len(unprocessed)

	if len(unprocessed) == 0 {
		sendProgress(progressCh, AnalyzeProgress{
			Phase:    PhaseComplete,
			Progress: 100,
			Message:  "All commits already processed",
		})
		return result, nil
	}

	// Sort by score
	sort.Slice(unprocessed, func(i, j int) bool {
		return unprocessed[i].Score > unprocessed[j].Score
	})

	sendProgress(progressCh, AnalyzeProgress{
		Phase:            PhaseFiltering,
		CommitsFound:     len(allCommits),
		CommitsToProcess: len(unprocessed),
		Progress:         30,
		Message:          fmt.Sprintf("%d commits to process", len(unprocessed)),
	})

	// Dry run check
	if cfg.DryRun {
		sendProgress(progressCh, AnalyzeProgress{
			Phase:            PhaseComplete,
			CommitsFound:     len(allCommits),
			CommitsToProcess: len(unprocessed),
			Progress:         100,
			Message:          "Dry run complete",
		})
		return result, nil
	}

	// Check API key
	if a.apiKey == "" {
		err := fmt.Errorf("CLAUDE_API_KEY not set")
		sendProgress(progressCh, AnalyzeProgress{
			Phase:   PhaseError,
			Message: err.Error(),
			Error:   err,
		})
		return nil, err
	}

	// Phase 3: Analyzing
	templateMgr := llm.NewTemplateManager()
	if err := templateMgr.SetTemplate(cfg.Template); err != nil {
		// Use default template
		templateMgr.SetTemplate("default")
	}

	client := llm.NewClientWithTemplate(a.apiKey, templateMgr)

	batches := createBatches(unprocessed, cfg.BatchSize, repoProjects)
	totalBatches := len(batches)

	for i, batch := range batches {
		// Check for cancellation
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		progress := 30 + int(float64(i)/float64(totalBatches)*60)
		sendProgress(progressCh, AnalyzeProgress{
			Phase:          PhaseAnalyzing,
			CurrentBatch:   i + 1,
			TotalBatches:   totalBatches,
			ResultsCreated: len(result.Results),
			Progress:       progress,
			Message:        fmt.Sprintf("Processing batch %d/%d", i+1, totalBatches),
		})

		batchResult, err := client.AnalyzeCommits(batch)
		if err != nil {
			continue // Skip failed batches
		}

		// Phase 4: Saving (inline)
		if err := a.cache.SaveResults(batchResult.Results); err != nil {
			continue
		}

		a.cache.RecordTokenUsage(
			fmt.Sprintf("batch-%d", i+1),
			batchResult.TokensUsed/2,
			batchResult.TokensUsed/2,
			0,
		)

		result.Results = append(result.Results, batchResult.Results...)
		result.TokensUsed += batchResult.TokensUsed
	}

	// Complete
	sendProgress(progressCh, AnalyzeProgress{
		Phase:          PhaseComplete,
		TotalBatches:   totalBatches,
		ResultsCreated: len(result.Results),
		Progress:       100,
		Message:        fmt.Sprintf("Generated %d bullet points", len(result.Results)),
	})

	return result, nil
}

// EstimateTokens estimates token usage for commits
func (a *Analyzer) EstimateTokens(commits []models.Commit) int {
	total := 0
	for _, c := range commits {
		total += llm.EstimateTokens(models.CommitBatch{Commits: []models.Commit{c}})
	}
	return total
}

// sendProgress safely sends progress to channel
func sendProgress(ch chan<- AnalyzeProgress, p AnalyzeProgress) {
	select {
	case ch <- p:
	default:
		// Channel full or closed, skip
	}
}

// createBatches groups commits into batches with project info
func createBatches(commits []models.Commit, size int, repoProjects map[string]string) []models.CommitBatch {
	var batches []models.CommitBatch

	for i := 0; i < len(commits); i += size {
		end := i + size
		if end > len(commits) {
			end = len(commits)
		}

		batchCommits := commits[i:end]

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
