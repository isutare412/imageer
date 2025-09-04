package apperr

import (
	"errors"
	"fmt"
	"strings"
)

type Error struct {
	Code    Code
	Stack   *StackTrace
	Cause   error  // optional
	Summary string // optional
	Detail  string // optional
}

func NewError(code Code) *Error {
	return &Error{
		Code:  code,
		Stack: NewStackTrace(3),
	}
}

func (e *Error) WithCause(err error) *Error {
	e.Cause = err
	return e
}

func (e *Error) WithSummary(format string, args ...any) *Error {
	e.Summary = fmt.Sprintf(format, args...)
	return e
}

func (e *Error) WithDetail(format string, args ...any) *Error {
	e.Detail = fmt.Sprintf(format, args...)
	return e
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func (e *Error) Error() string {
	var b strings.Builder
	b.WriteString("code=")
	b.WriteString(e.Code.String())

	if e.Summary != "" {
		b.WriteString(", summary=")
		b.WriteString(e.Summary)
	}

	if e.Detail != "" {
		b.WriteString(", detail=")
		b.WriteString(e.Detail)
	}

	if e.Cause != nil {
		b.WriteString(", cause=")
		b.WriteString(e.Cause.Error())
	}

	return b.String()
}

func (e *Error) ClientMessage() string {
	switch {
	case e.Summary != "" && e.Detail != "":
		return fmt.Sprintf("%s\n%s", e.Summary, e.Detail)
	case e.Summary != "":
		return e.Summary
	case e.Detail != "":
		return e.Detail
	default:
		return ""
	}
}

// AsError checks if the error is of type *Error and returns it if so.
func AsError(err error) (*Error, bool) {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// IsErrorCode checks if the error is of type *Error and matches the given Code.
func IsErrorCode(err error, code Code) bool {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

// IsErrorStatusCode checks if the error is of type *Error and matches the given
// HTTP status code.
func IsErrorStatusCode(err error, statusCode int) bool {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Code.HTTPStatusCode() == statusCode
	}
	return false
}
