package git

import (
	"regexp"
	"strings"

	"github.com/wootaiklee/git-resume/pkg/models"
)

// Filter handles commit and file filtering logic
type Filter struct {
	skipPatterns    []*regexp.Regexp
	skipFileExts    map[string]bool
	skipFilePaths   []string
	skipMessagePats []*regexp.Regexp
}

// NewFilter creates a new filter with default rules
func NewFilter() *Filter {
	return &Filter{
		skipPatterns: compilePatterns([]string{
			`^Merge\s+(branch|pull request)`,
			`^Revert\s+"`,
			`^WIP:?\s*`,
			`^fixup!`,
			`^squash!`,
		}),
		skipFileExts: map[string]bool{
			".lock":    true,
			".sum":     true,
			".mod":     true, // go.mod changes are usually not noteworthy
			".png":     true,
			".jpg":     true,
			".jpeg":    true,
			".gif":     true,
			".svg":     true,
			".ico":     true,
			".woff":    true,
			".woff2":   true,
			".ttf":     true,
			".eot":     true,
			".min.js":  true,
			".min.css": true,
		},
		skipFilePaths: []string{
			"vendor/",
			"node_modules/",
			"dist/",
			"build/",
			".git/",
			"__pycache__/",
			".idea/",
			".vscode/",
		},
		skipMessagePats: compilePatterns([]string{
			`(?i)^fix\s*typo`,
			`(?i)^typo`,
			`(?i)^update\s+changelog`,
			`(?i)^bump\s+version`,
			`(?i)^release\s+v?\d`,
			`(?i)^chore:\s*release`,
			`(?i)^\[skip\s*ci\]`,
			`(?i)^initial\s+commit`,
		}),
	}
}

// ShouldSkip determines if a commit should be skipped
func (f *Filter) ShouldSkip(commit models.Commit) bool {
	msg := strings.TrimSpace(commit.Message)

	// Skip merge commits and reverts
	for _, pat := range f.skipPatterns {
		if pat.MatchString(msg) {
			return true
		}
	}

	// Skip low-value commits based on message
	for _, pat := range f.skipMessagePats {
		if pat.MatchString(msg) {
			return true
		}
	}

	// Skip if no meaningful files changed
	if len(commit.Files) == 0 {
		return true
	}

	return false
}

// ShouldSkipFile determines if a file should be excluded from analysis
func (f *Filter) ShouldSkipFile(filename string) bool {
	lower := strings.ToLower(filename)

	// Check file extensions
	for ext := range f.skipFileExts {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}

	// Check file paths
	for _, path := range f.skipFilePaths {
		if strings.Contains(lower, path) {
			return true
		}
	}

	// Skip specific files
	skipFiles := []string{
		"package-lock.json",
		"yarn.lock",
		"go.sum",
		"poetry.lock",
		"Pipfile.lock",
		"composer.lock",
		".gitignore",
		".dockerignore",
		".editorconfig",
	}

	basename := filename
	if idx := strings.LastIndex(filename, "/"); idx != -1 {
		basename = filename[idx+1:]
	}

	for _, skip := range skipFiles {
		if strings.EqualFold(basename, skip) {
			return true
		}
	}

	return false
}

func compilePatterns(patterns []string) []*regexp.Regexp {
	var compiled []*regexp.Regexp
	for _, p := range patterns {
		if re, err := regexp.Compile(p); err == nil {
			compiled = append(compiled, re)
		}
	}
	return compiled
}
