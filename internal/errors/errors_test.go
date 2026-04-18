package errors

import (
	"errors"
	"io"
	"net/http"
	"testing"
)

func TestGitError(t *testing.T) {
	t.Run("Error message formatting", func(t *testing.T) {
		err := NewGitError("open", "/path/to/repo", io.EOF)
		expected := "git open: /path/to/repo: EOF"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("Error message without path", func(t *testing.T) {
		err := &GitError{Op: "init", Err: io.EOF}
		expected := "git init: EOF"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("Unwrap returns underlying error", func(t *testing.T) {
		err := NewGitError("open", "/path", io.EOF)
		if !errors.Is(err, io.EOF) {
			t.Error("expected error to wrap io.EOF")
		}
	})

	t.Run("RepositoryNotFoundError", func(t *testing.T) {
		err := NewRepositoryNotFoundError("/missing/repo")
		if !errors.Is(err, ErrRepositoryNotFound) {
			t.Error("expected error to wrap ErrRepositoryNotFound")
		}
		if err.Path != "/missing/repo" {
			t.Errorf("expected path %q, got %q", "/missing/repo", err.Path)
		}
	})

	t.Run("InvalidPathError", func(t *testing.T) {
		err := NewInvalidPathError("/bad/path")
		if !errors.Is(err, ErrInvalidPath) {
			t.Error("expected error to wrap ErrInvalidPath")
		}
	})
}

func TestAPIError(t *testing.T) {
	t.Run("Error message with status code", func(t *testing.T) {
		err := NewAPIError("generate", 500, "internal server error", nil)
		expected := "api generate: status 500: internal server error"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("Error message without status code", func(t *testing.T) {
		err := &APIError{Op: "generate", Message: "timeout"}
		expected := "api generate: timeout"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("Error message with underlying error only", func(t *testing.T) {
		err := &APIError{Op: "generate", Err: io.EOF}
		expected := "api generate: EOF"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("RateLimitError is retryable", func(t *testing.T) {
		err := NewRateLimitError("generate")
		if !err.Retryable {
			t.Error("expected rate limit error to be retryable")
		}
		if err.StatusCode != http.StatusTooManyRequests {
			t.Errorf("expected status %d, got %d", http.StatusTooManyRequests, err.StatusCode)
		}
		if !errors.Is(err, ErrRateLimited) {
			t.Error("expected error to wrap ErrRateLimited")
		}
	})

	t.Run("AuthError is not retryable", func(t *testing.T) {
		err := NewAuthError("generate", "invalid token")
		if err.Retryable {
			t.Error("expected auth error to not be retryable")
		}
		if !errors.Is(err, ErrUnauthorized) {
			t.Error("expected error to wrap ErrUnauthorized")
		}
	})

	t.Run("NetworkError is retryable", func(t *testing.T) {
		err := NewNetworkError("generate", io.EOF)
		if !err.Retryable {
			t.Error("expected network error to be retryable")
		}
		if !errors.Is(err, ErrNetworkFailure) {
			t.Error("expected error to wrap ErrNetworkFailure")
		}
	})

	t.Run("Retryable status codes", func(t *testing.T) {
		retryableCodes := []int{
			http.StatusTooManyRequests,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
			http.StatusBadGateway,
		}
		for _, code := range retryableCodes {
			err := NewAPIError("test", code, "error", nil)
			if !err.Retryable {
				t.Errorf("expected status %d to be retryable", code)
			}
		}
	})

	t.Run("Non-retryable status codes", func(t *testing.T) {
		nonRetryableCodes := []int{
			http.StatusBadRequest,
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusNotFound,
			http.StatusInternalServerError,
		}
		for _, code := range nonRetryableCodes {
			err := NewAPIError("test", code, "error", nil)
			if err.Retryable {
				t.Errorf("expected status %d to not be retryable", code)
			}
		}
	})
}

func TestConfigError(t *testing.T) {
	t.Run("Error message with key", func(t *testing.T) {
		err := NewConfigError("API_KEY", "value is empty", nil)
		expected := "config error [API_KEY]: value is empty"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("Error message without key", func(t *testing.T) {
		err := &ConfigError{Message: "file not found"}
		expected := "config error: file not found"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("MissingAPIKeyError", func(t *testing.T) {
		err := NewMissingAPIKeyError("OPENAI")
		if err.Key != "OPENAI_API_KEY" {
			t.Errorf("expected key %q, got %q", "OPENAI_API_KEY", err.Key)
		}
		if !errors.Is(err, ErrMissingAPIKey) {
			t.Error("expected error to wrap ErrMissingAPIKey")
		}
	})

	t.Run("InvalidConfigError", func(t *testing.T) {
		err := NewInvalidConfigError("TIMEOUT", "must be positive")
		if !errors.Is(err, ErrInvalidConfig) {
			t.Error("expected error to wrap ErrInvalidConfig")
		}
	})
}

func TestExportError(t *testing.T) {
	t.Run("Error message with format", func(t *testing.T) {
		err := NewExportError("write", "/output/resume.md", "markdown", io.EOF)
		expected := "export write [markdown]: /output/resume.md: EOF"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("Error message without format", func(t *testing.T) {
		err := &ExportError{Op: "write", Path: "/output/file.txt", Err: io.EOF}
		expected := "export write: /output/file.txt: EOF"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("FileWriteError", func(t *testing.T) {
		err := NewFileWriteError("/output/file.txt", io.EOF)
		if !errors.Is(err, ErrFileWriteFailed) {
			t.Error("expected error to wrap ErrFileWriteFailed")
		}
		if err.Op != "write" {
			t.Errorf("expected op %q, got %q", "write", err.Op)
		}
	})

	t.Run("FormatError", func(t *testing.T) {
		err := NewFormatError("json", io.EOF)
		if err.Format != "json" {
			t.Errorf("expected format %q, got %q", "json", err.Format)
		}
	})
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		retryable bool
	}{
		{"nil error", nil, false},
		{"rate limit error", NewRateLimitError("test"), true},
		{"network error", NewNetworkError("test", io.EOF), true},
		{"auth error", NewAuthError("test", "invalid"), false},
		{"git error", NewGitError("open", "/path", io.EOF), false},
		{"config error", NewConfigError("KEY", "missing", nil), false},
		{"export error", NewExportError("write", "/path", "md", io.EOF), false},
		{"wrapped rate limit", Wrap(ErrRateLimited, "context"), true},
		{"wrapped network failure", Wrap(ErrNetworkFailure, "context"), true},
		{"service unavailable", NewAPIError("test", http.StatusServiceUnavailable, "unavailable", nil), true},
		{"bad request", NewAPIError("test", http.StatusBadRequest, "bad request", nil), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRetryable(tt.err); got != tt.retryable {
				t.Errorf("IsRetryable() = %v, want %v", got, tt.retryable)
			}
		})
	}
}

func TestErrorTypeCheckers(t *testing.T) {
	t.Run("IsGitError", func(t *testing.T) {
		gitErr := NewGitError("open", "/path", io.EOF)
		wrappedGitErr := Wrap(gitErr, "context")

		if !IsGitError(gitErr) {
			t.Error("expected IsGitError to return true for GitError")
		}
		if !IsGitError(wrappedGitErr) {
			t.Error("expected IsGitError to return true for wrapped GitError")
		}
		if IsGitError(io.EOF) {
			t.Error("expected IsGitError to return false for non-GitError")
		}
	})

	t.Run("IsAPIError", func(t *testing.T) {
		apiErr := NewAPIError("test", 500, "error", nil)
		wrappedAPIErr := Wrap(apiErr, "context")

		if !IsAPIError(apiErr) {
			t.Error("expected IsAPIError to return true for APIError")
		}
		if !IsAPIError(wrappedAPIErr) {
			t.Error("expected IsAPIError to return true for wrapped APIError")
		}
		if IsAPIError(io.EOF) {
			t.Error("expected IsAPIError to return false for non-APIError")
		}
	})

	t.Run("IsConfigError", func(t *testing.T) {
		configErr := NewConfigError("KEY", "message", nil)
		wrappedConfigErr := Wrap(configErr, "context")

		if !IsConfigError(configErr) {
			t.Error("expected IsConfigError to return true for ConfigError")
		}
		if !IsConfigError(wrappedConfigErr) {
			t.Error("expected IsConfigError to return true for wrapped ConfigError")
		}
		if IsConfigError(io.EOF) {
			t.Error("expected IsConfigError to return false for non-ConfigError")
		}
	})

	t.Run("IsExportError", func(t *testing.T) {
		exportErr := NewExportError("write", "/path", "md", nil)
		wrappedExportErr := Wrap(exportErr, "context")

		if !IsExportError(exportErr) {
			t.Error("expected IsExportError to return true for ExportError")
		}
		if !IsExportError(wrappedExportErr) {
			t.Error("expected IsExportError to return true for wrapped ExportError")
		}
		if IsExportError(io.EOF) {
			t.Error("expected IsExportError to return false for non-ExportError")
		}
	})
}

func TestWrapFunctions(t *testing.T) {
	t.Run("Wrap with nil error", func(t *testing.T) {
		if Wrap(nil, "context") != nil {
			t.Error("expected Wrap(nil) to return nil")
		}
	})

	t.Run("Wrap preserves error chain", func(t *testing.T) {
		original := io.EOF
		wrapped := Wrap(original, "context")
		if !errors.Is(wrapped, io.EOF) {
			t.Error("expected wrapped error to preserve original error")
		}
	})

	t.Run("Wrapf with nil error", func(t *testing.T) {
		if Wrapf(nil, "context %s", "test") != nil {
			t.Error("expected Wrapf(nil) to return nil")
		}
	})

	t.Run("Wrapf formats message", func(t *testing.T) {
		wrapped := Wrapf(io.EOF, "operation %s failed", "read")
		expected := "operation read failed: EOF"
		if wrapped.Error() != expected {
			t.Errorf("expected %q, got %q", expected, wrapped.Error())
		}
	})
}

func TestErrorsAs(t *testing.T) {
	t.Run("errors.As with GitError", func(t *testing.T) {
		err := Wrap(NewGitError("open", "/path", io.EOF), "context")
		var gitErr *GitError
		if !errors.As(err, &gitErr) {
			t.Error("expected errors.As to find GitError")
		}
		if gitErr.Path != "/path" {
			t.Errorf("expected path %q, got %q", "/path", gitErr.Path)
		}
	})

	t.Run("errors.As with APIError", func(t *testing.T) {
		err := Wrap(NewAPIError("generate", 429, "rate limited", nil), "context")
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Error("expected errors.As to find APIError")
		}
		if apiErr.StatusCode != 429 {
			t.Errorf("expected status code 429, got %d", apiErr.StatusCode)
		}
	})
}
