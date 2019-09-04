package user

import "errors"

var (
	ErrUserNotFound = errors.New("requested user not found")
)
