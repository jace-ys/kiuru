package user

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/lib/pq"
)

var (
	ErrUserNotFound = errors.New("requested user not found")
	ErrUserExists   = errors.New("account already exists")

	ErrHashingPssword = errors.New("failed to encrypt password")

	ErrInvalidRequest = errors.New("invalid request")
	ErrUnknown        = errors.New("unknown error")
)

func ErrUserExistsContext(pqErr *pq.Error) error {
	duplicateKeyRegex := regexp.MustCompile(`duplicate key value \((\w+)\)`)
	match := duplicateKeyRegex.FindStringSubmatch(pqErr.Error())
	if len(match) != 2 {
		return pqErr
	}
	return fmt.Errorf("%w: %s unavailable", ErrUserExists, match[1])
}
