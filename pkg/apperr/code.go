package apperr

import (
	"net/http"
)

//go:generate go tool enumer -type=Code -trimprefix Code -output code_enum.go -transform snake-upper -text -json
type Code int

// 400 bad requests
const (
	CodeBadRequest Code = 1000 + iota
)

// 401 unauthorized
const (
	CodeUnauthorized Code = 2000 + iota
)

// 403 forbidden
const (
	CodeForbidden Code = 3000 + iota
)

// 404 not found
const (
	CodeNotFound Code = 4000 + iota
)

// 405 method not allowed
const (
	CodeMethodNotAllowed Code = 5000 + iota
)

// 406 not acceptable
const (
	CodeNotAcceptable Code = 6000 + iota
)

// 409 conflict
const (
	CodeConflict Code = 7000 + iota
)

// 410 gone
const (
	CodeGone Code = 8000 + iota
)

// 413 request entity too large
const (
	CodeRequestEntityTooLarge Code = 9000 + iota
)

// 418 I'm a teapot
const (
	CodeTeapot Code = 10000 + iota
)

// 422 unprocessable entity
const (
	CodeUnprocessableEntity Code = 11000 + iota
)

// 429 too many requests
const (
	CodeTooManyRequests Code = 12000 + iota
)

// 500 internal server error
const (
	CodeInternalServerError Code = 13000 + iota
)

// 501 not implemented
const (
	CodeNotImplemented Code = 14000 + iota
)

// 503 service unavailable
const (
	CodeServiceUnavailable Code = 15000 + iota
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

func (c Code) HTTPStatusCode() int {
	switch {
	case c >= 1000 && c < 2000:
		return http.StatusBadRequest
	case c >= 2000 && c < 3000:
		return http.StatusUnauthorized
	case c >= 3000 && c < 4000:
		return http.StatusForbidden
	case c >= 4000 && c < 5000:
		return http.StatusNotFound
	case c >= 5000 && c < 6000:
		return http.StatusMethodNotAllowed
	case c >= 6000 && c < 7000:
		return http.StatusNotAcceptable
	case c >= 7000 && c < 8000:
		return http.StatusConflict
	case c >= 8000 && c < 9000:
		return http.StatusGone
	case c >= 9000 && c < 10000:
		return http.StatusRequestEntityTooLarge
	case c >= 10000 && c < 11000:
		return http.StatusTeapot
	case c >= 11000 && c < 12000:
		return http.StatusUnprocessableEntity
	case c >= 12000 && c < 13000:
		return http.StatusTooManyRequests
	case c >= 13000 && c < 14000:
		return http.StatusInternalServerError
	case c >= 14000 && c < 15000:
		return http.StatusNotImplemented
	case c >= 15000 && c < 16000:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
