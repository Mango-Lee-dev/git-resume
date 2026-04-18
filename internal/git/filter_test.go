package git

import (
	"testing"

	"github.com/wootaiklee/git-resume/pkg/models"
)

func TestFilter_ShouldSkip(t *testing.T) {
	filter := NewFilter()

	tests := []struct {
		name     string
		commit   models.Commit
		expected bool
	}{
		// Merge commits should be skipped
		{
			name: "merge branch commit",
			commit: models.Commit{
				Message: "Merge branch 'feature/auth' into main",
				Files:   []string{"auth.go"},
			},
			expected: true,
		},
		{
			name: "merge pull request commit",
			commit: models.Commit{
				Message: "Merge pull request #123 from user/feature",
				Files:   []string{"feature.go"},
			},
			expected: true,
		},
		// Revert commits should be skipped
		{
			name: "revert commit",
			commit: models.Commit{
				Message: "Revert \"Add new feature\"",
				Files:   []string{"feature.go"},
			},
			expected: true,
		},
		// WIP commits should be skipped
		{
			name: "WIP commit with colon",
			commit: models.Commit{
				Message: "WIP: working on new feature",
				Files:   []string{"feature.go"},
			},
			expected: true,
		},
		{
			name: "WIP commit without colon",
			commit: models.Commit{
				Message: "WIP work in progress",
				Files:   []string{"feature.go"},
			},
			expected: true,
		},
		// Fixup and squash commits
		{
			name: "fixup commit",
			commit: models.Commit{
				Message: "fixup! original commit message",
				Files:   []string{"file.go"},
			},
			expected: true,
		},
		{
			name: "squash commit",
			commit: models.Commit{
				Message: "squash! original commit message",
				Files:   []string{"file.go"},
			},
			expected: true,
		},
		// Typo fixes should be skipped
		{
			name: "fix typo commit",
			commit: models.Commit{
				Message: "Fix typo in documentation",
				Files:   []string{"README.md"},
			},
			expected: true,
		},
		{
			name: "typo only commit",
			commit: models.Commit{
				Message: "Typo fix",
				Files:   []string{"file.go"},
			},
			expected: true,
		},
		// Changelog and version bumps
		{
			name: "update changelog commit",
			commit: models.Commit{
				Message: "Update changelog for release",
				Files:   []string{"CHANGELOG.md"},
			},
			expected: true,
		},
		{
			name: "bump version commit",
			commit: models.Commit{
				Message: "Bump version to 1.2.0",
				Files:   []string{"version.go"},
			},
			expected: true,
		},
		// Release commits
		{
			name: "release commit with v prefix",
			commit: models.Commit{
				Message: "Release v1.0.0",
				Files:   []string{"version.go"},
			},
			expected: true,
		},
		{
			name: "chore release commit",
			commit: models.Commit{
				Message: "chore: release 2.0.0",
				Files:   []string{"package.json"},
			},
			expected: true,
		},
		// Skip CI commits
		{
			name: "skip ci commit",
			commit: models.Commit{
				Message: "[skip ci] Update docs",
				Files:   []string{"docs.md"},
			},
			expected: true,
		},
		{
			name: "skip ci with space",
			commit: models.Commit{
				Message: "[skip CI] Minor update",
				Files:   []string{"file.go"},
			},
			expected: true,
		},
		// Initial commit
		{
			name: "initial commit",
			commit: models.Commit{
				Message: "Initial commit",
				Files:   []string{"README.md"},
			},
			expected: true,
		},
		// Empty files should be skipped
		{
			name: "commit with no files",
			commit: models.Commit{
				Message: "Some valid message",
				Files:   []string{},
			},
			expected: true,
		},
		{
			name: "commit with nil files",
			commit: models.Commit{
				Message: "Some valid message",
				Files:   nil,
			},
			expected: true,
		},
		// Normal commits should NOT be skipped
		{
			name: "normal feature commit",
			commit: models.Commit{
				Message: "Add user authentication module",
				Files:   []string{"auth/handler.go", "auth/service.go"},
			},
			expected: false,
		},
		{
			name: "conventional commit feature",
			commit: models.Commit{
				Message: "feat(auth): implement OAuth2 login",
				Files:   []string{"auth/oauth.go"},
			},
			expected: false,
		},
		{
			name: "bug fix commit",
			commit: models.Commit{
				Message: "fix: resolve race condition in worker pool",
				Files:   []string{"worker/pool.go"},
			},
			expected: false,
		},
		{
			name: "refactor commit",
			commit: models.Commit{
				Message: "Refactor database connection handling",
				Files:   []string{"db/connection.go"},
			},
			expected: false,
		},
		// Edge cases
		{
			name: "commit with leading whitespace",
			commit: models.Commit{
				Message: "   Merge branch 'feature'",
				Files:   []string{"file.go"},
			},
			expected: true,
		},
		{
			name: "commit with trailing whitespace",
			commit: models.Commit{
				Message: "Add new feature   ",
				Files:   []string{"feature.go"},
			},
			expected: false,
		},
		{
			name: "commit containing merge in middle",
			commit: models.Commit{
				Message: "Fix merge conflict resolution logic",
				Files:   []string{"merge.go"},
			},
			expected: false, // "merge" in middle of message is fine
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filter.ShouldSkip(tt.commit)
			if got != tt.expected {
				t.Errorf("ShouldSkip() = %v, expected %v for message: %q", got, tt.expected, tt.commit.Message)
			}
		})
	}
}

func TestFilter_ShouldSkipFile(t *testing.T) {
	filter := NewFilter()

	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		// Lock files should be skipped
		{
			name:     "package-lock.json",
			filename: "package-lock.json",
			expected: true,
		},
		{
			name:     "yarn.lock",
			filename: "yarn.lock",
			expected: true,
		},
		{
			name:     "go.sum",
			filename: "go.sum",
			expected: true,
		},
		{
			name:     "poetry.lock",
			filename: "poetry.lock",
			expected: true,
		},
		{
			name:     "Pipfile.lock",
			filename: "Pipfile.lock",
			expected: true,
		},
		{
			name:     "composer.lock",
			filename: "composer.lock",
			expected: true,
		},
		// File extensions to skip
		{
			name:     ".lock extension",
			filename: "some-package.lock",
			expected: true,
		},
		{
			name:     ".sum extension",
			filename: "checksum.sum",
			expected: true,
		},
		{
			name:     ".mod extension",
			filename: "go.mod",
			expected: true,
		},
		// Image files
		{
			name:     "PNG image",
			filename: "logo.png",
			expected: true,
		},
		{
			name:     "JPG image",
			filename: "photo.jpg",
			expected: true,
		},
		{
			name:     "JPEG image",
			filename: "image.jpeg",
			expected: true,
		},
		{
			name:     "GIF image",
			filename: "animation.gif",
			expected: true,
		},
		{
			name:     "SVG image",
			filename: "icon.svg",
			expected: true,
		},
		{
			name:     "ICO favicon",
			filename: "favicon.ico",
			expected: true,
		},
		// Font files
		{
			name:     "WOFF font",
			filename: "font.woff",
			expected: true,
		},
		{
			name:     "WOFF2 font",
			filename: "font.woff2",
			expected: true,
		},
		{
			name:     "TTF font",
			filename: "font.ttf",
			expected: true,
		},
		{
			name:     "EOT font",
			filename: "font.eot",
			expected: true,
		},
		// Minified files
		{
			name:     "minified JS",
			filename: "app.min.js",
			expected: true,
		},
		{
			name:     "minified CSS",
			filename: "styles.min.css",
			expected: true,
		},
		// Vendor directories
		{
			name:     "vendor Go file",
			filename: "vendor/github.com/pkg/errors/errors.go",
			expected: true,
		},
		{
			name:     "node_modules file",
			filename: "node_modules/lodash/index.js",
			expected: true,
		},
		{
			name:     "dist directory",
			filename: "dist/bundle.js",
			expected: true,
		},
		{
			name:     "build directory",
			filename: "build/output.js",
			expected: true,
		},
		{
			name:     ".git directory",
			filename: ".git/config",
			expected: true,
		},
		{
			name:     "__pycache__ directory",
			filename: "__pycache__/module.cpython-39.pyc",
			expected: true,
		},
		{
			name:     ".idea directory",
			filename: ".idea/workspace.xml",
			expected: true,
		},
		{
			name:     ".vscode directory",
			filename: ".vscode/settings.json",
			expected: true,
		},
		// Config files to skip
		{
			name:     ".gitignore",
			filename: ".gitignore",
			expected: true,
		},
		{
			name:     ".dockerignore",
			filename: ".dockerignore",
			expected: true,
		},
		{
			name:     ".editorconfig",
			filename: ".editorconfig",
			expected: true,
		},
		// Normal source files should NOT be skipped
		{
			name:     "Go source file",
			filename: "main.go",
			expected: false,
		},
		{
			name:     "Python file",
			filename: "app.py",
			expected: false,
		},
		{
			name:     "JavaScript file",
			filename: "index.js",
			expected: false,
		},
		{
			name:     "TypeScript file",
			filename: "component.tsx",
			expected: false,
		},
		{
			name:     "Java file",
			filename: "Main.java",
			expected: false,
		},
		{
			name:     "Rust file",
			filename: "lib.rs",
			expected: false,
		},
		{
			name:     "SQL file",
			filename: "migrations/001_init.sql",
			expected: false,
		},
		{
			name:     "Markdown file",
			filename: "README.md",
			expected: false,
		},
		{
			name:     "YAML config",
			filename: "config.yaml",
			expected: false,
		},
		{
			name:     "JSON config",
			filename: "config.json",
			expected: false,
		},
		// Deep nested source files
		{
			name:     "nested Go file in internal",
			filename: "internal/git/parser.go",
			expected: false,
		},
		{
			name:     "nested file in src",
			filename: "src/components/Button.tsx",
			expected: false,
		},
		// Edge cases
		{
			name:     "case insensitive PNG",
			filename: "IMAGE.PNG",
			expected: true,
		},
		{
			name:     "case insensitive package-lock",
			filename: "PACKAGE-LOCK.JSON",
			expected: true,
		},
		{
			name:     "file with path containing vendor-like name",
			filename: "src/vendor-utils/helper.go",
			expected: false, // Does NOT match "vendor/" pattern, vendor-utils is different
		},
		{
			name:     "actual vendor directory file",
			filename: "src/vendor/github.com/pkg/errors.go",
			expected: true, // Contains "vendor/" pattern
		},
		{
			name:     "non-minified JS",
			filename: "app.js",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filter.ShouldSkipFile(tt.filename)
			if got != tt.expected {
				t.Errorf("ShouldSkipFile(%q) = %v, expected %v", tt.filename, got, tt.expected)
			}
		})
	}
}

func TestNewFilter(t *testing.T) {
	filter := NewFilter()

	if filter == nil {
		t.Fatal("NewFilter() returned nil")
	}

	if len(filter.skipPatterns) == 0 {
		t.Error("skipPatterns should not be empty")
	}

	if len(filter.skipFileExts) == 0 {
		t.Error("skipFileExts should not be empty")
	}

	if len(filter.skipFilePaths) == 0 {
		t.Error("skipFilePaths should not be empty")
	}

	if len(filter.skipMessagePats) == 0 {
		t.Error("skipMessagePats should not be empty")
	}
}

func TestFilter_ShouldSkip_CombinedConditions(t *testing.T) {
	filter := NewFilter()

	// A commit that matches multiple skip patterns should still return true
	commit := models.Commit{
		Message: "Merge branch 'fix-typo'", // Matches merge pattern
		Files:   []string{"README.md"},
	}

	if !filter.ShouldSkip(commit) {
		t.Error("ShouldSkip should return true for commit matching multiple patterns")
	}
}
