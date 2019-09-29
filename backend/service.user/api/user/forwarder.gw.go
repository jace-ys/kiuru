package user

import (
	"github.com/kru-travel/airdrop-go/pkg/gorpc"
)

func init() {
	forward_UserService_DeleteUser_0 = gorpc.Forward(gorpc.StatusNoContentHandler)
}
