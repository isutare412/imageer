package apperr

import (
	"fmt"
	"net/http"
)

type Code struct {
	id   int
	name string
}

// 400 Bad Request
var (
	CodeBadRequest = Code{1000, "BAD_REQUEST"}
)

// 401 Unauthorized
var (
	CodeUnauthorized = Code{2000, "UNAUTHORIZED"}
)

// 403 Forbidden
var (
	CodeForbidden = Code{3000, "FORBIDDEN"}
)

// 404 Not Found
var (
	CodeNotFound = Code{4000, "NOT_FOUND"}
)

// 405 Method Not Allowed
var (
	CodeMethodNotAllowed = Code{5000, "METHOD_NOT_ALLOWED"}
)

// 406 Not Acceptable
var (
	CodeNotAcceptable = Code{6000, "NOT_ACCEPTABLE"}
)

// 409 Conflict
var (
	CodeConflict = Code{7000, "CONFLICT"}
)

// 410 Gone
var (
	CodeGone = Code{8000, "GONE"}
)

// 413 Request Entity Too Large
var (
	CodeRequestEntityTooLarge = Code{9000, "REQUEST_ENTITY_TOO_LARGE"}
)

// 418 I'm a Teapot
var (
	CodeTeapot = Code{10000, "IM_A_TEAPOT"}
)

// 422 Unprocessable Entity
var (
	CodeUnprocessableEntity = Code{11000, "UNPROCESSABLE_ENTITY"}
)

// 429 Too Many Requests
var (
	CodeTooManyRequests = Code{12000, "TOO_MANY_REQUESTS"}
)

// 500 Internal Server Error
var (
	CodeInternalServerError = Code{13000, "INTERNAL_SERVER_ERROR"}
)

// 501 Not Implemented
var (
	CodeNotImplemented = Code{14000, "NOT_IMPLEMENTED"}
)

// 503 Service Unavailable
var (
	CodeServiceUnavailable = Code{15000, "SERVICE_UNAVAILABLE"}
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
	case http.StatusNotAcceptable:
		return CodeNotAcceptable
	case http.StatusConflict:
		return CodeConflict
	case http.StatusGone:
		return CodeGone
	case http.StatusRequestEntityTooLarge:
		return CodeRequestEntityTooLarge
	case http.StatusTeapot:
		return CodeTeapot
	case http.StatusUnprocessableEntity:
		return CodeUnprocessableEntity
	case http.StatusTooManyRequests:
		return CodeTooManyRequests
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

func (c Code) String() string {
	return fmt.Sprintf("%s(%d)", c.name, c.id)
}
