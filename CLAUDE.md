# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Git Resume Analyzer transforms Git commit history into STAR-format resume bullet points using Claude AI. It parses commits from local repositories, filters noise, scores importance, and generates achievement-focused summaries suitable for resumes.

## Build and Test Commands

```bash
# Build binary to bin/git-resume
make build

# Run tests
make test
go test -v ./...                    # all tests
go test -v ./internal/git/...       # specific package

# Lint (requires golangci-lint)
make lint

# Run without building
make run ARGS="analyze --dry-run"
go run . analyze --dry-run

# Clean build artifacts and cache databases
make clean
```

## Architecture

### Data Flow
1. **Git Parsing** (`internal/git/parser.go`) - Uses go-git to extract commits within date range
2. **Filtering** (`internal/git/filter.go`) - Removes noise commits (merges, typos, lock files)
3. **Scoring** (`internal/git/scorer.go`) - Calculates importance score (0-100) based on file count, line changes, file types
4. **Caching** (`internal/db/cache.go`) - SQLite tracks processed commits to avoid duplicate API calls
5. **LLM Analysis** (`internal/llm/client.go`) - Batches commits to Claude API with retry logic
6. **Export** (`internal/export/`) - Outputs CSV, Markdown, or JSON

### Key Patterns

**CLI Structure**: Uses Cobra for commands (`cmd/`). Main commands:
- `analyze` - Core analysis workflow
- `export` - Export cached results
- `estimate` - Token cost estimation
- `templates` - List available templates
- `tui` - Interactive terminal UI (Bubble Tea)

**Terminal UI** (`internal/tui/`): Bubble Tea-based interactive interface with screen-based navigation:
- `app.go` - Main model with screen routing
- `screens/` - Individual screens (welcome, repo_select, date_range, template, options, analysis, results, export, done)
- `state/state.go` - Shared application state
- `styles/styles.go` - Lip Gloss styling

**Template System** (`internal/llm/templates.go`): Seven built-in templates (default, startup, enterprise, backend, frontend, devops, data) customize LLM prompts with different personas, tones, and action verbs.

**CGO-free SQLite**: Uses `modernc.org/sqlite` for cross-platform compatibility without CGO.

**Retry with Backoff** (`internal/llm/retry.go`): Handles rate limits (429) with exponential backoff.

### Models (`pkg/models/`)
- `Commit` - Parsed git commit with files, additions/deletions, score
- `CommitBatch` - Group of commits for batch API processing
- `AnalysisResult` - LLM output with category and impact summary
- `ExportFormat` - CSV/JSON/Markdown enum

## Environment Configuration

Required in `.env`:
```
CLAUDE_API_KEY=your_api_key_here
```

Optional:
```
DEFAULT_REPO_PATH=/path/to/repo
DB_PATH=./data/cache.db
OUTPUT_DIR=./output
SLACK_WEBHOOK_URL=https://hooks.slack.com/...
```

## Common CLI Patterns

```bash
# Launch interactive TUI
./bin/git-resume tui

# Analyze current month
./bin/git-resume analyze

# Analyze specific month/date range
./bin/git-resume analyze --month=4 --year=2024
./bin/git-resume analyze --from=2024-01-01 --to=2024-03-31

# Multiple repositories
./bin/git-resume analyze --repos=/path/repo1,/path/repo2

# Use specific template
./bin/git-resume analyze --template=backend

# Dry run (no API calls)
./bin/git-resume analyze --dry-run
```
