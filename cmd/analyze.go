package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wootaiklee/git-resume/internal/db"
	"github.com/wootaiklee/git-resume/internal/export"
	"github.com/wootaiklee/git-resume/internal/git"
	"github.com/wootaiklee/git-resume/internal/llm"
	"github.com/wootaiklee/git-resume/internal/notify"
	"github.com/wootaiklee/git-resume/internal/ui"
	"github.com/wootaiklee/git-resume/pkg/models"
)

var (
	month        int
	year         int
	fromDate     string
	toDate       string
	dryRun       bool
	batchSize    int
	repos        string
	templateName string
	notifySlack  bool
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze git commits and generate resume bullet points",
	Long: `Analyze git commits from the specified repository and generate
STAR-format resume bullet points using Claude AI.

Example:
  git-resume analyze --month=4
  git-resume analyze --from=2024-01-01 --to=2024-03-31
  git-resume analyze --repos=/path/to/repo1,/path/to/repo2
  git-resume analyze --template=startup
  git-resume analyze --dry-run --notify`,
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().IntVarP(&month, "month", "m", 0, "month to analyze (1-12)")
	analyzeCmd.Flags().IntVarP(&year, "year", "y", 0, "year to analyze (defaults to current year)")
	analyzeCmd.Flags().StringVar(&fromDate, "from", "", "start date (YYYY-MM-DD)")
	analyzeCmd.Flags().StringVar(&toDate, "to", "", "end date (YYYY-MM-DD)")
	analyzeCmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be analyzed without calling API")
	analyzeCmd.Flags().IntVar(&batchSize, "batch-size", 5, "number of commits per API call")
	analyzeCmd.Flags().StringVar(&repos, "repos", "", "comma-separated list of repository paths")
	analyzeCmd.Flags().StringVarP(&templateName, "template", "t", "default", "prompt template (default, startup, enterprise, backend, frontend, devops, data)")
	analyzeCmd.Flags().BoolVar(&notifySlack, "notify", false, "send Slack notification when complete")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	// Determine date range
	from, to, err := parseDateRange()
	if err != nil {
		return err
	}

	// Get repository paths
	repoPaths := getRepoPaths()
	if len(repoPaths) == 0 {
		return fmt.Errorf("no repository paths specified")
	}

	fmt.Printf("Date range: %s to %s\n", from.Format("2006-01-02"), to.Format("2006-01-02"))
	fmt.Printf("Repositories to analyze: %d\n\n", len(repoPaths))

	// Collect commits from all repositories
	var allCommits []models.Commit
	repoProjects := make(map[string]string) // hash -> project name mapping

	for _, repoPath := range repoPaths {
		projectName := filepath.Base(repoPath)
		fmt.Printf("Scanning repository: %s\n", repoPath)

		parser, err := git.NewParser(repoPath)
		if err != nil {
			fmt.Printf("  Warning: failed to open repository: %v\n", err)
			continue
		}

		commits, err := parser.GetCommits(from, to)
		if err != nil {
			fmt.Printf("  Warning: failed to get commits: %v\n", err)
			continue
		}

		fmt.Printf("  Found %d commits after filtering\n", len(commits))

		// Track project for each commit
		for i := range commits {
			repoProjects[commits[i].Hash] = projectName
		}

		allCommits = append(allCommits, commits...)
	}

	fmt.Printf("\nTotal commits found: %d\n", len(allCommits))

	if len(allCommits) == 0 {
		fmt.Println("No commits to analyze")
		return nil
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
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer database.Close()

	cache := db.NewCache(database)

	// Filter out already processed commits
	unprocessed, err := cache.FilterUnprocessed(allCommits)
	if err != nil {
		return fmt.Errorf("failed to filter commits: %w", err)
	}

	fmt.Printf("New commits to process: %d\n", len(unprocessed))

	if len(unprocessed) == 0 {
		fmt.Println("All commits have already been processed")
		return nil
	}

	// Sort by score (highest first)
	sortByScore(unprocessed)

	// Show commits to be processed
	fmt.Println("\nCommits to analyze:")
	for i, c := range unprocessed {
		fmt.Printf("  %d. [Score: %d] %s - %s\n",
			i+1, c.Score, c.Hash[:7], truncate(c.Message, 50))
	}

	if dryRun {
		fmt.Println("\n[Dry run] Would process these commits. Exiting.")
		estimatedTokens := 0
		for _, c := range unprocessed {
			estimatedTokens += llm.EstimateTokens(models.CommitBatch{Commits: []models.Commit{c}})
		}
		fmt.Printf("Estimated input tokens: ~%d\n", estimatedTokens)
		return nil
	}

	// Get API key
	apiKey := viper.GetString("CLAUDE_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("CLAUDE_API_KEY not set. Please set it in .env file or environment")
	}

	// Initialize template manager
	templateMgr := llm.NewTemplateManager()
	if err := templateMgr.SetTemplate(templateName); err != nil {
		fmt.Printf("Warning: %v, using default template\n", err)
	} else {
		fmt.Printf("Using template: %s\n", templateMgr.GetTemplate().Name)
	}

	// Initialize LLM client with template
	client := llm.NewClientWithTemplate(apiKey, templateMgr)

	// Process in batches
	var allResults []models.AnalysisResult

	batches := createBatchesMultiRepo(unprocessed, batchSize, repoProjects)
	fmt.Printf("\nProcessing %d batches...\n", len(batches))

	// Initialize progress bar
	progress := ui.NewProgress(len(batches))
	progress.Start("Initializing batch processing")

	for i, batch := range batches {
		progress.Update(i, fmt.Sprintf("Batch %d/%d (%d commits)", i+1, len(batches), len(batch.Commits)))

		result, err := client.AnalyzeCommits(batch)
		if err != nil {
			fmt.Printf("\nWarning: batch %d failed: %v\n", i+1, err)
			continue
		}

		// Save results to database
		if err := cache.SaveResults(result.Results); err != nil {
			fmt.Printf("\nWarning: failed to save results: %v\n", err)
		}

		// Track token usage
		cache.RecordTokenUsage(
			fmt.Sprintf("batch-%d", i+1),
			result.TokensUsed/2, // Rough split
			result.TokensUsed/2,
			0, // Cost calculation would go here
		)

		allResults = append(allResults, result.Results...)
	}

	progress.Complete()
	fmt.Printf("Generated %d resume bullet points\n", len(allResults))

	// Export results
	if len(allResults) > 0 {
		outputDir := viper.GetString("output")
		if outputDir == "" {
			outputDir = viper.GetString("OUTPUT_DIR")
		}
		if outputDir == "" {
			outputDir = "./output"
		}

		exporter := export.NewCSVExporter(outputDir)
		outputPath, err := exporter.ExportWithTimestamp(allResults, "resume")
		if err != nil {
			return fmt.Errorf("failed to export results: %w", err)
		}

		fmt.Printf("\nResults exported to: %s\n", outputPath)
		fmt.Printf("Total bullet points generated: %d\n", len(allResults))

		// Send Slack notification if enabled
		if notifySlack {
			slackURL := viper.GetString("SLACK_WEBHOOK_URL")
			if slackURL != "" {
				notifier := notify.NewSlackNotifier(slackURL)
				if err := notifier.SendAnalysisComplete(allResults, from, to); err != nil {
					fmt.Printf("Warning: failed to send Slack notification: %v\n", err)
				} else {
					fmt.Println("Slack notification sent")
				}
			} else {
				fmt.Println("Warning: SLACK_WEBHOOK_URL not set, skipping notification")
			}
		}
	}

	return nil
}

func parseDateRange() (time.Time, time.Time, error) {
	now := time.Now()

	// If specific dates provided
	if fromDate != "" && toDate != "" {
		from, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid from date: %w", err)
		}
		to, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid to date: %w", err)
		}
		return from, to.Add(24*time.Hour - time.Second), nil
	}

	// If month specified
	if month > 0 {
		y := year
		if y == 0 {
			y = now.Year()
		}
		from := time.Date(y, time.Month(month), 1, 0, 0, 0, 0, time.Local)
		to := from.AddDate(0, 1, 0).Add(-time.Second)
		return from, to, nil
	}

	// Default: current month
	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	to := now
	return from, to, nil
}

func createBatches(commits []models.Commit, size int, project string) []models.CommitBatch {
	var batches []models.CommitBatch

	for i := 0; i < len(commits); i += size {
		end := i + size
		if end > len(commits) {
			end = len(commits)
		}

		batch := models.CommitBatch{
			Commits: commits[i:end],
			Project: project,
		}

		if len(batch.Commits) > 0 {
			batch.FromDate = batch.Commits[len(batch.Commits)-1].Date
			batch.ToDate = batch.Commits[0].Date
		}

		batches = append(batches, batch)
	}

	return batches
}

func sortByScore(commits []models.Commit) {
	// Simple bubble sort (commits are usually < 100)
	for i := 0; i < len(commits); i++ {
		for j := i + 1; j < len(commits); j++ {
			if commits[j].Score > commits[i].Score {
				commits[i], commits[j] = commits[j], commits[i]
			}
		}
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// getRepoPaths returns a list of repository paths to analyze
func getRepoPaths() []string {
	var paths []string

	// Check --repos flag first
	if repos != "" {
		for _, path := range strings.Split(repos, ",") {
			path = strings.TrimSpace(path)
			if path != "" {
				paths = append(paths, path)
			}
		}
		return paths
	}

	// Check --repo flag
	repoPath := viper.GetString("repo")
	if repoPath != "" {
		return []string{repoPath}
	}

	// Check DEFAULT_REPO_PATH env
	repoPath = viper.GetString("DEFAULT_REPO_PATH")
	if repoPath != "" {
		return []string{repoPath}
	}

	// Use current directory
	cwd, err := os.Getwd()
	if err == nil {
		return []string{cwd}
	}

	return nil
}

// createBatchesMultiRepo creates batches from commits across multiple repositories
func createBatchesMultiRepo(commits []models.Commit, size int, repoProjects map[string]string) []models.CommitBatch {
	var batches []models.CommitBatch

	for i := 0; i < len(commits); i += size {
		end := i + size
		if end > len(commits) {
			end = len(commits)
		}

		batchCommits := commits[i:end]

		// Determine project name (use first commit's project or "multi-repo")
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
