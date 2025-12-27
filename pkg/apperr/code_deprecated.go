package apperr

import (
	"fmt"
	"net/http"
)

type CodeDeprecated struct {
	id   int
	name string
}

var (
	// 400 Bad Request
	CodeDeprecatedBadRequest = CodeDeprecated{1000, "BAD_REQUEST"}

	// 401 Unauthorized
	CodeDeprecatedUnauthorized = CodeDeprecated{2000, "UNAUTHORIZED"}

	// 403 Forbidden
	CodeDeprecatedForbidden = CodeDeprecated{3000, "FORBIDDEN"}

	// 404 Not Found
	CodeDeprecatedNotFound = CodeDeprecated{4000, "NOT_FOUND"}

	// 405 Method Not Allowed
	CodeDeprecatedMethodNotAllowed = CodeDeprecated{5000, "METHOD_NOT_ALLOWED"}

	// 406 Not Acceptable
	CodeDeprecatedNotAcceptable = CodeDeprecated{6000, "NOT_ACCEPTABLE"}

	// 409 Conflict
	CodeDeprecatedConflict = CodeDeprecated{7000, "CONFLICT"}

	// 410 Gone
	CodeDeprecatedGone = CodeDeprecated{8000, "GONE"}

	// 413 Request Entity Too Large
	CodeDeprecatedRequestEntityTooLarge = CodeDeprecated{9000, "REQUEST_ENTITY_TOO_LARGE"}

	// 418 I'm a Teapot
	CodeDeprecatedTeapot = CodeDeprecated{10000, "IM_A_TEAPOT"}

	// 422 Unprocessable Entity
	CodeDeprecatedUnprocessableEntity = CodeDeprecated{11000, "UNPROCESSABLE_ENTITY"}

	// 429 Too Many Requests
	CodeDeprecatedTooManyRequests = CodeDeprecated{12000, "TOO_MANY_REQUESTS"}

	// 500 Internal Server Error
	CodeDeprecatedInternalServerError = CodeDeprecated{13000, "INTERNAL_SERVER_ERROR"}

	// 501 Not Implemented
	CodeDeprecatedNotImplemented = CodeDeprecated{14000, "NOT_IMPLEMENTED"}

	// 503 Service Unavailable
	CodeDeprecatedServiceUnavailable = CodeDeprecated{15000, "SERVICE_UNAVAILABLE"}
)

func DefaultCodeDeprecated(statusCode int) CodeDeprecated {
	switch statusCode {
	case http.StatusBadRequest:
		return CodeDeprecatedBadRequest
	case http.StatusUnauthorized:
		return CodeDeprecatedUnauthorized
	case http.StatusForbidden:
		return CodeDeprecatedForbidden
	case http.StatusNotFound:
		return CodeDeprecatedNotFound
	case http.StatusMethodNotAllowed:
		return CodeDeprecatedMethodNotAllowed
	case http.StatusNotAcceptable:
		return CodeDeprecatedNotAcceptable
	case http.StatusConflict:
		return CodeDeprecatedConflict
	case http.StatusGone:
		return CodeDeprecatedGone
	case http.StatusRequestEntityTooLarge:
		return CodeDeprecatedRequestEntityTooLarge
	case http.StatusTeapot:
		return CodeDeprecatedTeapot
	case http.StatusUnprocessableEntity:
		return CodeDeprecatedUnprocessableEntity
	case http.StatusTooManyRequests:
		return CodeDeprecatedTooManyRequests
	case http.StatusInternalServerError:
		return CodeDeprecatedInternalServerError
	case http.StatusNotImplemented:
		return CodeDeprecatedNotImplemented
	case http.StatusServiceUnavailable:
		return CodeDeprecatedServiceUnavailable
	default:
		return CodeDeprecatedInternalServerError
	}
}

func (c CodeDeprecated) ID() int      { return c.id }
func (c CodeDeprecated) Name() string { return c.name }

func (c CodeDeprecated) HTTPStatusCode() int {
	switch {
	case c.id >= 1000 && c.id < 2000:
		return http.StatusBadRequest
	case c.id >= 2000 && c.id < 3000:
		return http.StatusUnauthorized
	case c.id >= 3000 && c.id < 4000:
		return http.StatusForbidden
	case c.id >= 4000 && c.id < 5000:
		return http.StatusNotFound
	case c.id >= 5000 && c.id < 6000:
		return http.StatusMethodNotAllowed
	case c.id >= 6000 && c.id < 7000:
		return http.StatusNotAcceptable
	case c.id >= 7000 && c.id < 8000:
		return http.StatusConflict
	case c.id >= 8000 && c.id < 9000:
		return http.StatusGone
	case c.id >= 9000 && c.id < 10000:
		return http.StatusRequestEntityTooLarge
	case c.id >= 10000 && c.id < 11000:
		return http.StatusTeapot
	case c.id >= 11000 && c.id < 12000:
		return http.StatusUnprocessableEntity
	case c.id >= 12000 && c.id < 13000:
		return http.StatusTooManyRequests
	case c.id >= 13000 && c.id < 14000:
		return http.StatusInternalServerError
	case c.id >= 14000 && c.id < 15000:
		return http.StatusNotImplemented
	case c.id >= 15000 && c.id < 16000:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

func (c CodeDeprecated) String() string {
	return fmt.Sprintf("%s(%d)", c.name, c.id)
}
