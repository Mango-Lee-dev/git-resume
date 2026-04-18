package llm

import (
	"errors"
	"math"
	"time"
)

// RateLimitError indicates an API rate limit was hit
type RateLimitError struct {
	RetryAfter time.Duration
}

func (e *RateLimitError) Error() string {
	return "rate limit exceeded"
}

// Retrier handles retry logic with exponential backoff
type Retrier struct {
	maxAttempts int
	baseDelay   time.Duration
	maxDelay    time.Duration
}

// NewRetrier creates a new retrier with default settings
func NewRetrier() *Retrier {
	return &Retrier{
		maxAttempts: 5,
		baseDelay:   1 * time.Second,
		maxDelay:    60 * time.Second,
	}
}

// Do executes the function with retry logic
func (r *Retrier) Do(fn func() error) error {
	var lastErr error

	for attempt := 0; attempt < r.maxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if it's a rate limit error
		var rateLimitErr *RateLimitError
		if errors.As(err, &rateLimitErr) {
			time.Sleep(rateLimitErr.RetryAfter)
			continue
		}

		// Check if error is retryable
		if !isRetryable(err) {
			return err
		}

		// Calculate delay with exponential backoff
		delay := r.calculateDelay(attempt)
		time.Sleep(delay)
	}

	return lastErr
}

func (r *Retrier) calculateDelay(attempt int) time.Duration {
	// Exponential backoff: baseDelay * 2^attempt
	delay := float64(r.baseDelay) * math.Pow(2, float64(attempt))

	// Add jitter (±20%)
	jitter := delay * 0.2
	delay = delay - jitter + (jitter * 2 * float64(time.Now().UnixNano()%100) / 100)

	if time.Duration(delay) > r.maxDelay {
		return r.maxDelay
	}

	return time.Duration(delay)
}

// isRetryable determines if an error should trigger a retry
func isRetryable(err error) bool {
	// Rate limit errors are handled separately
	var rateLimitErr *RateLimitError
	if errors.As(err, &rateLimitErr) {
		return true
	}

	// Network errors are generally retryable
	errMsg := err.Error()
	retryablePatterns := []string{
		"connection refused",
		"connection reset",
		"timeout",
		"temporary failure",
		"server error",
		"status 500",
		"status 502",
		"status 503",
		"status 504",
	}

	for _, pattern := range retryablePatterns {
		if contains(errMsg, pattern) {
			return true
		}
	}

	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsImpl(s, substr))
}

func containsImpl(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
