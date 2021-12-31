package user

import "errors"

var (
	ErrPasswordNotCorrect = errors.New("password not correct")
	ErrUserNotFound       = errors.New("user not found")
)
