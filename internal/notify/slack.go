package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/wootaiklee/git-resume/pkg/models"
)

// SlackNotifier sends notifications to Slack
type SlackNotifier struct {
	webhookURL string
	client     *http.Client
}

// SlackMessage represents a Slack webhook message
type SlackMessage struct {
	Text        string       `json:"text,omitempty"`
	Blocks      []SlackBlock `json:"blocks,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// SlackBlock represents a Slack block element
type SlackBlock struct {
	Type   string      `json:"type"`
	Text   *SlackText  `json:"text,omitempty"`
	Fields []SlackText `json:"fields,omitempty"`
}

// SlackText represents text content in Slack
type SlackText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Attachment represents a Slack attachment
type Attachment struct {
	Color  string `json:"color,omitempty"`
	Title  string `json:"title,omitempty"`
	Text   string `json:"text,omitempty"`
	Footer string `json:"footer,omitempty"`
}

// NewSlackNotifier creates a new Slack notifier
func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendAnalysisComplete sends a notification when analysis is complete
func (s *SlackNotifier) SendAnalysisComplete(results []models.AnalysisResult, fromDate, toDate time.Time) error {
	if s.webhookURL == "" {
		return fmt.Errorf("slack webhook URL not configured")
	}

	// Count by category
	categoryCounts := make(map[models.Category]int)
	for _, r := range results {
		categoryCounts[r.Category]++
	}

	// Build category summary
	categoryText := ""
	for cat, count := range categoryCounts {
		emoji := getCategoryEmoji(cat)
		categoryText += fmt.Sprintf("%s %s: %d\n", emoji, cat, count)
	}

	// Build message
	msg := SlackMessage{
		Blocks: []SlackBlock{
			{
				Type: "header",
				Text: &SlackText{
					Type: "plain_text",
					Text: "📊 Resume Generation Complete",
				},
			},
			{
				Type: "section",
				Text: &SlackText{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Period:* %s to %s\n*Total bullet points:* %d",
						fromDate.Format("2006-01-02"),
						toDate.Format("2006-01-02"),
						len(results)),
				},
			},
			{
				Type: "section",
				Text: &SlackText{
					Type: "mrkdwn",
					Text: "*Breakdown by Category:*\n" + categoryText,
				},
			},
			{
				Type: "divider",
			},
		},
	}

	// Add top achievements (up to 5)
	if len(results) > 0 {
		achievementsText := "*Top Achievements:*\n"
		limit := 5
		if len(results) < limit {
			limit = len(results)
		}

		for i := 0; i < limit; i++ {
			r := results[i]
			achievementsText += fmt.Sprintf("• %s\n", truncateText(r.ImpactSummary, 100))
		}

		msg.Blocks = append(msg.Blocks, SlackBlock{
			Type: "section",
			Text: &SlackText{
				Type: "mrkdwn",
				Text: achievementsText,
			},
		})
	}

	return s.send(msg)
}

// SendError sends an error notification
func (s *SlackNotifier) SendError(err error, context string) error {
	if s.webhookURL == "" {
		return fmt.Errorf("slack webhook URL not configured")
	}

	msg := SlackMessage{
		Blocks: []SlackBlock{
			{
				Type: "header",
				Text: &SlackText{
					Type: "plain_text",
					Text: "⚠️ Resume Generation Error",
				},
			},
			{
				Type: "section",
				Text: &SlackText{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Context:* %s\n*Error:* %s", context, err.Error()),
				},
			},
		},
		Attachments: []Attachment{
			{
				Color:  "#FF0000",
				Footer: fmt.Sprintf("Error occurred at %s", time.Now().Format(time.RFC3339)),
			},
		},
	}

	return s.send(msg)
}

// SendDailyDigest sends a daily summary of work
func (s *SlackNotifier) SendDailyDigest(results []models.AnalysisResult) error {
	if s.webhookURL == "" {
		return fmt.Errorf("slack webhook URL not configured")
	}

	if len(results) == 0 {
		return nil // Skip if no results
	}

	achievementsText := ""
	for _, r := range results {
		achievementsText += fmt.Sprintf("• %s `%s`\n", r.ImpactSummary, r.CommitHash[:7])
	}

	msg := SlackMessage{
		Blocks: []SlackBlock{
			{
				Type: "header",
				Text: &SlackText{
					Type: "plain_text",
					Text: fmt.Sprintf("📝 Daily Digest - %s", time.Now().Format("Jan 2, 2006")),
				},
			},
			{
				Type: "section",
				Text: &SlackText{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*%d achievements generated today:*\n%s", len(results), achievementsText),
				},
			},
		},
	}

	return s.send(msg)
}

func (s *SlackNotifier) send(msg SlackMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequest("POST", s.webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack returned status %d", resp.StatusCode)
	}

	return nil
}

func getCategoryEmoji(cat models.Category) string {
	switch cat {
	case models.CategoryFeature:
		return "✨"
	case models.CategoryFix:
		return "🐛"
	case models.CategoryRefactor:
		return "♻️"
	case models.CategoryTest:
		return "🧪"
	case models.CategoryDocs:
		return "📚"
	case models.CategoryChore:
		return "🔧"
	default:
		return "📌"
	}
}

func truncateText(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
