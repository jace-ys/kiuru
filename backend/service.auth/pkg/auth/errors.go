package auth

import (
	"errors"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrIncorrectPassword   = errors.New("incorrect password")
	ErrRefreshRateExceeded = errors.New("refresh rate exceeded")
)
