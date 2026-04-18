package git

import (
	"regexp"
	"strings"

	"github.com/wootaiklee/git-resume/pkg/models"
)

// Scorer calculates importance scores for commits
type Scorer struct {
	highValuePatterns []*regexp.Regexp
	fileTypeWeights   map[string]int
}

// NewScorer creates a new scorer with default rules
func NewScorer() *Scorer {
	return &Scorer{
		highValuePatterns: compilePatterns([]string{
			`(?i)^feat(\(.+\))?:`,    // Conventional commit: feature
			`(?i)^fix(\(.+\))?:`,     // Conventional commit: fix
			`(?i)^perf(\(.+\))?:`,    // Conventional commit: performance
			`(?i)implement`,          // Implementation work
			`(?i)add\s+.+\s+feature`, // Feature addition
			`(?i)improve`,            // Improvements
			`(?i)optimize`,           // Optimizations
			`(?i)refactor`,           // Refactoring
			`(?i)migrate`,            // Migrations
			`(?i)integrate`,          // Integrations
			`(?i)security`,           // Security fixes
			`(?i)api`,                // API changes
			`(?i)database`,           // Database changes
			`(?i)authentication`,     // Auth changes
			`(?i)authorization`,      // Auth changes
		}),
		fileTypeWeights: map[string]int{
			".go":    10,
			".py":    10,
			".java":  10,
			".ts":    10,
			".tsx":   10,
			".js":    8,
			".jsx":   8,
			".rs":    10,
			".rb":    10,
			".sql":   12, // Database changes are often impactful
			".proto": 12, // API definitions
			".yaml":  5,
			".yml":   5,
			".json":  3,
			".md":    2,
			".txt":   1,
		},
	}
}

// Calculate computes an importance score (0-100) for a commit
func (s *Scorer) Calculate(commit models.Commit) int {
	score := 0

	// Base score from message patterns (max 40)
	score += s.messageScore(commit.Message)

	// Score from file changes (max 30)
	score += s.fileScore(commit.Files)

	// Score from change volume (max 20)
	score += s.volumeScore(commit.Additions, commit.Deletions)

	// Bonus for multiple files (max 10)
	score += s.multiFileBonus(len(commit.Files))

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return score
}

func (s *Scorer) messageScore(message string) int {
	score := 0
	msg := strings.ToLower(message)

	// Check high-value patterns
	for _, pat := range s.highValuePatterns {
		if pat.MatchString(message) {
			score += 15
			break // Only count once
		}
	}

	// Bonus for descriptive messages
	wordCount := len(strings.Fields(msg))
	if wordCount >= 5 && wordCount <= 20 {
		score += 10
	}

	// Bonus for conventional commits
	if strings.Contains(message, ":") && len(message) > 10 {
		score += 5
	}

	// Penalty for very short messages
	if wordCount < 3 {
		score -= 10
	}

	if score < 0 {
		score = 0
	}
	if score > 40 {
		score = 40
	}

	return score
}

func (s *Scorer) fileScore(files []string) int {
	score := 0

	for _, file := range files {
		lower := strings.ToLower(file)

		// Check file type weights
		for ext, weight := range s.fileTypeWeights {
			if strings.HasSuffix(lower, ext) {
				score += weight
				break
			}
		}

		// Bonus for core directories
		corePatterns := []string{
			"/src/", "/lib/", "/pkg/", "/internal/",
			"/core/", "/services/", "/handlers/",
			"/api/", "/models/", "/controllers/",
		}
		for _, pat := range corePatterns {
			if strings.Contains(lower, pat) {
				score += 3
				break
			}
		}
	}

	if score > 30 {
		score = 30
	}

	return score
}

func (s *Scorer) volumeScore(additions, deletions int) int {
	total := additions + deletions

	switch {
	case total >= 500:
		return 20 // Large change
	case total >= 200:
		return 15
	case total >= 50:
		return 10
	case total >= 10:
		return 5
	default:
		return 2
	}
}

func (s *Scorer) multiFileBonus(fileCount int) int {
	switch {
	case fileCount >= 10:
		return 10
	case fileCount >= 5:
		return 7
	case fileCount >= 3:
		return 5
	case fileCount >= 2:
		return 3
	default:
		return 0
	}
}
