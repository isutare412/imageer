package auth

import "errors"

var (
	ErrTokenExpired       = errors.New("token expired")
	ErrCtxSessionNotFound = errors.New("session not found in context")
	ErrCtxInvalidSession  = errors.New("invalid type of session in context")
)
