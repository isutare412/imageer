package apperr

import (
	"fmt"
	"net/http"
)

type Code struct {
	id   int
	name string
}

var (
	// 400 Bad Request
	CodeBadRequest = Code{1000, "BAD_REQUEST"}

	// 401 Unauthorized
	CodeUnauthorized = Code{2000, "UNAUTHORIZED"}

	// 403 Forbidden
	CodeForbidden = Code{3000, "FORBIDDEN"}

	// 404 Not Found
	CodeNotFound = Code{4000, "NOT_FOUND"}

	// 405 Method Not Allowed
	CodeMethodNotAllowed = Code{5000, "METHOD_NOT_ALLOWED"}

	// 409 Conflict
	CodeConflict = Code{6000, "CONFLICT"}

	// 413 Entity Too Large
	CodeEntityTooLarge = Code{7000, "ENTITY_TOO_LARGE"}

	// 422 Unprocessable Entity
	CodeUnprocessableEntity = Code{8000, "UNPROCESSABLE_ENTITY"}

	// 500 Internal Server Error
	CodeInternalServerError = Code{9000, "INTERNAL_SERVER_ERROR"}

	// 501 Not Implemented
	CodeNotImplemented = Code{10000, "NOT_IMPLEMENTED"}

	// 503 Service Unavailable
	CodeServiceUnavailable = Code{11000, "SERVICE_UNAVAILABLE"}
)

func DefaultCode(statusCode int) Code {
	switch statusCode {
	case http.StatusBadRequest:
		return CodeBadRequest
	case http.StatusUnauthorized:
		return CodeUnauthorized
	case http.StatusForbidden:
		return CodeForbidden
	case http.StatusNotFound:
		return CodeNotFound
	case http.StatusMethodNotAllowed:
		return CodeMethodNotAllowed
	case http.StatusConflict:
		return CodeConflict
	case http.StatusRequestEntityTooLarge:
		return CodeEntityTooLarge
	case http.StatusUnprocessableEntity:
		return CodeUnprocessableEntity
	case http.StatusInternalServerError:
		return CodeInternalServerError
	case http.StatusNotImplemented:
		return CodeNotImplemented
	case http.StatusServiceUnavailable:
		return CodeServiceUnavailable
	default:
		return CodeInternalServerError
	}
}

func (c Code) ID() int      { return c.id }
func (c Code) Name() string { return c.name }

func (c Code) HTTPStatusCode() int {
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
		return http.StatusConflict
	case c.id >= 7000 && c.id < 8000:
		return http.StatusRequestEntityTooLarge
	case c.id >= 8000 && c.id < 9000:
		return http.StatusUnprocessableEntity
	case c.id >= 9000 && c.id < 10000:
		return http.StatusInternalServerError
	case c.id >= 10000 && c.id < 11000:
		return http.StatusNotImplemented
	case c.id >= 11000 && c.id < 12000:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

func (c Code) String() string {
	return fmt.Sprintf("%s(%d)", c.name, c.id)
}
