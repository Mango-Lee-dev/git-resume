package git

import (
	"testing"

	"github.com/wootaiklee/git-resume/pkg/models"
)

func TestScorer_Calculate(t *testing.T) {
	scorer := NewScorer()

	tests := []struct {
		name     string
		commit   models.Commit
		minScore int
		maxScore int
	}{
		// High-value commits with conventional commit format
		{
			name: "feature commit with good description",
			commit: models.Commit{
				Message:   "feat(auth): implement OAuth2 authentication with Google and GitHub providers",
				Files:     []string{"internal/auth/oauth.go", "internal/auth/providers.go", "internal/auth/handler.go"},
				Additions: 250,
				Deletions: 10,
			},
			minScore: 50,
			maxScore: 100,
		},
		{
			name: "fix commit with security mention",
			commit: models.Commit{
				Message:   "fix(api): resolve critical security vulnerability in JWT validation",
				Files:     []string{"internal/api/middleware.go", "internal/auth/jwt.go"},
				Additions: 100,
				Deletions: 50,
			},
			minScore: 50,
			maxScore: 100,
		},
		{
			name: "performance improvement commit",
			commit: models.Commit{
				Message:   "perf(db): optimize database query performance with index tuning",
				Files:     []string{"internal/db/queries.go", "migrations/002_add_indexes.sql"},
				Additions: 80,
				Deletions: 20,
			},
			minScore: 45,
			maxScore: 100,
		},
		// Implementation commits
		{
			name: "implement new feature",
			commit: models.Commit{
				Message:   "Implement real-time notification system using WebSocket",
				Files:     []string{"internal/ws/handler.go", "internal/ws/hub.go", "internal/ws/client.go"},
				Additions: 300,
				Deletions: 0,
			},
			minScore: 50,
			maxScore: 100,
		},
		// Refactoring commits
		{
			name: "refactoring commit",
			commit: models.Commit{
				Message:   "Refactor user service to use clean architecture patterns",
				Files:     []string{"internal/user/service.go", "internal/user/repository.go"},
				Additions: 150,
				Deletions: 100,
			},
			minScore: 40,
			maxScore: 100,
		},
		// Migration commits
		{
			name: "database migration commit",
			commit: models.Commit{
				Message:   "Migrate from PostgreSQL to MySQL with schema changes",
				Files:     []string{"internal/db/connection.go", "migrations/001_init.sql"},
				Additions: 200,
				Deletions: 150,
			},
			minScore: 45,
			maxScore: 100,
		},
		// API changes
		{
			name: "API endpoint commit",
			commit: models.Commit{
				Message:   "Add new REST API endpoints for user management",
				Files:     []string{"internal/api/users.go", "internal/api/routes.go"},
				Additions: 180,
				Deletions: 20,
			},
			minScore: 40,
			maxScore: 100,
		},
		// Low-value commits
		{
			name: "very short message single file",
			commit: models.Commit{
				Message:   "fix",
				Files:     []string{"main.go"},
				Additions: 2,
				Deletions: 1,
			},
			minScore: 0,
			maxScore: 30,
		},
		{
			name: "documentation only commit",
			commit: models.Commit{
				Message:   "Update README",
				Files:     []string{"README.md"},
				Additions: 10,
				Deletions: 5,
			},
			minScore: 0,
			maxScore: 25,
		},
		{
			name: "config file only commit",
			commit: models.Commit{
				Message:   "Update config",
				Files:     []string{"config.yaml"},
				Additions: 5,
				Deletions: 3,
			},
			minScore: 0,
			maxScore: 25,
		},
		// Medium-value commits
		{
			name: "test file changes",
			commit: models.Commit{
				Message:   "Add unit tests for user service",
				Files:     []string{"internal/user/service_test.go"},
				Additions: 100,
				Deletions: 0,
			},
			minScore: 15,
			maxScore: 60,
		},
		// Large volume changes
		{
			name: "large change volume",
			commit: models.Commit{
				Message:   "Major restructuring of codebase",
				Files:     []string{"pkg/a.go", "pkg/b.go", "pkg/c.go", "pkg/d.go", "pkg/e.go", "pkg/f.go", "pkg/g.go", "pkg/h.go", "pkg/i.go", "pkg/j.go"},
				Additions: 600,
				Deletions: 200,
			},
			minScore: 40,
			maxScore: 100,
		},
		// Edge cases
		{
			name: "empty files list",
			commit: models.Commit{
				Message:   "Some commit message with good description",
				Files:     []string{},
				Additions: 0,
				Deletions: 0,
			},
			minScore: 0,
			maxScore: 30,
		},
		{
			name: "commit with only json files",
			commit: models.Commit{
				Message:   "Update package json dependencies",
				Files:     []string{"package.json"},
				Additions: 20,
				Deletions: 10,
			},
			minScore: 0,
			maxScore: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := scorer.Calculate(tt.commit)

			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("Calculate() = %d, expected between %d and %d for commit: %q",
					score, tt.minScore, tt.maxScore, tt.commit.Message)
			}
		})
	}
}

func TestScorer_Calculate_BoundedScore(t *testing.T) {
	scorer := NewScorer()

	tests := []struct {
		name   string
		commit models.Commit
	}{
		{
			name: "extremely high value commit",
			commit: models.Commit{
				Message:   "feat(auth): implement OAuth2 authentication security feature for API endpoints with database integration",
				Files:     []string{"a.go", "b.go", "c.go", "d.go", "e.go", "f.go", "g.go", "h.go", "i.go", "j.go", "k.go", "l.sql", "m.proto"},
				Additions: 1000,
				Deletions: 500,
			},
		},
		{
			name: "minimum value commit",
			commit: models.Commit{
				Message:   "x",
				Files:     []string{"file.txt"},
				Additions: 1,
				Deletions: 0,
			},
		},
		{
			name: "zero files and changes",
			commit: models.Commit{
				Message:   "",
				Files:     []string{},
				Additions: 0,
				Deletions: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := scorer.Calculate(tt.commit)

			if score < 0 {
				t.Errorf("Calculate() = %d, score should never be negative", score)
			}

			if score > 100 {
				t.Errorf("Calculate() = %d, score should never exceed 100", score)
			}
		})
	}
}

func TestScorer_messageScore(t *testing.T) {
	scorer := NewScorer()

	tests := []struct {
		name     string
		message  string
		minScore int
		maxScore int
	}{
		{
			name:     "conventional commit feature",
			message:  "feat(auth): add login endpoint",
			minScore: 15,
			maxScore: 40,
		},
		{
			name:     "conventional commit fix",
			message:  "fix(api): resolve null pointer exception",
			minScore: 15,
			maxScore: 40,
		},
		{
			name:     "conventional commit perf",
			message:  "perf(db): optimize query execution",
			minScore: 15,
			maxScore: 40,
		},
		{
			name:     "message with implement keyword",
			message:  "Implement caching layer for API responses",
			minScore: 15,
			maxScore: 40,
		},
		{
			name:     "message with optimize keyword",
			message:  "Optimize database connection pooling",
			minScore: 15,
			maxScore: 40,
		},
		{
			name:     "message with security keyword",
			message:  "Add security headers to HTTP responses",
			minScore: 15,
			maxScore: 40,
		},
		{
			name:     "descriptive message 5-20 words",
			message:  "Add new feature for handling user registration with email verification",
			minScore: 10,
			maxScore: 40,
		},
		{
			name:     "very short message",
			message:  "fix bug",
			minScore: 0,
			maxScore: 10,
		},
		{
			name:     "single word message",
			message:  "update",
			minScore: 0,
			maxScore: 5,
		},
		{
			name:     "empty message",
			message:  "",
			minScore: 0,
			maxScore: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := scorer.messageScore(tt.message)

			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("messageScore(%q) = %d, expected between %d and %d",
					tt.message, score, tt.minScore, tt.maxScore)
			}

			// Ensure bounded
			if score < 0 || score > 40 {
				t.Errorf("messageScore(%q) = %d, should be between 0 and 40", tt.message, score)
			}
		})
	}
}

func TestScorer_fileScore(t *testing.T) {
	scorer := NewScorer()

	tests := []struct {
		name     string
		files    []string
		minScore int
		maxScore int
	}{
		{
			name:     "high value Go files",
			files:    []string{"internal/service.go", "pkg/handler.go"},
			minScore: 20,
			maxScore: 30,
		},
		{
			name:     "SQL files",
			files:    []string{"migrations/001.sql"},
			minScore: 12,
			maxScore: 30,
		},
		{
			name:     "proto files",
			files:    []string{"api/service.proto"},
			minScore: 12,
			maxScore: 30,
		},
		{
			name:     "low value text files",
			files:    []string{"notes.txt"},
			minScore: 1,
			maxScore: 5,
		},
		{
			name:     "markdown files",
			files:    []string{"README.md", "CONTRIBUTING.md"},
			minScore: 4,
			maxScore: 10,
		},
		{
			name:     "core directory bonus",
			files:    []string{"src/main.go"},
			minScore: 10,
			maxScore: 20,
		},
		{
			name:     "internal directory bonus",
			files:    []string{"internal/handler/user.go"},
			minScore: 10,
			maxScore: 20,
		},
		{
			name:     "services directory bonus",
			files:    []string{"services/user/service.go"},
			minScore: 10,
			maxScore: 20,
		},
		{
			name:     "empty files list",
			files:    []string{},
			minScore: 0,
			maxScore: 0,
		},
		{
			name:     "unknown file type",
			files:    []string{"unknown.xyz"},
			minScore: 0,
			maxScore: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := scorer.fileScore(tt.files)

			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("fileScore(%v) = %d, expected between %d and %d",
					tt.files, score, tt.minScore, tt.maxScore)
			}

			// Ensure bounded to max 30
			if score > 30 {
				t.Errorf("fileScore(%v) = %d, should not exceed 30", tt.files, score)
			}
		})
	}
}

func TestScorer_volumeScore(t *testing.T) {
	scorer := NewScorer()

	tests := []struct {
		name      string
		additions int
		deletions int
		expected  int
	}{
		{
			name:      "large change >= 500",
			additions: 400,
			deletions: 200,
			expected:  20,
		},
		{
			name:      "medium-large change >= 200",
			additions: 150,
			deletions: 100,
			expected:  15,
		},
		{
			name:      "medium change >= 50",
			additions: 40,
			deletions: 30,
			expected:  10,
		},
		{
			name:      "small change >= 10",
			additions: 8,
			deletions: 5,
			expected:  5,
		},
		{
			name:      "tiny change < 10",
			additions: 3,
			deletions: 2,
			expected:  2,
		},
		{
			name:      "zero changes",
			additions: 0,
			deletions: 0,
			expected:  2,
		},
		{
			name:      "exact boundary 500",
			additions: 500,
			deletions: 0,
			expected:  20,
		},
		{
			name:      "exact boundary 200",
			additions: 200,
			deletions: 0,
			expected:  15,
		},
		{
			name:      "exact boundary 50",
			additions: 50,
			deletions: 0,
			expected:  10,
		},
		{
			name:      "exact boundary 10",
			additions: 10,
			deletions: 0,
			expected:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := scorer.volumeScore(tt.additions, tt.deletions)

			if score != tt.expected {
				t.Errorf("volumeScore(%d, %d) = %d, expected %d",
					tt.additions, tt.deletions, score, tt.expected)
			}
		})
	}
}

func TestScorer_multiFileBonus(t *testing.T) {
	scorer := NewScorer()

	tests := []struct {
		name      string
		fileCount int
		expected  int
	}{
		{
			name:      "10 or more files",
			fileCount: 15,
			expected:  10,
		},
		{
			name:      "exactly 10 files",
			fileCount: 10,
			expected:  10,
		},
		{
			name:      "5-9 files",
			fileCount: 7,
			expected:  7,
		},
		{
			name:      "exactly 5 files",
			fileCount: 5,
			expected:  7,
		},
		{
			name:      "3-4 files",
			fileCount: 4,
			expected:  5,
		},
		{
			name:      "exactly 3 files",
			fileCount: 3,
			expected:  5,
		},
		{
			name:      "2 files",
			fileCount: 2,
			expected:  3,
		},
		{
			name:      "1 file",
			fileCount: 1,
			expected:  0,
		},
		{
			name:      "0 files",
			fileCount: 0,
			expected:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := scorer.multiFileBonus(tt.fileCount)

			if score != tt.expected {
				t.Errorf("multiFileBonus(%d) = %d, expected %d",
					tt.fileCount, score, tt.expected)
			}
		})
	}
}

func TestNewScorer(t *testing.T) {
	scorer := NewScorer()

	if scorer == nil {
		t.Fatal("NewScorer() returned nil")
	}

	if len(scorer.highValuePatterns) == 0 {
		t.Error("highValuePatterns should not be empty")
	}

	if len(scorer.fileTypeWeights) == 0 {
		t.Error("fileTypeWeights should not be empty")
	}

	// Verify some expected file type weights exist
	expectedTypes := []string{".go", ".py", ".java", ".sql", ".proto"}
	for _, ext := range expectedTypes {
		if _, ok := scorer.fileTypeWeights[ext]; !ok {
			t.Errorf("fileTypeWeights should contain %s", ext)
		}
	}
}

func BenchmarkScorer_Calculate(b *testing.B) {
	scorer := NewScorer()
	commit := models.Commit{
		Message:   "feat(auth): implement OAuth2 authentication with Google and GitHub providers",
		Files:     []string{"internal/auth/oauth.go", "internal/auth/providers.go", "internal/auth/handler.go"},
		Additions: 250,
		Deletions: 10,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scorer.Calculate(commit)
	}
}
