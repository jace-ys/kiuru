package user

import (
	"context"

	"google.golang.org/grpc"

	"github.com/kru-travel/airdrop-go/pkg/slogger"

	pb "github.com/jace-ys/kru-travel/backend/service.user/api/grpc/user"
)

type userService struct {
	*grpc.Server
}

func NewService() *userService {
	opts := []grpc.ServerOption{}
	return &userService{grpc.NewServer(opts...)}
}

func (s *userService) Init() error {
	pb.RegisterUserServiceServer(s.Server, s)
	return nil
}

func (s *userService) Shutdown() {
	s.Server.Stop()
}

func (s *userService) GetUser(ctx context.Context, r *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	slogger.Info().Log("event", "request.get_user", "user_id", r.Id)
	return &pb.GetUserResponse{
		Id: r.Id,
	}, nil
}
