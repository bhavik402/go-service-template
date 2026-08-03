// Package apperror defines typed application errors.
//
// Services return these instead of raw errors so that handlers (and any
// other entry point: a Kafka consumer, a cron job) can map them to the
// right transport-level response without needing to know persistence or
// client-library details. Named apperror rather than errors to avoid
// colliding with the standard library package of the same name.
package apperror

import "fmt"

// Code identifies the class of failure, independent of transport.
type Code string

const (
	CodeNotFound     Code = "NOT_FOUND"
	CodeConflict     Code = "CONFLICT"
	CodeInvalidInput Code = "INVALID_INPUT"
	CodeUnauthorized Code = "UNAUTHORIZED"
	CodeForbidden    Code = "FORBIDDEN"
	CodeInternal     Code = "INTERNAL"
)

// Error is the application-level error type. It wraps an optional
// underlying cause so callers can still use errors.Is/errors.As/errors.Unwrap.
type Error struct {
	Code    Code
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *Error) Unwrap() error { return e.Err }

func NotFound(message string) *Error {
	return &Error{Code: CodeNotFound, Message: message}
}

func Conflict(message string) *Error {
	return &Error{Code: CodeConflict, Message: message}
}

func InvalidInput(message string) *Error {
	return &Error{Code: CodeInvalidInput, Message: message}
}

func Unauthorized(message string) *Error {
	return &Error{Code: CodeUnauthorized, Message: message}
}

func Forbidden(message string) *Error {
	return &Error{Code: CodeForbidden, Message: message}
}

// Internal wraps an unexpected error. The underlying err is logged but
// should never be exposed to the client directly.
func Internal(message string, err error) *Error {
	return &Error{Code: CodeInternal, Message: message, Err: err}
}
