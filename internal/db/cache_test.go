package db

import (
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/wootaiklee/git-resume/pkg/models"
)

// newTestDB creates an in-memory SQLite database for testing
func newTestDB(t *testing.T) *DB {
	t.Helper()

	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	db := &DB{conn: conn}

	// Initialize schema
	if err := db.migrate(); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func TestCache_IsProcessed(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	cache := NewCache(db)

	tests := []struct {
		name           string
		hash           string
		markFirst      bool
		expectedResult bool
	}{
		{
			name:           "unprocessed commit",
			hash:           "abc1234",
			markFirst:      false,
			expectedResult: false,
		},
		{
			name:           "processed commit",
			hash:           "def5678",
			markFirst:      true,
			expectedResult: true,
		},
		{
			name:           "empty hash unprocessed",
			hash:           "",
			markFirst:      false,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.markFirst {
				if err := cache.MarkProcessed(tt.hash); err != nil {
					t.Fatalf("MarkProcessed() error = %v", err)
				}
			}

			processed, err := cache.IsProcessed(tt.hash)
			if err != nil {
				t.Fatalf("IsProcessed() error = %v", err)
			}

			if processed != tt.expectedResult {
				t.Errorf("IsProcessed(%q) = %v, expected %v", tt.hash, processed, tt.expectedResult)
			}
		})
	}
}

func TestCache_MarkProcessed(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	cache := NewCache(db)

	tests := []struct {
		name string
		hash string
	}{
		{
			name: "mark single commit",
			hash: "commit1234567890",
		},
		{
			name: "mark short hash",
			hash: "abc1234",
		},
		{
			name: "mark with special characters",
			hash: "abc_123-def",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mark as processed
			err := cache.MarkProcessed(tt.hash)
			if err != nil {
				t.Fatalf("MarkProcessed() error = %v", err)
			}

			// Verify it's processed
			processed, err := cache.IsProcessed(tt.hash)
			if err != nil {
				t.Fatalf("IsProcessed() error = %v", err)
			}

			if !processed {
				t.Errorf("hash %q should be marked as processed", tt.hash)
			}
		})
	}
}

func TestCache_MarkProcessed_Idempotent(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	cache := NewCache(db)
	hash := "test1234567"

	// Mark the same hash multiple times
	for i := 0; i < 3; i++ {
		err := cache.MarkProcessed(hash)
		if err != nil {
			t.Fatalf("MarkProcessed() iteration %d error = %v", i, err)
		}
	}

	// Verify it's still processed
	processed, err := cache.IsProcessed(hash)
	if err != nil {
		t.Fatalf("IsProcessed() error = %v", err)
	}

	if !processed {
		t.Error("hash should be marked as processed after multiple marks")
	}

	// Verify only one record exists
	var count int
	err = db.conn.QueryRow("SELECT COUNT(*) FROM processed_commits WHERE hash = ?", hash).Scan(&count)
	if err != nil {
		t.Fatalf("count query error = %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 record for hash, got %d", count)
	}
}

func TestCache_FilterUnprocessed(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	cache := NewCache(db)

	// Create test commits
	commits := []models.Commit{
		{Hash: "hash1", Message: "Commit 1"},
		{Hash: "hash2", Message: "Commit 2"},
		{Hash: "hash3", Message: "Commit 3"},
		{Hash: "hash4", Message: "Commit 4"},
		{Hash: "hash5", Message: "Commit 5"},
	}

	// Mark some as processed
	processedHashes := []string{"hash2", "hash4"}
	for _, hash := range processedHashes {
		if err := cache.MarkProcessed(hash); err != nil {
			t.Fatalf("MarkProcessed() error = %v", err)
		}
	}

	// Filter unprocessed
	unprocessed, err := cache.FilterUnprocessed(commits)
	if err != nil {
		t.Fatalf("FilterUnprocessed() error = %v", err)
	}

	// Expected unprocessed: hash1, hash3, hash5
	expectedHashes := map[string]bool{
		"hash1": true,
		"hash3": true,
		"hash5": true,
	}

	if len(unprocessed) != len(expectedHashes) {
		t.Errorf("FilterUnprocessed() returned %d commits, expected %d", len(unprocessed), len(expectedHashes))
	}

	for _, commit := range unprocessed {
		if !expectedHashes[commit.Hash] {
			t.Errorf("unexpected commit hash in result: %s", commit.Hash)
		}
	}
}

func TestCache_FilterUnprocessed_EmptyInput(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	cache := NewCache(db)

	unprocessed, err := cache.FilterUnprocessed([]models.Commit{})
	if err != nil {
		t.Fatalf("FilterUnprocessed() error = %v", err)
	}

	if len(unprocessed) != 0 {
		t.Errorf("FilterUnprocessed(empty) should return empty slice, got %d", len(unprocessed))
	}
}

func TestCache_FilterUnprocessed_AllProcessed(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	cache := NewCache(db)

	commits := []models.Commit{
		{Hash: "hash1", Message: "Commit 1"},
		{Hash: "hash2", Message: "Commit 2"},
	}

	// Mark all as processed
	for _, commit := range commits {
		if err := cache.MarkProcessed(commit.Hash); err != nil {
			t.Fatalf("MarkProcessed() error = %v", err)
		}
	}

	unprocessed, err := cache.FilterUnprocessed(commits)
	if err != nil {
		t.Fatalf("FilterUnprocessed() error = %v", err)
	}

	if len(unprocessed) != 0 {
		t.Errorf("FilterUnprocessed() should return empty when all processed, got %d", len(unprocessed))
	}
}

func TestCache_FilterUnprocessed_NoneProcessed(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	cache := NewCache(db)

	commits := []models.Commit{
		{Hash: "hash1", Message: "Commit 1"},
		{Hash: "hash2", Message: "Commit 2"},
		{Hash: "hash3", Message: "Commit 3"},
	}

	unprocessed, err := cache.FilterUnprocessed(commits)
	if err != nil {
		t.Fatalf("FilterUnprocessed() error = %v", err)
	}

	if len(unprocessed) != len(commits) {
		t.Errorf("FilterUnprocessed() should return all commits when none processed, got %d", len(unprocessed))
	}
}

func TestCache_SaveResult(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	cache := NewCache(db)

	result := models.AnalysisResult{
		CommitHash:    "abc1234567890",
		Date:          time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		Project:       "test-project",
		Category:      models.CategoryFeature,
		ImpactSummary: "Implemented new authentication module with OAuth2 support",
		CreatedAt:     time.Now(),
	}

	err := cache.SaveResult(result)
	if err != nil {
		t.Fatalf("SaveResult() error = %v", err)
	}

	// Verify result was saved
	var count int
	err = db.conn.QueryRow("SELECT COUNT(*) FROM analysis_results WHERE commit_hash = ?", result.CommitHash).Scan(&count)
	if err != nil {
		t.Fatalf("count query error = %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 result record, got %d", count)
	}

	// Verify commit was marked as processed
	processed, err := cache.IsProcessed(result.CommitHash)
	if err != nil {
		t.Fatalf("IsProcessed() error = %v", err)
	}

	if !processed {
		t.Error("commit should be marked as processed after SaveResult")
	}
}

func TestCache_SaveResults(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	cache := NewCache(db)

	results := []models.AnalysisResult{
		{
			CommitHash:    "hash1",
			Date:          time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Project:       "project-a",
			Category:      models.CategoryFeature,
			ImpactSummary: "Feature 1 description",
			CreatedAt:     time.Now(),
		},
		{
			CommitHash:    "hash2",
			Date:          time.Date(2024, 1, 16, 11, 0, 0, 0, time.UTC),
			Project:       "project-a",
			Category:      models.CategoryFix,
			ImpactSummary: "Bug fix description",
			CreatedAt:     time.Now(),
		},
		{
			CommitHash:    "hash3",
			Date:          time.Date(2024, 1, 17, 12, 0, 0, 0, time.UTC),
			Project:       "project-b",
			Category:      models.CategoryRefactor,
			ImpactSummary: "Refactor description",
			CreatedAt:     time.Now(),
		},
	}

	err := cache.SaveResults(results)
	if err != nil {
		t.Fatalf("SaveResults() error = %v", err)
	}

	// Verify all results were saved
	var count int
	err = db.conn.QueryRow("SELECT COUNT(*) FROM analysis_results").Scan(&count)
	if err != nil {
		t.Fatalf("count query error = %v", err)
	}

	if count != len(results) {
		t.Errorf("expected %d results, got %d", len(results), count)
	}

	// Verify all commits were marked as processed
	for _, result := range results {
		processed, err := cache.IsProcessed(result.CommitHash)
		if err != nil {
			t.Fatalf("IsProcessed() error = %v", err)
		}

		if !processed {
			t.Errorf("commit %s should be marked as processed", result.CommitHash)
		}
	}
}

func TestCache_SaveResults_Empty(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	cache := NewCache(db)

	err := cache.SaveResults([]models.AnalysisResult{})
	if err != nil {
		t.Fatalf("SaveResults(empty) error = %v", err)
	}

	// Verify no results were saved
	var count int
	err = db.conn.QueryRow("SELECT COUNT(*) FROM analysis_results").Scan(&count)
	if err != nil {
		t.Fatalf("count query error = %v", err)
	}

	if count != 0 {
		t.Errorf("expected 0 results, got %d", count)
	}
}

func TestCache_GetResults(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	cache := NewCache(db)

	// Insert test data
	results := []models.AnalysisResult{
		{
			CommitHash:    "hash1",
			Date:          time.Date(2024, 1, 10, 10, 0, 0, 0, time.UTC),
			Project:       "project-a",
			Category:      models.CategoryFeature,
			ImpactSummary: "Feature 1",
		},
		{
			CommitHash:    "hash2",
			Date:          time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
			Project:       "project-a",
			Category:      models.CategoryFix,
			ImpactSummary: "Fix 1",
		},
		{
			CommitHash:    "hash3",
			Date:          time.Date(2024, 1, 20, 12, 0, 0, 0, time.UTC),
			Project:       "project-b",
			Category:      models.CategoryRefactor,
			ImpactSummary: "Refactor 1",
		},
		{
			CommitHash:    "hash4",
			Date:          time.Date(2024, 2, 5, 13, 0, 0, 0, time.UTC),
			Project:       "project-a",
			Category:      models.CategoryFeature,
			ImpactSummary: "Feature 2",
		},
	}

	if err := cache.SaveResults(results); err != nil {
		t.Fatalf("SaveResults() error = %v", err)
	}

	tests := []struct {
		name          string
		project       string
		from          time.Time
		to            time.Time
		expectedCount int
	}{
		{
			name:          "all results within date range",
			project:       "",
			from:          time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			to:            time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
			expectedCount: 4,
		},
		{
			name:          "filter by project",
			project:       "project-a",
			from:          time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			to:            time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
			expectedCount: 3,
		},
		{
			name:          "filter by date range",
			project:       "",
			from:          time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
			to:            time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
			expectedCount: 2,
		},
		{
			name:          "filter by project and date range",
			project:       "project-a",
			from:          time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			to:            time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			expectedCount: 2,
		},
		{
			name:          "no results in date range",
			project:       "",
			from:          time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			to:            time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			expectedCount: 0,
		},
		{
			name:          "non-existent project",
			project:       "non-existent",
			from:          time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			to:            time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cache.GetResults(tt.project, tt.from, tt.to)
			if err != nil {
				t.Fatalf("GetResults() error = %v", err)
			}

			if len(got) != tt.expectedCount {
				t.Errorf("GetResults() returned %d results, expected %d", len(got), tt.expectedCount)
			}
		})
	}
}

func TestCache_GetResults_OrderedByDateDesc(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	cache := NewCache(db)

	results := []models.AnalysisResult{
		{CommitHash: "hash1", Date: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC), Project: "p", Category: models.CategoryFeature, ImpactSummary: "s1"},
		{CommitHash: "hash2", Date: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC), Project: "p", Category: models.CategoryFeature, ImpactSummary: "s2"},
		{CommitHash: "hash3", Date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), Project: "p", Category: models.CategoryFeature, ImpactSummary: "s3"},
	}

	if err := cache.SaveResults(results); err != nil {
		t.Fatalf("SaveResults() error = %v", err)
	}

	got, err := cache.GetResults("", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("GetResults() error = %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("GetResults() returned %d results, expected 3", len(got))
	}

	// Verify descending order by date
	expectedOrder := []string{"hash2", "hash3", "hash1"} // newest to oldest
	for i, hash := range expectedOrder {
		if got[i].CommitHash != hash {
			t.Errorf("result[%d].CommitHash = %s, expected %s", i, got[i].CommitHash, hash)
		}
	}
}

func TestCache_RecordTokenUsage(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	cache := NewCache(db)

	err := cache.RecordTokenUsage("batch-123", 1000, 500, 0.05)
	if err != nil {
		t.Fatalf("RecordTokenUsage() error = %v", err)
	}

	// Verify token usage was recorded
	var count int
	err = db.conn.QueryRow("SELECT COUNT(*) FROM token_usage WHERE batch_id = ?", "batch-123").Scan(&count)
	if err != nil {
		t.Fatalf("count query error = %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 token usage record, got %d", count)
	}
}

func TestCache_GetTotalTokenUsage(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	cache := NewCache(db)

	// Record multiple batches
	batches := []struct {
		batchID      string
		inputTokens  int
		outputTokens int
		costEstimate float64
	}{
		{"batch-1", 1000, 500, 0.05},
		{"batch-2", 2000, 800, 0.10},
		{"batch-3", 1500, 600, 0.07},
	}

	for _, b := range batches {
		if err := cache.RecordTokenUsage(b.batchID, b.inputTokens, b.outputTokens, b.costEstimate); err != nil {
			t.Fatalf("RecordTokenUsage() error = %v", err)
		}
	}

	inputTokens, outputTokens, totalCost, err := cache.GetTotalTokenUsage()
	if err != nil {
		t.Fatalf("GetTotalTokenUsage() error = %v", err)
	}

	expectedInput := 1000 + 2000 + 1500
	expectedOutput := 500 + 800 + 600
	expectedCost := 0.05 + 0.10 + 0.07

	if inputTokens != expectedInput {
		t.Errorf("inputTokens = %d, expected %d", inputTokens, expectedInput)
	}

	if outputTokens != expectedOutput {
		t.Errorf("outputTokens = %d, expected %d", outputTokens, expectedOutput)
	}

	// Compare float with tolerance
	if totalCost < expectedCost-0.001 || totalCost > expectedCost+0.001 {
		t.Errorf("totalCost = %f, expected %f", totalCost, expectedCost)
	}
}

func TestCache_GetTotalTokenUsage_Empty(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	cache := NewCache(db)

	inputTokens, outputTokens, totalCost, err := cache.GetTotalTokenUsage()
	if err != nil {
		t.Fatalf("GetTotalTokenUsage() error = %v", err)
	}

	if inputTokens != 0 {
		t.Errorf("inputTokens = %d, expected 0", inputTokens)
	}

	if outputTokens != 0 {
		t.Errorf("outputTokens = %d, expected 0", outputTokens)
	}

	if totalCost != 0 {
		t.Errorf("totalCost = %f, expected 0", totalCost)
	}
}

func TestNewCache(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	cache := NewCache(db)

	if cache == nil {
		t.Fatal("NewCache() returned nil")
	}

	if cache.db != db {
		t.Error("NewCache() did not set db correctly")
	}
}

// Benchmark tests
func BenchmarkCache_IsProcessed(b *testing.B) {
	conn, _ := sql.Open("sqlite", ":memory:")
	db := &DB{conn: conn}
	db.migrate()
	defer db.Close()

	cache := NewCache(db)
	cache.MarkProcessed("benchmark-hash")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.IsProcessed("benchmark-hash")
	}
}

func BenchmarkCache_FilterUnprocessed(b *testing.B) {
	conn, _ := sql.Open("sqlite", ":memory:")
	db := &DB{conn: conn}
	db.migrate()
	defer db.Close()

	cache := NewCache(db)

	// Create 100 commits, mark 50 as processed
	commits := make([]models.Commit, 100)
	for i := 0; i < 100; i++ {
		commits[i] = models.Commit{Hash: "hash" + string(rune('a'+i%26)) + string(rune('0'+i/26))}
		if i%2 == 0 {
			cache.MarkProcessed(commits[i].Hash)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.FilterUnprocessed(commits)
	}
}
