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
	ErrUnknown      = errors.New("unknown error")
)

func ErrUserExistsCtx(pqErr *pq.Error) error {
	duplicateKeyRegex := regexp.MustCompile(`duplicate key value \((\w+)\)`)
	match := duplicateKeyRegex.FindStringSubmatch(pqErr.Error())
	if len(match) != 2 {
		return fmt.Errorf("%w: %w", ErrUnknown, pqErr)
	}
	return fmt.Errorf("%w: %s taken", ErrUserExists, match[1])
}
