package auth

import (
	"errors"
)

var (
	ErrMissingSecret = errors.New("secret key not provided")
)

var (
	ErrUserNotFound      = errors.New("requested user not found")
	ErrIncorrectPassword = errors.New("incorrect password provided")

	ErrGeneratingToken     = errors.New("failed to generate token")
	ErrInvalidToken        = errors.New("invalid token")
	ErrRefreshRateExceeded = errors.New("refresh rate exceeded")

	ErrInvalidRequest = errors.New("invalid request")
	ErrUnknown        = errors.New("unknown error")
)
