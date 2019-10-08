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
	Init(ctx context.Context, server pb.UserServiceServer) error
	Serve(port int) error
	Shutdown(ctx context.Context) error
}

type userService struct {
	db DbClient
}

func NewService() *userService {
	return &userService{}
}

func (s *userService) Init(dbClient DbClient) error {
	s.db = dbClient
	return nil
}

func (s *userService) StartServer(ctx context.Context, server Server, port int) error {
	if err := server.Init(ctx, s); err != nil {
		return err
	}
	defer server.Shutdown(ctx)
	return server.Serve(port)
}

func (s *userService) Teardown() error {
	return s.db.Close()
}
