package user

import (
	"fmt"

	"github.com/golang/protobuf/proto"

	"github.com/kiuru-travel/airdrop-go/pkg/gorpc"
)

func init() {
	forward_UserService_CreateUser_0 = gorpc.Forward(
		gorpc.NewHeaderLocationHandler(createUserHeaderLocationFunc),
	)
}

func createUserHeaderLocationFunc(resp proto.Message) string {
	data, err := proto.Marshal(resp)
	if err != nil {
		return ""
	}

	var user CreateUserResponse
	err = proto.Unmarshal(data, &user)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("/v1/users/%s", user.Id)
}
