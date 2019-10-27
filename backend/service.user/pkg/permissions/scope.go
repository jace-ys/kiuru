package permissions

import (
	"github.com/kru-travel/airdrop-go/pkg/authr"
)

type ScopeFunc func(userMD *authr.UserMD, param interface{}) bool

func AdminScope(userMD *authr.UserMD, param interface{}) bool {
	return userMD.Admin
}

func UserScope(userMD *authr.UserMD, param interface{}) bool {
	switch {
	case userMD.Admin:
		return true
	case userMD.Id == param:
		return true
	default:
		return false
	}
}
