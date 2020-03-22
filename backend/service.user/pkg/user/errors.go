package user

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/lib/pq"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserExists       = errors.New("user already exists")
	ErrPermissionDenied = errors.New("permission denied")
)

func ErrUserExistsContext(pqErr *pq.Error) error {
	duplicateKeyRegex := regexp.MustCompile(`duplicate key value \((\w+)\)`)
	match := duplicateKeyRegex.FindStringSubmatch(pqErr.Error())
	if len(match) != 2 {
		return pqErr
	}
	return fmt.Errorf("%s unavailable: %w", match[1], ErrUserExists)
}
