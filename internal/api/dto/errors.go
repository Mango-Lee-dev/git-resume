package dto

import "net/http"

// Error codes
const (
	ErrCodeBadRequest      = "BAD_REQUEST"
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeValidation      = "VALIDATION_ERROR"
	ErrCodeUnauthorized    = "UNAUTHORIZED"
	ErrCodeInternal        = "INTERNAL_ERROR"
	ErrCodeJobNotFound     = "JOB_NOT_FOUND"
	ErrCodeRepositoryError = "REPOSITORY_ERROR"
	ErrCodeAnalysisError   = "ANALYSIS_ERROR"
	ErrCodeExportError     = "EXPORT_ERROR"
)

// APIError represents a structured error response
type APIError struct {
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Details    map[string]string `json:"details,omitempty"`
	RequestID  string            `json:"request_id,omitempty"`
	StatusCode int               `json:"-"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return e.Message
}

// NewBadRequestError creates a 400 error
func NewBadRequestError(message string) *APIError {
	return &APIError{
		Code:       ErrCodeBadRequest,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NewNotFoundError creates a 404 error
func NewNotFoundError(resource string) *APIError {
	return &APIError{
		Code:       ErrCodeNotFound,
		Message:    resource + " not found",
		StatusCode: http.StatusNotFound,
	}
}

// NewValidationError creates a validation error with details
func NewValidationError(details map[string]string) *APIError {
	return &APIError{
		Code:       ErrCodeValidation,
		Message:    "validation failed",
		Details:    details,
		StatusCode: http.StatusBadRequest,
	}
}

// NewInternalError creates a 500 error
func NewInternalError(message string) *APIError {
	return &APIError{
		Code:       ErrCodeInternal,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

// NewJobNotFoundError creates a job not found error
func NewJobNotFoundError(jobID string) *APIError {
	return &APIError{
		Code:       ErrCodeJobNotFound,
		Message:    "job not found: " + jobID,
		StatusCode: http.StatusNotFound,
	}
}

// NewRepositoryError creates a repository error
func NewRepositoryError(message string) *APIError {
	return &APIError{
		Code:       ErrCodeRepositoryError,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NewAnalysisError creates an analysis error
func NewAnalysisError(message string) *APIError {
	return &APIError{
		Code:       ErrCodeAnalysisError,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

// NewExportError creates an export error
func NewExportError(message string) *APIError {
	return &APIError{
		Code:       ErrCodeExportError,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

// WithRequestID adds request ID to the error
func (e *APIError) WithRequestID(requestID string) *APIError {
	e.RequestID = requestID
	return e
}
