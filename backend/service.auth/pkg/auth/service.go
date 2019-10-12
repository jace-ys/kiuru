package auth

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"

	pb "github.com/jace-ys/kru-travel/backend/service.auth/api/auth"
)

type DBClient interface {
	Connect() error
	Transact(ctx context.Context, fn func(*sqlx.Tx) error) error
	Close() error
}

type RedisClient interface {
	Connect() error
	Transact(ctx context.Context, fn func(redis.Conn) error) error
	Close() error
}

type Server interface {
	Init(ctx context.Context, server pb.AuthServiceServer) error
	Serve(port int) error
	Shutdown(ctx context.Context) error
}

type authService struct {
	db        DBClient
	redis     RedisClient
	jwtConfig JWTConfig
}

func NewService(dbClient DBClient, redisClient RedisClient, jwtConfig JWTConfig) (*authService, error) {
	if jwtConfig.SecretKey == "" {
		return nil, fmt.Errorf("failed to create service: %w", ErrMissingSecret)
	}
	return &authService{
		db:        dbClient,
		redis:     redisClient,
		jwtConfig: jwtConfig,
	}, nil
}

func (s *authService) Init() error {
	if err := s.db.Connect(); err != nil {
		return err
	}
	if err := s.redis.Connect(); err != nil {
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
	if err := s.db.Close(); err != nil {
		return err
	}
	if err := s.redis.Close(); err != nil {
		return err
	}
	return nil
}
