package auth

import (
	"context"

	"github.com/jmoiron/sqlx"

	pb "github.com/jace-ys/kru-travel/backend/service.auth/api/auth"
)

type DbClient interface {
	Transact(ctx context.Context, fn func(*sqlx.Tx) error) error
	Close() error
}

type Server interface {
	Init(ctx context.Context, server pb.AuthServiceServer) error
	Serve(port int) error
	Shutdown(ctx context.Context) error
}

type authService struct {
	db DbClient
}

func NewService() *authService {
	return &authService{}
}

func (s *authService) Init(dbClient DbClient) error {
	s.db = dbClient
	return nil
}

func (s *authService) StartServer(ctx context.Context, server Server, port int) error {
	if err := server.Init(ctx, s); err != nil {
		return err
	}
	defer server.Shutdown(ctx)
	return server.Serve(port)
}

func (s *authService) Teardown() error {
	return s.db.Close()
}
