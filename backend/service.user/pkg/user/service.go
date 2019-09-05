package user

import (
	"context"

	"github.com/jmoiron/sqlx"

	pb "github.com/jace-ys/kru-travel/backend/service.user/api/user"
)

type DbClient interface {
	Transact(ctx context.Context, fn func(*sqlx.Tx) error) error
	Close() error
}

type Server interface {
	Init(ctx context.Context, s pb.UserServiceServer) error
	Serve(port int) error
	Shutdown(ctx context.Context) error
}

type userService struct {
	db DbClient
}

func NewService() *userService {
	return &userService{}
}

func (u *userService) Init(dbClient DbClient) error {
	u.db = dbClient
	return nil
}

func (u *userService) StartServer(ctx context.Context, s Server, port int) error {
	if err := s.Init(ctx, u); err != nil {
		return err
	}
	defer s.Shutdown(ctx)
	return s.Serve(port)
}

func (u *userService) Teardown() error {
	return u.db.Close()
}
