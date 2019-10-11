package auth

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	pb "github.com/jace-ys/kru-travel/backend/service.auth/api/auth"
)

type DBClient interface {
	Connect() error
	Transact(ctx context.Context, fn func(*sqlx.Tx) error) error
	Close() error
}

type Server interface {
	Init(ctx context.Context, server pb.AuthServiceServer) error
	Serve(port int) error
	Shutdown(ctx context.Context) error
}

type authService struct {
	db        DBClient
	jwtConfig JWTConfig
}

func NewService(dbClient DBClient, jwtConfig JWTConfig) (*authService, error) {
	if jwtConfig.SecretKey == "" {
		return nil, fmt.Errorf("failed to create service: %w", ErrMissingSecret)
	}
	return &authService{
		db:        dbClient,
		jwtConfig: jwtConfig,
	}, nil
}

func (s *authService) Init() error {
	if err := s.db.Connect(); err != nil {
		return err
	}
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
