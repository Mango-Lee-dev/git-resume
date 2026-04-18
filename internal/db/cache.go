package db

import (
	"time"

	"github.com/wootaiklee/git-resume/pkg/models"
)

// Cache handles commit and result caching operations
type Cache struct {
	db *DB
}

// NewCache creates a new cache instance
func NewCache(db *DB) *Cache {
	return &Cache{db: db}
}

// IsProcessed checks if a commit hash has already been processed
func (c *Cache) IsProcessed(hash string) (bool, error) {
	var count int
	err := c.db.conn.QueryRow(
		"SELECT COUNT(*) FROM processed_commits WHERE hash = ?",
		hash,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// MarkProcessed marks a commit as processed
func (c *Cache) MarkProcessed(hash string) error {
	_, err := c.db.conn.Exec(
		"INSERT OR IGNORE INTO processed_commits (hash) VALUES (?)",
		hash,
	)
	return err
}

// FilterUnprocessed returns only commits that haven't been processed
func (c *Cache) FilterUnprocessed(commits []models.Commit) ([]models.Commit, error) {
	var unprocessed []models.Commit

	for _, commit := range commits {
		processed, err := c.IsProcessed(commit.Hash)
		if err != nil {
			return nil, err
		}
		if !processed {
			unprocessed = append(unprocessed, commit)
		}
	}

	return unprocessed, nil
}

// SaveResult saves an analysis result to the database
func (c *Cache) SaveResult(result models.AnalysisResult) error {
	_, err := c.db.conn.Exec(`
		INSERT INTO analysis_results (commit_hash, date, project, category, impact_summary)
		VALUES (?, ?, ?, ?, ?)
	`, result.CommitHash, result.Date, result.Project, result.Category, result.ImpactSummary)

	if err != nil {
		return err
	}

	// Mark commit as processed
	return c.MarkProcessed(result.CommitHash)
}

// SaveResults saves multiple analysis results
func (c *Cache) SaveResults(results []models.AnalysisResult) error {
	tx, err := c.db.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO analysis_results (commit_hash, date, project, category, impact_summary)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	commitStmt, err := tx.Prepare("INSERT OR IGNORE INTO processed_commits (hash) VALUES (?)")
	if err != nil {
		return err
	}
	defer commitStmt.Close()

	for _, result := range results {
		_, err = stmt.Exec(result.CommitHash, result.Date, result.Project, result.Category, result.ImpactSummary)
		if err != nil {
			return err
		}
		_, err = commitStmt.Exec(result.CommitHash)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetResults retrieves analysis results with optional filters
func (c *Cache) GetResults(project string, from, to time.Time) ([]models.AnalysisResult, error) {
	query := `
		SELECT id, commit_hash, date, project, category, impact_summary, created_at
		FROM analysis_results
		WHERE date >= ? AND date <= ?
	`
	args := []interface{}{from, to}

	if project != "" {
		query += " AND project = ?"
		args = append(args, project)
	}

	query += " ORDER BY date DESC"

	rows, err := c.db.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.AnalysisResult
	for rows.Next() {
		var r models.AnalysisResult
		err := rows.Scan(&r.ID, &r.CommitHash, &r.Date, &r.Project, &r.Category, &r.ImpactSummary, &r.CreatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, rows.Err()
}

// RecordTokenUsage records API token usage for cost tracking
func (c *Cache) RecordTokenUsage(batchID string, inputTokens, outputTokens int, costEstimate float64) error {
	_, err := c.db.conn.Exec(`
		INSERT INTO token_usage (batch_id, input_tokens, output_tokens, cost_estimate)
		VALUES (?, ?, ?, ?)
	`, batchID, inputTokens, outputTokens, costEstimate)
	return err
}

// GetTotalTokenUsage returns total token usage statistics
func (c *Cache) GetTotalTokenUsage() (inputTokens, outputTokens int, totalCost float64, err error) {
	err = c.db.conn.QueryRow(`
		SELECT COALESCE(SUM(input_tokens), 0), COALESCE(SUM(output_tokens), 0), COALESCE(SUM(cost_estimate), 0)
		FROM token_usage
	`).Scan(&inputTokens, &outputTokens, &totalCost)
	return
}

// GetAllResults retrieves all analysis results from the database
func (c *Cache) GetAllResults() ([]models.AnalysisResult, error) {
	rows, err := c.db.conn.Query(`
		SELECT id, commit_hash, date, project, category, impact_summary, created_at
		FROM analysis_results
		ORDER BY date DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.AnalysisResult
	for rows.Next() {
		var r models.AnalysisResult
		err := rows.Scan(&r.ID, &r.CommitHash, &r.Date, &r.Project, &r.Category, &r.ImpactSummary, &r.CreatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, rows.Err()
}
