// Package errors menyediakan error types standar untuk seluruh platform.
// Error codes bersifat global dan konsisten antar service.
package errors

import "fmt"

// Error codes standar.
const (
	// General
	CodeInternal       = "INTERNAL_ERROR"
	CodeValidation     = "VALIDATION_ERROR"
	CodeNotFound       = "NOT_FOUND"
	CodeDuplicate      = "DUPLICATE"
	CodeInvalidContext = "INVALID_CONTEXT"

	// Auth
	CodeUnauthorized = "UNAUTHORIZED"
	CodeForbidden    = "FORBIDDEN"
	CodeTokenExpired = "TOKEN_EXPIRED"

	// Business
	CodeBusinessRule = "BUSINESS_RULE_VIOLATION"
	CodeConflict     = "CONFLICT"

	// Event
	CodeEventPublish = "EVENT_PUBLISH_FAILED"
	CodeEventConsume = "EVENT_CONSUME_FAILED"

	// Integration
	CodeIntegration = "INTEGRATION_ERROR"
	CodeTimeout     = "TIMEOUT"
)

// AppError adalah error standar aplikasi yang membawa code, message, dan details.
type AppError struct {
	Code    string
	Message string
	Details interface{}
	Cause   error
}

// Error mengimplementasikan error interface.
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap mengembalikan cause error (untuk errors.Is/As).
func (e *AppError) Unwrap() error {
	return e.Cause
}

// New membuat AppError baru.
func New(code, message string, details interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Wrap membuat AppError baru dengan cause.
func Wrap(code, message string, cause error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// IsNotFound mengecek apakah error adalah NOT_FOUND.
func IsNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == CodeNotFound
	}
	return false
}

// IsValidation mengecek apakah error adalah VALIDATION_ERROR.
func IsValidation(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == CodeValidation
	}
	return false
}

// IsDuplicate mengecek apakah error adalah DUPLICATE.
func IsDuplicate(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == CodeDuplicate
	}
	return false
}

// === Convenience constructors ===

// NotFound membuat error NOT_FOUND.
func NotFound(entity, id string) *AppError {
	return New(CodeNotFound, fmt.Sprintf("%s dengan ID '%s' tidak ditemukan", entity, id), nil)
}

// Validation membuat error VALIDATION_ERROR.
func Validation(message string, details interface{}) *AppError {
	return New(CodeValidation, message, details)
}

// Duplicate membuat error DUPLICATE.
func Duplicate(entity, field, value string) *AppError {
	return New(CodeDuplicate, fmt.Sprintf("%s dengan %s '%s' sudah ada", entity, field, value), nil)
}

// InvalidContext membuat error INVALID_CONTEXT.
func InvalidContext(message string) *AppError {
	return New(CodeInvalidContext, message, nil)
}

// Internal membuat error INTERNAL_ERROR.
func Internal(message string, cause error) *AppError {
	return Wrap(CodeInternal, message, cause)
}

// BusinessRule membuat error BUSINESS_RULE_VIOLATION.
func BusinessRule(message string) *AppError {
	return New(CodeBusinessRule, message, nil)
}