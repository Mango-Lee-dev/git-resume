package llm

import (
	"strings"
	"testing"
	"time"

	"github.com/wootaiklee/git-resume/pkg/models"
)

func TestBuildPrompt(t *testing.T) {
	tests := []struct {
		name               string
		batch              models.CommitBatch
		expectedSubstrings []string
		notExpected        []string
	}{
		{
			name: "single commit batch",
			batch: models.CommitBatch{
				Project: "test-project",
				Commits: []models.Commit{
					{
						Hash:      "abc1234567890def",
						Date:      time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
						Message:   "Add user authentication feature",
						Files:     []string{"auth/handler.go", "auth/service.go"},
						Additions: 150,
						Deletions: 10,
					},
				},
			},
			expectedSubstrings: []string{
				"test-project",
				"abc1234", // First 7 chars of hash
				"2024-01-15",
				"Add user authentication feature",
				"auth/handler.go",
				"auth/service.go",
				"+150/-10",
				"2 (+150/-10 lines)", // Files changed count
				"Commit 1",
				"JSON format",
			},
			notExpected: []string{},
		},
		{
			name: "multiple commits batch",
			batch: models.CommitBatch{
				Project: "multi-commit-project",
				Commits: []models.Commit{
					{
						Hash:      "hash1_abc",
						Date:      time.Date(2024, 2, 1, 9, 0, 0, 0, time.UTC),
						Message:   "First commit message",
						Files:     []string{"file1.go"},
						Additions: 50,
						Deletions: 5,
					},
					{
						Hash:      "hash2_def",
						Date:      time.Date(2024, 2, 2, 10, 0, 0, 0, time.UTC),
						Message:   "Second commit message",
						Files:     []string{"file2.go", "file3.go"},
						Additions: 100,
						Deletions: 20,
					},
				},
			},
			expectedSubstrings: []string{
				"multi-commit-project",
				"Commit 1",
				"Commit 2",
				"hash1_a", // First 7 chars
				"hash2_d", // First 7 chars
				"First commit message",
				"Second commit message",
			},
			notExpected: []string{},
		},
		{
			name: "batch with many files truncated",
			batch: models.CommitBatch{
				Project: "large-project",
				Commits: []models.Commit{
					{
						Hash:    "largehash123",
						Date:    time.Date(2024, 3, 1, 12, 0, 0, 0, time.UTC),
						Message: "Large commit with many files",
						Files: []string{
							"file1.go", "file2.go", "file3.go", "file4.go", "file5.go",
							"file6.go", "file7.go", "file8.go", "file9.go", "file10.go",
							"file11.go", "file12.go", "file13.go",
						},
						Additions: 500,
						Deletions: 100,
					},
				},
			},
			expectedSubstrings: []string{
				"file1.go",
				"file10.go",
				"and 3 more files", // 13 - 10 = 3 more files
			},
			notExpected: []string{
				"file11.go", // Should be truncated
				"file12.go",
				"file13.go",
			},
		},
		{
			name: "commit with empty files",
			batch: models.CommitBatch{
				Project: "empty-files-project",
				Commits: []models.Commit{
					{
						Hash:      "emptyhash",
						Date:      time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
						Message:   "Commit with no files",
						Files:     []string{},
						Additions: 0,
						Deletions: 0,
					},
				},
			},
			expectedSubstrings: []string{
				"empty-files-project",
				"Commit with no files",
				"0 (+0/-0 lines)",
			},
			notExpected: []string{
				"Key files:", // Should not appear when files is empty
			},
		},
		{
			name: "batch with whitespace in message",
			batch: models.CommitBatch{
				Project: "whitespace-project",
				Commits: []models.Commit{
					{
						Hash:      "wshash123",
						Date:      time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
						Message:   "  Message with leading/trailing whitespace  \n\n",
						Files:     []string{"file.go"},
						Additions: 10,
						Deletions: 5,
					},
				},
			},
			expectedSubstrings: []string{
				"Message with leading/trailing whitespace", // Should be trimmed
			},
			notExpected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := BuildPrompt(tt.batch)

			for _, substr := range tt.expectedSubstrings {
				if !strings.Contains(prompt, substr) {
					t.Errorf("BuildPrompt() output should contain %q\nGot:\n%s", substr, prompt)
				}
			}

			for _, substr := range tt.notExpected {
				if strings.Contains(prompt, substr) {
					t.Errorf("BuildPrompt() output should NOT contain %q\nGot:\n%s", substr, prompt)
				}
			}
		})
	}
}

func TestBuildPrompt_Format(t *testing.T) {
	batch := models.CommitBatch{
		Project: "format-test",
		Commits: []models.Commit{
			{
				Hash:      "formathash1234567",
				Date:      time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC),
				Message:   "Test commit for format validation",
				Files:     []string{"main.go", "util.go"},
				Additions: 75,
				Deletions: 25,
			},
		},
	}

	prompt := BuildPrompt(batch)

	// Verify proper structure
	if !strings.HasPrefix(prompt, "Analyze the following git commits") {
		t.Error("Prompt should start with 'Analyze the following git commits'")
	}

	if !strings.HasSuffix(prompt, "Generate resume bullet points in JSON format as specified.") {
		t.Error("Prompt should end with JSON format instruction")
	}

	// Verify commit delimiter format
	if !strings.Contains(prompt, "--- Commit 1 ---") {
		t.Error("Prompt should contain commit delimiter")
	}

	// Verify hash is truncated to 7 characters
	if strings.Contains(prompt, "formathash1234567") {
		t.Error("Full hash should not appear, only first 7 characters")
	}
	if !strings.Contains(prompt, "formath") { // First 7 chars: f-o-r-m-a-t-h
		t.Error("First 7 characters of hash should appear")
	}
}

func TestParseResponse_ValidJSON(t *testing.T) {
	tests := []struct {
		name          string
		responseText  string
		batch         models.CommitBatch
		expectedCount int
		checkFirst    func(t *testing.T, result models.AnalysisResult)
	}{
		{
			name: "simple valid JSON array",
			responseText: `Here are the results:
[
  {"hash": "abc1234", "category": "Feature", "summary": "Implemented OAuth2 authentication system"}
]`,
			batch: models.CommitBatch{
				Project: "test-project",
				Commits: []models.Commit{
					{
						Hash: "abc1234567890",
						Date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			expectedCount: 1,
			checkFirst: func(t *testing.T, result models.AnalysisResult) {
				if result.Category != models.CategoryFeature {
					t.Errorf("expected category Feature, got %s", result.Category)
				}
				if result.ImpactSummary != "Implemented OAuth2 authentication system" {
					t.Errorf("unexpected summary: %s", result.ImpactSummary)
				}
				if result.Project != "test-project" {
					t.Errorf("expected project test-project, got %s", result.Project)
				}
			},
		},
		{
			name: "multiple results",
			responseText: `[
  {"hash": "hash1_a", "category": "Feature", "summary": "Added new feature"},
  {"hash": "hash2_b", "category": "Fix", "summary": "Fixed critical bug"},
  {"hash": "hash3_c", "category": "Refactor", "summary": "Refactored code"}
]`,
			batch: models.CommitBatch{
				Project: "multi-project",
				Commits: []models.Commit{
					{Hash: "hash1_abc123", Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
					{Hash: "hash2_bcd234", Date: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
					{Hash: "hash3_cde345", Date: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectedCount: 3,
			checkFirst:    nil,
		},
		{
			name: "JSON with surrounding text",
			responseText: `Based on my analysis, here are the resume bullet points:

[{"hash": "def5678", "category": "Fix", "summary": "Resolved memory leak issue"}]

I hope this helps with your resume!`,
			batch: models.CommitBatch{
				Project: "text-surrounding-project",
				Commits: []models.Commit{
					{Hash: "def5678901234", Date: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectedCount: 1,
			checkFirst: func(t *testing.T, result models.AnalysisResult) {
				if result.Category != models.CategoryFix {
					t.Errorf("expected category Fix, got %s", result.Category)
				}
			},
		},
		{
			name: "all category types",
			responseText: `[
  {"hash": "hash_fe", "category": "Feature", "summary": "Feature work"},
  {"hash": "hash_fi", "category": "fix", "summary": "Bug fix"},
  {"hash": "hash_re", "category": "REFACTOR", "summary": "Code refactor"},
  {"hash": "hash_do", "category": "Docs", "summary": "Documentation"},
  {"hash": "hash_te", "category": "Test", "summary": "Test coverage"},
  {"hash": "hash_ch", "category": "Chore", "summary": "Maintenance"}
]`,
			batch: models.CommitBatch{
				Project: "all-categories",
				Commits: []models.Commit{
					{Hash: "hash_fe_1234", Date: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)},
					{Hash: "hash_fi_1234", Date: time.Date(2024, 3, 2, 0, 0, 0, 0, time.UTC)},
					{Hash: "hash_re_1234", Date: time.Date(2024, 3, 3, 0, 0, 0, 0, time.UTC)},
					{Hash: "hash_do_1234", Date: time.Date(2024, 3, 4, 0, 0, 0, 0, time.UTC)},
					{Hash: "hash_te_1234", Date: time.Date(2024, 3, 5, 0, 0, 0, 0, time.UTC)},
					{Hash: "hash_ch_1234", Date: time.Date(2024, 3, 6, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectedCount: 6,
			checkFirst:    nil,
		},
		{
			name: "unknown category defaults to chore",
			responseText: `[{"hash": "unknown", "category": "Unknown", "summary": "Some work"}]`,
			batch: models.CommitBatch{
				Project: "unknown-cat-project",
				Commits: []models.Commit{
					{Hash: "unknown1234567", Date: time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectedCount: 1,
			checkFirst: func(t *testing.T, result models.AnalysisResult) {
				if result.Category != models.CategoryChore {
					t.Errorf("unknown category should default to Chore, got %s", result.Category)
				}
			},
		},
		{
			name: "hash mismatch - result skipped",
			responseText: `[
  {"hash": "match12", "category": "Feature", "summary": "Matching hash"},
  {"hash": "nomatch", "category": "Fix", "summary": "Non-matching hash"}
]`,
			batch: models.CommitBatch{
				Project: "mismatch-project",
				Commits: []models.Commit{
					{Hash: "match1234567890", Date: time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)},
				},
			},
			expectedCount: 1, // Only the matching one
			checkFirst: func(t *testing.T, result models.AnalysisResult) {
				if result.CommitHash != "match1234567890" {
					t.Errorf("expected full hash match1234567890, got %s", result.CommitHash)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := ParseResponse(tt.responseText, tt.batch)
			if err != nil {
				t.Fatalf("ParseResponse() unexpected error = %v", err)
			}

			if len(results) != tt.expectedCount {
				t.Errorf("ParseResponse() returned %d results, expected %d", len(results), tt.expectedCount)
			}

			if tt.checkFirst != nil && len(results) > 0 {
				tt.checkFirst(t, results[0])
			}
		})
	}
}

func TestParseResponse_InvalidJSON(t *testing.T) {
	batch := models.CommitBatch{
		Project: "test-project",
		Commits: []models.Commit{
			{Hash: "abc1234567890", Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	tests := []struct {
		name         string
		responseText string
		expectError  bool
	}{
		{
			name:         "no JSON array",
			responseText: "This is a response without any JSON",
			expectError:  true,
		},
		{
			name:         "empty string",
			responseText: "",
			expectError:  true,
		},
		{
			name:         "only opening bracket",
			responseText: "[",
			expectError:  true,
		},
		{
			name:         "only closing bracket",
			responseText: "]",
			expectError:  true,
		},
		{
			name:         "malformed JSON",
			responseText: "[{\"hash\": \"abc1234\", \"category\": \"Feature\", \"summary\": }]",
			expectError:  true,
		},
		{
			name:         "JSON object not array",
			responseText: "{\"hash\": \"abc1234\", \"category\": \"Feature\", \"summary\": \"test\"}",
			expectError:  true,
		},
		{
			name:         "brackets in wrong order",
			responseText: "]some text[",
			expectError:  true,
		},
		{
			name:         "incomplete JSON structure",
			responseText: "[{\"hash\": \"abc1234\"",
			expectError:  true,
		},
		{
			name:         "nested array confusing parser",
			responseText: "text [[nested]] text",
			expectError:  true, // Inner array is not valid JSON
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := ParseResponse(tt.responseText, batch)

			if tt.expectError {
				if err == nil {
					t.Errorf("ParseResponse() expected error for %q, got results: %v", tt.responseText, results)
				}
			} else {
				if err != nil {
					t.Errorf("ParseResponse() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestParseResponse_EmptyArray(t *testing.T) {
	batch := models.CommitBatch{
		Project: "test-project",
		Commits: []models.Commit{
			{Hash: "abc1234567890", Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	results, err := ParseResponse("[]", batch)
	if err != nil {
		t.Fatalf("ParseResponse() unexpected error for empty array = %v", err)
	}

	if len(results) != 0 {
		t.Errorf("ParseResponse([]) should return empty slice, got %d results", len(results))
	}
}

func TestParseResponse_PreservesCommitData(t *testing.T) {
	commitDate := time.Date(2024, 6, 15, 14, 30, 45, 0, time.UTC)
	batch := models.CommitBatch{
		Project: "preserve-data-project",
		Commits: []models.Commit{
			{
				Hash:      "preserve1234567",
				Date:      commitDate,
				Message:   "Original message",
				Files:     []string{"file1.go", "file2.go"},
				Additions: 100,
				Deletions: 50,
			},
		},
	}

	responseText := `[{"hash": "preserv", "category": "Feature", "summary": "Test summary"}]`

	results, err := ParseResponse(responseText, batch)
	if err != nil {
		t.Fatalf("ParseResponse() error = %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	result := results[0]

	// Verify data is preserved from original commit
	if result.CommitHash != "preserve1234567" {
		t.Errorf("CommitHash should be full hash, got %s", result.CommitHash)
	}

	if !result.Date.Equal(commitDate) {
		t.Errorf("Date should match original commit date, got %v", result.Date)
	}

	if result.Project != "preserve-data-project" {
		t.Errorf("Project should match batch project, got %s", result.Project)
	}

	// CreatedAt should be set (not zero)
	if result.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
}

func TestParseCategory(t *testing.T) {
	tests := []struct {
		input    string
		expected models.Category
	}{
		{"Feature", models.CategoryFeature},
		{"feature", models.CategoryFeature},
		{"FEATURE", models.CategoryFeature},
		{"Fix", models.CategoryFix},
		{"fix", models.CategoryFix},
		{"FIX", models.CategoryFix},
		{"Refactor", models.CategoryRefactor},
		{"refactor", models.CategoryRefactor},
		{"Docs", models.CategoryDocs},
		{"docs", models.CategoryDocs},
		{"Test", models.CategoryTest},
		{"test", models.CategoryTest},
		{"Chore", models.CategoryChore},
		{"chore", models.CategoryChore},
		// Unknown defaults to Chore
		{"Unknown", models.CategoryChore},
		{"random", models.CategoryChore},
		{"", models.CategoryChore},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseCategory(tt.input)
			if result != tt.expected {
				t.Errorf("parseCategory(%q) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSystemPrompt_Content(t *testing.T) {
	// Verify system prompt contains essential instructions
	requiredPhrases := []string{
		"technical writer",
		"resume",
		"STAR method",
		"business impact",
		"action verbs",
		"JSON",
		"hash",
		"category",
		"summary",
		"Feature",
		"Fix",
		"Refactor",
		"Docs",
		"Test",
		"Chore",
	}

	for _, phrase := range requiredPhrases {
		if !strings.Contains(SystemPrompt, phrase) {
			t.Errorf("SystemPrompt should contain %q", phrase)
		}
	}
}

// Benchmark tests
func BenchmarkBuildPrompt(b *testing.B) {
	batch := models.CommitBatch{
		Project: "benchmark-project",
		Commits: make([]models.Commit, 10),
	}

	for i := 0; i < 10; i++ {
		batch.Commits[i] = models.Commit{
			Hash:      "benchmarkhash" + string(rune('0'+i)),
			Date:      time.Now(),
			Message:   "Benchmark commit message for testing performance",
			Files:     []string{"file1.go", "file2.go", "file3.go"},
			Additions: 100,
			Deletions: 50,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildPrompt(batch)
	}
}

func BenchmarkParseResponse(b *testing.B) {
	responseText := `[
		{"hash": "hash1_a", "category": "Feature", "summary": "Feature 1"},
		{"hash": "hash2_a", "category": "Fix", "summary": "Fix 1"},
		{"hash": "hash3_a", "category": "Refactor", "summary": "Refactor 1"},
		{"hash": "hash4_a", "category": "Docs", "summary": "Docs 1"},
		{"hash": "hash5_a", "category": "Test", "summary": "Test 1"}
	]`

	batch := models.CommitBatch{
		Project: "benchmark-project",
		Commits: []models.Commit{
			{Hash: "hash1_abc123", Date: time.Now()},
			{Hash: "hash2_abc123", Date: time.Now()},
			{Hash: "hash3_abc123", Date: time.Now()},
			{Hash: "hash4_abc123", Date: time.Now()},
			{Hash: "hash5_abc123", Date: time.Now()},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseResponse(responseText, batch)
	}
}
