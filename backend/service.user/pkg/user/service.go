package user

import (
	"context"

	"github.com/kru-travel/airdrop-go/pkg/slogger"
	"github.com/pkg/errors"

	"github.com/jace-ys/kru-travel/backend/service.user/pkg/server"

	pb "github.com/jace-ys/kru-travel/backend/service.user/api/user"
)

type userService struct {
}

func NewService() (*userService, error) {
	u := &userService{}
	if err := u.Init(); err != nil {
		return nil, errors.Wrap(err, "service init")
	}
	return u, nil
}

func (u *userService) Init() error {
	return nil
}

func (u *userService) ListenAndServe(ctx context.Context, s server.Server, port int) error {
	if err := s.Init(u); err != nil {
		return err
	}
	defer s.Shutdown(ctx)
	return s.Serve(ctx, port)
}

func (s *userService) GetUser(ctx context.Context, r *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	slogger.Info().Log("event", "request.get_user", "user_id", r.Id)
	return &pb.GetUserResponse{
		Id: r.Id,
	}, nil
}
