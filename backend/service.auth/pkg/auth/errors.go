package auth

import (
	"errors"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrPasswordIncorrect   = errors.New("password incorrect")
	ErrRefreshRateExceeded = errors.New("refresh rate exceeded")
)
