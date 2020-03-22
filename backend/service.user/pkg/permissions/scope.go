package permissions

import (
	"github.com/kiuru-travel/airdrop-go/authr"
)

type ScopeFunc func(userMD *authr.UserMD, userID string) bool

func AdminScope(userMD *authr.UserMD, userID string) bool {
	return userMD.Admin
}

func UserScope(userMD *authr.UserMD, userID string) bool {
	switch {
	case userMD.Admin:
		return true
	case userMD.Id == userID:
		return true
	default:
		return false
	}
}
