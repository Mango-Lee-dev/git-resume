package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/wootaiklee/git-resume/pkg/models"
)

const (
	claudeAPIURL   = "https://api.anthropic.com/v1/messages"
	claudeModel    = "claude-sonnet-4-20250514"
	maxTokens      = 1024
	defaultTimeout = 60 * time.Second
)

// Client handles Claude API interactions
type Client struct {
	apiKey      string
	httpClient  *http.Client
	retrier     *Retrier
	templateMgr *TemplateManager
}

// NewClient creates a new Claude API client with default template
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		retrier:     NewRetrier(),
		templateMgr: nil, // Use default SystemPrompt
	}
}

// NewClientWithTemplate creates a new Claude API client with custom template
func NewClientWithTemplate(apiKey string, templateMgr *TemplateManager) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		retrier:     NewRetrier(),
		templateMgr: templateMgr,
	}
}

// claudeRequest represents the API request structure
type claudeRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []message `json:"messages"`
	System    string    `json:"system,omitempty"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// claudeResponse represents the API response structure
type claudeResponse struct {
	ID      string `json:"id"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// AnalyzeCommits analyzes a batch of commits and returns resume bullet points
func (c *Client) AnalyzeCommits(batch models.CommitBatch) (*models.BatchResult, error) {
	prompt := BuildPrompt(batch)

	var response *claudeResponse
	var err error

	// Use retrier for resilient API calls
	err = c.retrier.Do(func() error {
		response, err = c.sendRequest(prompt)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to analyze commits: %w", err)
	}

	// Parse the response
	results, err := ParseResponse(response.Content[0].Text, batch)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &models.BatchResult{
		Results:     results,
		TokensUsed:  response.Usage.InputTokens + response.Usage.OutputTokens,
		ProcessedAt: time.Now(),
	}, nil
}

func (c *Client) sendRequest(prompt string) (*claudeResponse, error) {
	// Use template's system prompt if available, otherwise use default
	systemPrompt := SystemPrompt
	if c.templateMgr != nil {
		systemPrompt = c.templateMgr.BuildSystemPrompt()
	}

	reqBody := claudeRequest{
		Model:     claudeModel,
		MaxTokens: maxTokens,
		System:    systemPrompt,
		Messages: []message{
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", claudeAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 429 {
		return nil, &RateLimitError{RetryAfter: parseRetryAfter(resp)}
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return nil, err
	}

	if claudeResp.Error != nil {
		return nil, fmt.Errorf("API error: %s", claudeResp.Error.Message)
	}

	return &claudeResp, nil
}

func parseRetryAfter(resp *http.Response) time.Duration {
	// Default retry after
	return 30 * time.Second
}

// EstimateTokens estimates token count for a batch (rough approximation)
func EstimateTokens(batch models.CommitBatch) int {
	// Rough estimate: ~4 chars per token
	totalChars := len(SystemPrompt)
	for _, commit := range batch.Commits {
		totalChars += len(commit.Message) + len(commit.Hash)
		for _, file := range commit.Files {
			totalChars += len(file)
		}
	}
	return totalChars / 4
}
