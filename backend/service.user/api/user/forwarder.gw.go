package user

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/kru-travel/airdrop-go/pkg/gorpc"
)

func init() {
	forward_UserService_CreateUser_0 = gorpc.Forward(
		gorpc.NewHeaderLocationHandler(createUserHeaderLocationFunc),
	)
}

func createUserHeaderLocationFunc(resp proto.Message) string {
	var user CreateUserResponse
	data, err := proto.Marshal(resp)
	if err != nil {
		return ""
	}
	err = proto.Unmarshal(data, &user)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("/v1/users/%s", user.Id)
}
