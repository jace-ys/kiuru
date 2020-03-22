package user

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/kiuru-travel/airdrop-go/crdb"

	pb "github.com/jace-ys/kiuru/backend/service.user/api/user"
)

var authenticatedMethods = map[string]bool{
	"/user.UserService/ListUsers":  false,
	"/user.UserService/GetUser":    false,
	"/user.UserService/CreateUser": false,
	"/user.UserService/DeleteUser": true,
}

type Server interface {
	Init(ctx context.Context, server pb.UserServiceServer) error
	Serve() error
	Shutdown(ctx context.Context) error
}

type UserService struct {
	authMethods map[string]bool
	logger      log.Logger
	db          crdb.Client
}

func NewService(logger log.Logger, dbClient crdb.Client) (*UserService, error) {
	return &UserService{
		authMethods: authenticatedMethods,
		logger:      logger,
		db:          dbClient,
	}, nil
}

func (s *UserService) GetAuthenticatedMethods() map[string]bool {
	return s.authMethods
}

func (s *UserService) StartServer(ctx context.Context, server Server) error {
	if err := server.Init(ctx, s); err != nil {
		return err
	}
	if err := server.Serve(); err != nil {
		return err
	}
	return nil
}

func (s *UserService) Teardown() error {
	if err := s.db.Close(); err != nil {
		return err
	}
	return nil
}
