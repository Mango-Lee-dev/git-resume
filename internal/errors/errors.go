// Package errors provides custom error types for the git-resume project.
// Each error type includes relevant context and supports error wrapping
// for proper error chain handling with errors.Is and errors.As.
package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Sentinel errors for common scenarios.
// Use these with errors.Is() for error checking.
var (
	ErrRepositoryNotFound = errors.New("repository not found")
	ErrInvalidPath        = errors.New("invalid repository path")
	ErrRateLimited        = errors.New("rate limit exceeded")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrNetworkFailure     = errors.New("network failure")
	ErrMissingAPIKey      = errors.New("missing API key")
	ErrInvalidConfig      = errors.New("invalid configuration")
	ErrFileWriteFailed    = errors.New("file write failed")
)

// GitError represents errors related to git repository operations.
type GitError struct {
	Op   string // Operation that failed (e.g., "open", "clone", "log")
	Path string // Repository path
	Err  error  // Underlying error
}

// Error implements the error interface.
func (e *GitError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("git %s: %s: %v", e.Op, e.Path, e.Err)
	}
	return fmt.Sprintf("git %s: %v", e.Op, e.Err)
}

// Unwrap returns the underlying error for error chain support.
func (e *GitError) Unwrap() error {
	return e.Err
}

// NewGitError creates a new GitError with the given operation, path, and underlying error.
func NewGitError(op, path string, err error) *GitError {
	return &GitError{
		Op:   op,
		Path: path,
		Err:  err,
	}
}

// NewRepositoryNotFoundError creates a GitError for a missing repository.
func NewRepositoryNotFoundError(path string) *GitError {
	return &GitError{
		Op:   "open",
		Path: path,
		Err:  ErrRepositoryNotFound,
	}
}

// NewInvalidPathError creates a GitError for an invalid repository path.
func NewInvalidPathError(path string) *GitError {
	return &GitError{
		Op:   "validate",
		Path: path,
		Err:  ErrInvalidPath,
	}
}

// APIError represents errors from external API calls (e.g., LLM providers).
type APIError struct {
	Op         string // Operation that failed (e.g., "generate", "summarize")
	StatusCode int    // HTTP status code (0 if not applicable)
	Message    string // Error message from API
	Err        error  // Underlying error
	Retryable  bool   // Whether this error is retryable
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("api %s: status %d: %s", e.Op, e.StatusCode, e.Message)
	}
	if e.Message != "" {
		return fmt.Sprintf("api %s: %s", e.Op, e.Message)
	}
	return fmt.Sprintf("api %s: %v", e.Op, e.Err)
}

// Unwrap returns the underlying error for error chain support.
func (e *APIError) Unwrap() error {
	return e.Err
}

// NewAPIError creates a new APIError with full context.
func NewAPIError(op string, statusCode int, message string, err error) *APIError {
	retryable := isRetryableStatusCode(statusCode)
	return &APIError{
		Op:         op,
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
		Retryable:  retryable,
	}
}

// NewRateLimitError creates an APIError for rate limiting.
func NewRateLimitError(op string) *APIError {
	return &APIError{
		Op:         op,
		StatusCode: http.StatusTooManyRequests,
		Message:    "rate limit exceeded, please retry later",
		Err:        ErrRateLimited,
		Retryable:  true,
	}
}

// NewAuthError creates an APIError for authentication failures.
func NewAuthError(op, message string) *APIError {
	return &APIError{
		Op:         op,
		StatusCode: http.StatusUnauthorized,
		Message:    message,
		Err:        ErrUnauthorized,
		Retryable:  false,
	}
}

// NewNetworkError creates an APIError for network-related failures.
func NewNetworkError(op string, err error) *APIError {
	return &APIError{
		Op:        op,
		Message:   "network error occurred",
		Err:       fmt.Errorf("%w: %v", ErrNetworkFailure, err),
		Retryable: true,
	}
}

// isRetryableStatusCode determines if an HTTP status code indicates a retryable error.
func isRetryableStatusCode(code int) bool {
	switch code {
	case http.StatusTooManyRequests,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusBadGateway:
		return true
	default:
		return false
	}
}

// ConfigError represents configuration-related errors.
type ConfigError struct {
	Key     string // Configuration key that caused the error
	Message string // Descriptive error message
	Err     error  // Underlying error
}

// Error implements the error interface.
func (e *ConfigError) Error() string {
	if e.Key != "" {
		return fmt.Sprintf("config error [%s]: %s", e.Key, e.Message)
	}
	return fmt.Sprintf("config error: %s", e.Message)
}

// Unwrap returns the underlying error for error chain support.
func (e *ConfigError) Unwrap() error {
	return e.Err
}

// NewConfigError creates a new ConfigError with the given key, message, and underlying error.
func NewConfigError(key, message string, err error) *ConfigError {
	return &ConfigError{
		Key:     key,
		Message: message,
		Err:     err,
	}
}

// NewMissingAPIKeyError creates a ConfigError for a missing API key.
func NewMissingAPIKeyError(provider string) *ConfigError {
	return &ConfigError{
		Key:     provider + "_API_KEY",
		Message: fmt.Sprintf("API key for %s is not configured", provider),
		Err:     ErrMissingAPIKey,
	}
}

// NewInvalidConfigError creates a ConfigError for invalid configuration values.
func NewInvalidConfigError(key, reason string) *ConfigError {
	return &ConfigError{
		Key:     key,
		Message: fmt.Sprintf("invalid value: %s", reason),
		Err:     ErrInvalidConfig,
	}
}

// ExportError represents errors during export operations (file writing, formatting).
type ExportError struct {
	Op       string // Operation that failed (e.g., "write", "format", "create")
	Path     string // File path involved
	Format   string // Export format (e.g., "markdown", "json", "pdf")
	Err      error  // Underlying error
}

// Error implements the error interface.
func (e *ExportError) Error() string {
	if e.Format != "" {
		return fmt.Sprintf("export %s [%s]: %s: %v", e.Op, e.Format, e.Path, e.Err)
	}
	return fmt.Sprintf("export %s: %s: %v", e.Op, e.Path, e.Err)
}

// Unwrap returns the underlying error for error chain support.
func (e *ExportError) Unwrap() error {
	return e.Err
}

// NewExportError creates a new ExportError with full context.
func NewExportError(op, path, format string, err error) *ExportError {
	return &ExportError{
		Op:     op,
		Path:   path,
		Format: format,
		Err:    err,
	}
}

// NewFileWriteError creates an ExportError for file write failures.
func NewFileWriteError(path string, err error) *ExportError {
	return &ExportError{
		Op:   "write",
		Path: path,
		Err:  fmt.Errorf("%w: %v", ErrFileWriteFailed, err),
	}
}

// NewFormatError creates an ExportError for formatting failures.
func NewFormatError(format string, err error) *ExportError {
	return &ExportError{
		Op:     "format",
		Format: format,
		Err:    err,
	}
}

// IsRetryable determines if an error should trigger a retry.
// It checks for known retryable error types and conditions.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Check for APIError with Retryable flag
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Retryable
	}

	// Check for specific sentinel errors that are retryable
	if errors.Is(err, ErrRateLimited) || errors.Is(err, ErrNetworkFailure) {
		return true
	}

	// Git errors are generally not retryable (local operations)
	var gitErr *GitError
	if errors.As(err, &gitErr) {
		return false
	}

	// Config errors are not retryable (require user intervention)
	var configErr *ConfigError
	if errors.As(err, &configErr) {
		return false
	}

	// Export errors are generally not retryable (disk/permission issues)
	var exportErr *ExportError
	if errors.As(err, &exportErr) {
		return false
	}

	return false
}

// IsGitError checks if the error is or wraps a GitError.
func IsGitError(err error) bool {
	var gitErr *GitError
	return errors.As(err, &gitErr)
}

// IsAPIError checks if the error is or wraps an APIError.
func IsAPIError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr)
}

// IsConfigError checks if the error is or wraps a ConfigError.
func IsConfigError(err error) bool {
	var configErr *ConfigError
	return errors.As(err, &configErr)
}

// IsExportError checks if the error is or wraps an ExportError.
func IsExportError(err error) bool {
	var exportErr *ExportError
	return errors.As(err, &exportErr)
}

// Wrap wraps an error with additional context message.
// Returns nil if err is nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// Wrapf wraps an error with a formatted context message.
// Returns nil if err is nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}
