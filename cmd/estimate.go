package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wootaiklee/git-resume/internal/db"
	"github.com/wootaiklee/git-resume/internal/git"
	"github.com/wootaiklee/git-resume/internal/llm"
	"github.com/wootaiklee/git-resume/pkg/models"
)

// Claude API pricing (as of 2024)
const (
	inputTokenCostPer1K  = 0.003  // $3 per 1M input tokens
	outputTokenCostPer1K = 0.015  // $15 per 1M output tokens
	avgOutputTokens      = 150    // Average output tokens per commit
)

var estimateCmd = &cobra.Command{
	Use:   "estimate",
	Short: "Estimate API costs for analyzing commits",
	Long: `Estimate the approximate API costs for analyzing git commits
without actually making any API calls.

Example:
  git-resume estimate --month=4
  git-resume estimate --from=2024-01-01 --to=2024-03-31`,
	RunE: runEstimate,
}

func init() {
	rootCmd.AddCommand(estimateCmd)

	estimateCmd.Flags().IntVarP(&month, "month", "m", 0, "month to estimate (1-12)")
	estimateCmd.Flags().IntVarP(&year, "year", "y", 0, "year to estimate (defaults to current year)")
	estimateCmd.Flags().StringVar(&fromDate, "from", "", "start date (YYYY-MM-DD)")
	estimateCmd.Flags().StringVar(&toDate, "to", "", "end date (YYYY-MM-DD)")
	estimateCmd.Flags().IntVar(&batchSize, "batch-size", 5, "number of commits per API call")
}

func runEstimate(cmd *cobra.Command, args []string) error {
	// Determine date range
	from, to, err := parseDateRange()
	if err != nil {
		return err
	}

	// Get repository path
	repoPath := viper.GetString("repo")
	if repoPath == "" {
		repoPath = viper.GetString("DEFAULT_REPO_PATH")
	}
	if repoPath == "" {
		repoPath, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	fmt.Printf("Estimating costs for: %s\n", repoPath)
	fmt.Printf("Date range: %s to %s\n\n", from.Format("2006-01-02"), to.Format("2006-01-02"))

	// Initialize git parser
	parser, err := git.NewParser(repoPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	// Get commits
	commits, err := parser.GetCommits(from, to)
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	fmt.Printf("Total commits found: %d\n", len(commits))

	if len(commits) == 0 {
		fmt.Println("No commits to analyze")
		return nil
	}

	// Check cache for already processed
	dbPath := viper.GetString("db")
	if dbPath == "" {
		dbPath = viper.GetString("DB_PATH")
	}
	if dbPath == "" {
		dbPath = "./data/cache.db"
	}

	var unprocessedCount int
	database, err := db.New(dbPath)
	if err == nil {
		defer database.Close()
		cache := db.NewCache(database)
		unprocessed, _ := cache.FilterUnprocessed(commits)
		unprocessedCount = len(unprocessed)
		commits = unprocessed
	} else {
		unprocessedCount = len(commits)
	}

	fmt.Printf("Already processed: %d\n", len(commits)-unprocessedCount+len(commits))
	fmt.Printf("New commits to process: %d\n\n", unprocessedCount)

	if unprocessedCount == 0 {
		fmt.Println("All commits have already been processed. No additional cost.")
		return nil
	}

	// Calculate estimates
	estimate := calculateEstimate(commits, batchSize)

	// Print estimate
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("           COST ESTIMATE")
	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("Commits to analyze:     %d\n", estimate.CommitCount)
	fmt.Printf("API calls needed:       %d\n", estimate.APICallCount)
	fmt.Printf("Batch size:             %d commits/call\n\n", batchSize)

	fmt.Println("Token Estimates:")
	fmt.Printf("  Input tokens:         ~%d\n", estimate.InputTokens)
	fmt.Printf("  Output tokens:        ~%d\n", estimate.OutputTokens)
	fmt.Printf("  Total tokens:         ~%d\n\n", estimate.InputTokens+estimate.OutputTokens)

	fmt.Println("Cost Breakdown:")
	fmt.Printf("  Input cost:           $%.4f\n", estimate.InputCost)
	fmt.Printf("  Output cost:          $%.4f\n", estimate.OutputCost)
	fmt.Printf("  ─────────────────────────────\n")
	fmt.Printf("  ESTIMATED TOTAL:      $%.4f\n", estimate.TotalCost)
	fmt.Println("═══════════════════════════════════════")

	// Provide context
	fmt.Println("\nNote: Actual costs may vary based on:")
	fmt.Println("  • Commit message complexity")
	fmt.Println("  • Number of files changed")
	fmt.Println("  • API response length")

	return nil
}

// CostEstimate holds the estimated costs
type CostEstimate struct {
	CommitCount   int
	APICallCount  int
	InputTokens   int
	OutputTokens  int
	InputCost     float64
	OutputCost    float64
	TotalCost     float64
}

func calculateEstimate(commits []models.Commit, batchSize int) CostEstimate {
	estimate := CostEstimate{
		CommitCount: len(commits),
	}

	// Calculate number of API calls
	estimate.APICallCount = (len(commits) + batchSize - 1) / batchSize

	// Estimate input tokens
	for i := 0; i < len(commits); i += batchSize {
		end := i + batchSize
		if end > len(commits) {
			end = len(commits)
		}

		batch := models.CommitBatch{
			Commits:  commits[i:end],
			Project:  "estimate",
			FromDate: time.Now(),
			ToDate:   time.Now(),
		}

		estimate.InputTokens += llm.EstimateTokens(batch)
	}

	// Estimate output tokens (based on average per commit)
	estimate.OutputTokens = len(commits) * avgOutputTokens

	// Calculate costs
	estimate.InputCost = float64(estimate.InputTokens) / 1000.0 * inputTokenCostPer1K
	estimate.OutputCost = float64(estimate.OutputTokens) / 1000.0 * outputTokenCostPer1K
	estimate.TotalCost = estimate.InputCost + estimate.OutputCost

	return estimate
}
