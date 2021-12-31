package user

import "errors"

var (
	ErrDuplicate          = errors.New("duplicate key")
	ErrPasswordNotCorrect = errors.New("password not correct")
	ErrUserNotFound       = errors.New("user not found")
)
