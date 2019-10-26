package user

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"

	pb "github.com/jace-ys/kru-travel/backend/service.user/api/user"
)

var authenticatedMethods = map[string]bool{
	"/user.UserService/GetAllUsers": false,
	"/user.UserService/GetUser":     false,
	"/user.UserService/CreateUser":  false,
	"/user.UserService/DeleteUser":  true,
}

type DBClient interface {
	Connect() error
	Transact(ctx context.Context, fn func(*sqlx.Tx) error) error
	Close() error
}

type Server interface {
	Init(ctx context.Context, server pb.UserServiceServer) error
	Serve() error
	Shutdown(ctx context.Context) error
}

type userService struct {
	authMethods map[string]bool
	logger      log.Logger
	db          DBClient
}

func NewService(logger log.Logger, dbClient DBClient) (*userService, error) {
	return &userService{
		authMethods: authenticatedMethods,
		logger:      logger,
		db:          dbClient,
	}, nil
}

func (s *userService) GetAuthenticatedMethods() map[string]bool {
	return s.authMethods
}

func (s *userService) Init() error {
	if err := s.db.Connect(); err != nil {
		return err
	}
	return nil
}

func (s *userService) StartServer(ctx context.Context, server Server) error {
	if err := server.Init(ctx, s); err != nil {
		return err
	}
	return server.Serve()
}

func (s *userService) Teardown() error {
	return s.db.Close()
}
