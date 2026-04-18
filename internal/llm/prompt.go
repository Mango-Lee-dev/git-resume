package llm

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/wootaiklee/git-resume/pkg/models"
)

// SystemPrompt defines the AI persona for resume generation
const SystemPrompt = `You are an expert technical writer specializing in resume optimization for software engineers. Your task is to transform git commit information into impactful resume bullet points.

Guidelines:
1. Use the STAR method (Situation, Task, Action, Result) implicitly
2. Focus on business impact and technical achievements
3. Use action verbs (Implemented, Developed, Optimized, etc.)
4. Quantify impact when possible (performance improvements, code reduction, etc.)
5. Keep each bullet point concise (1-2 sentences max)
6. Categorize each change appropriately

Output format: Return a JSON array with objects containing:
- "hash": commit hash (first 7 chars)
- "category": one of "Feature", "Fix", "Refactor", "Docs", "Test", "Chore"
- "summary": the resume bullet point (in English)

Example output:
[
  {"hash": "abc1234", "category": "Feature", "summary": "Implemented real-time notification system using WebSocket, reducing user response time by 40%"},
  {"hash": "def5678", "category": "Fix", "summary": "Resolved critical authentication vulnerability affecting 10K+ users by implementing secure token validation"}
]`

// BuildPrompt constructs the analysis prompt for a batch of commits
func BuildPrompt(batch models.CommitBatch) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Analyze the following git commits from project '%s' and generate resume-worthy bullet points.\n\n", batch.Project))
	sb.WriteString("Commits:\n")

	for i, commit := range batch.Commits {
		sb.WriteString(fmt.Sprintf("\n--- Commit %d ---\n", i+1))
		sb.WriteString(fmt.Sprintf("Hash: %s\n", commit.Hash[:7]))
		sb.WriteString(fmt.Sprintf("Date: %s\n", commit.Date.Format("2006-01-02")))
		sb.WriteString(fmt.Sprintf("Message: %s\n", strings.TrimSpace(commit.Message)))
		sb.WriteString(fmt.Sprintf("Files changed: %d (+%d/-%d lines)\n",
			len(commit.Files), commit.Additions, commit.Deletions))

		if len(commit.Files) > 0 {
			sb.WriteString("Key files:\n")
			// Limit to 10 files to save tokens
			maxFiles := 10
			if len(commit.Files) < maxFiles {
				maxFiles = len(commit.Files)
			}
			for _, file := range commit.Files[:maxFiles] {
				sb.WriteString(fmt.Sprintf("  - %s\n", file))
			}
			if len(commit.Files) > 10 {
				sb.WriteString(fmt.Sprintf("  ... and %d more files\n", len(commit.Files)-10))
			}
		}
	}

	sb.WriteString("\nGenerate resume bullet points in JSON format as specified.")

	return sb.String()
}

// llmOutput represents the expected JSON output from Claude
type llmOutput struct {
	Hash     string `json:"hash"`
	Category string `json:"category"`
	Summary  string `json:"summary"`
}

// ParseResponse parses the Claude API response into analysis results
func ParseResponse(responseText string, batch models.CommitBatch) ([]models.AnalysisResult, error) {
	// Find JSON array in response
	start := strings.Index(responseText, "[")
	end := strings.LastIndex(responseText, "]")

	if start == -1 || end == -1 || end <= start {
		return nil, fmt.Errorf("no valid JSON array found in response")
	}

	jsonStr := responseText[start : end+1]

	var outputs []llmOutput
	if err := json.Unmarshal([]byte(jsonStr), &outputs); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Map outputs back to commits
	hashToCommit := make(map[string]models.Commit)
	for _, commit := range batch.Commits {
		hashToCommit[commit.Hash[:7]] = commit
	}

	var results []models.AnalysisResult
	for _, out := range outputs {
		commit, ok := hashToCommit[out.Hash]
		if !ok {
			continue // Skip if hash doesn't match
		}

		results = append(results, models.AnalysisResult{
			CommitHash:    commit.Hash,
			Date:          commit.Date,
			Project:       batch.Project,
			Category:      parseCategory(out.Category),
			ImpactSummary: out.Summary,
			CreatedAt:     time.Now(),
		})
	}

	return results, nil
}

func parseCategory(cat string) models.Category {
	switch strings.ToLower(cat) {
	case "feature":
		return models.CategoryFeature
	case "fix":
		return models.CategoryFix
	case "refactor":
		return models.CategoryRefactor
	case "docs":
		return models.CategoryDocs
	case "test":
		return models.CategoryTest
	default:
		return models.CategoryChore
	}
}
