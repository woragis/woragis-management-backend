package errors

import (
	"fmt"
)

// AppError represents a structured application error
type AppError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details string                 `json:"details,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
	Err     error                  `json:"-"` // Original error, not exposed in JSON
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError with the given code
func New(code string) *AppError {
	return &AppError{
		Code:    code,
		Message: GetMessage(code),
		Context: make(map[string]interface{}),
	}
}

// NewWithDetails creates a new AppError with additional details
func NewWithDetails(code string, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: GetMessage(code),
		Details: details,
		Context: make(map[string]interface{}),
	}
}

// Wrap creates a new AppError wrapping an existing error
func Wrap(code string, err error) *AppError {
	details := ""
	if err != nil {
		details = err.Error()
	}
	return &AppError{
		Code:    code,
		Message: GetMessage(code),
		Details: details,
		Err:     err,
		Context: make(map[string]interface{}),
	}
}

// WithContext adds context information to the error
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithDetails adds or updates details
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// Is checks if the error matches a specific code
func (e *AppError) Is(code string) bool {
	return e.Code == code
}

// ToHTTPStatus maps error codes to HTTP status codes
func (e *AppError) ToHTTPStatus() int {
	switch {
	// Authentication errors - 401
	case e.Code == AUTH_JWT_INVALID_SIGNATURE,
		e.Code == AUTH_JWT_EXPIRED,
		e.Code == AUTH_JWT_MALFORMED,
		e.Code == AUTH_JWT_MISSING_CLAIMS,
		e.Code == AUTH_TOKEN_MISSING,
		e.Code == AUTH_TOKEN_INVALID_FORMAT,
		e.Code == AUTH_UNAUTHORIZED,
		e.Code == API_KEY_INVALID,
		e.Code == API_KEY_EXPIRED,
		e.Code == API_KEY_REVOKED:
		return 401

	// CSRF errors - 403
	case e.Code == CSRF_TOKEN_EXPIRED,
		e.Code == CSRF_TOKEN_INVALID,
		e.Code == CSRF_TOKEN_MISSING,
		e.Code == CSRF_TOKEN_MISMATCH:
		return 403

	// Not found - 404
	case e.Code == DB_RECORD_NOT_FOUND,
		e.Code == PROJECT_NOT_FOUND,
		e.Code == EXPERIENCE_NOT_FOUND,
		e.Code == LANGUAGE_NOT_FOUND,
		e.Code == CLIENT_NOT_FOUND,
		e.Code == FINANCE_NOT_FOUND,
		e.Code == API_KEY_NOT_FOUND,
		e.Code == TESTIMONIAL_NOT_FOUND:
		return 404

	// Conflict - 409
	case e.Code == DB_DUPLICATE_ENTRY,
		e.Code == PROJECT_SLUG_EXISTS,
		e.Code == EXPERIENCE_EXISTS,
		e.Code == LANGUAGE_EXISTS,
		e.Code == CLIENT_EXISTS,
		e.Code == TESTIMONIAL_EXISTS:
		return 409

	// Validation errors - 400
	case e.Code == VALIDATION_INVALID_INPUT,
		e.Code == VALIDATION_MISSING_FIELD,
		e.Code == VALIDATION_FIELD_TOO_LONG,
		e.Code == VALIDATION_FIELD_TOO_SHORT,
		e.Code == VALIDATION_INVALID_UUID,
		e.Code == VALIDATION_INVALID_DATE,
		e.Code == VALIDATION_INVALID_URL,
		e.Code == VALIDATION_INVALID_EMAIL,
		e.Code == VALIDATION_INVALID_SLUG,
		e.Code == EXPERIENCE_INVALID_DATE_RANGE,
		e.Code == LANGUAGE_INVALID_LEVEL,
		e.Code == FINANCE_INVALID_AMOUNT,
		e.Code == CHAT_INVALID_MESSAGE,
		e.Code == CHAT_CONTEXT_TOO_LARGE:
		return 400

	// Service unavailable - 503
	case e.Code == DB_CONNECTION_FAILED,
		e.Code == REDIS_CONNECTION_FAILED,
		e.Code == AUTH_SERVICE_UNAVAILABLE,
		e.Code == AI_SERVICE_UNAVAILABLE,
		e.Code == SERVER_SERVICE_UNAVAILABLE:
		return 503

	// Timeout - 504
	case e.Code == SERVER_TIMEOUT,
		e.Code == AI_SERVICE_TIMEOUT:
		return 504

	// Too many requests - 429
	case e.Code == SERVER_RATE_LIMIT_EXCEEDED:
		return 429

	// Default to 500 for everything else
	default:
		return 500
	}
}
