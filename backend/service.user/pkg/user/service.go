package user

import (
	"context"

	"github.com/kru-travel/airdrop-go/pkg/slogger"

	"github.com/jace-ys/kru-travel/backend/service.user/pkg/server"

	pb "github.com/jace-ys/kru-travel/backend/service.user/api/user"
)

type userService struct {
}

func NewService() *userService {
	return &userService{}
}

func (u *userService) Init() error {
	slogger.Info().Log("event", "service.initialise")
	return nil
}

func (u *userService) StartServer(ctx context.Context, s server.Server, port int) error {
	if err := s.Init(ctx, u); err != nil {
		return err
	}
	defer s.Shutdown(ctx)
	return s.Serve(port)
}

func (u *userService) Teardown() {
	slogger.Info().Log("event", "service.teardown")
}

func (s *userService) GetUser(ctx context.Context, r *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	slogger.Info().Log("event", "request.get_user", "user_id", r.Id)
	return &pb.GetUserResponse{
		Id: r.Id,
	}, nil
}