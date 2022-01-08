package auth

import "errors"

var (
	ErrTokenExpired  = errors.New("token expired")
	ErrCtxIDNotFound = errors.New("ID not found in context")
	ErrCtxInvalidID  = errors.New("invalid type of ID in context")
)
