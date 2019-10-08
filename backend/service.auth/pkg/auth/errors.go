package auth

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound      = errors.New("requested user not found")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrInvalidRequest    = errors.New("invalid request")
	ErrUnknown           = errors.New("unknown error")
)

func ErrInvalidRequestCtx(errCtx string) error {
	return fmt.Errorf("%w: %s", ErrInvalidRequest, errCtx)
}
