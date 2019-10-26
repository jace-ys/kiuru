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

	ErrTokenRevoked        = errors.New("token has been revoked")
	ErrRefreshRateExceeded = errors.New("refresh rate exceeded")
	ErrRevokingToken       = errors.New("failed to revoke token")

	ErrInvalidRequest = errors.New("invalid request")
	ErrUnknown        = errors.New("unknown error")
)
