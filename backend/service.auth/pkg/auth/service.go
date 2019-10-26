package auth

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"github.com/kru-travel/airdrop-go/pkg/authr"

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
	Serve() error
	Shutdown(ctx context.Context) error
}

type authService struct {
	logger log.Logger
	db     DBClient
	redis  RedisClient
	jwt    authr.JWTConfig
}

func NewService(logger log.Logger, dbClient DBClient, redisClient RedisClient, jwtConfig authr.JWTConfig) (*authService, error) {
	if jwtConfig.SecretKey == "" {
		return nil, fmt.Errorf("could not create service: %w", ErrMissingSecret)
	}
	return &authService{
		logger: logger,
		db:     dbClient,
		redis:  redisClient,
		jwt:    jwtConfig,
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

func (s *authService) StartServer(ctx context.Context, server Server) error {
	if err := server.Init(ctx, s); err != nil {
		return err
	}
	if err := server.Serve(); err != nil {
		return err
	}
	return nil
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
