package auth

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/kiuru-travel/airdrop-go/pkg/authr"
	"github.com/kiuru-travel/airdrop-go/pkg/crdb"

	pb "github.com/jace-ys/kiuru/backend/service.auth/api/auth"
)

var authenticatedMethods = map[string]bool{}

type Server interface {
	Init(ctx context.Context, server pb.AuthServiceServer) error
	Serve() error
	Shutdown(ctx context.Context) error
}

type authService struct {
	authMethods map[string]bool
	logger      log.Logger
	db          crdb.Client
	token       authr.JWTHandler
}

func NewService(logger log.Logger, dbClient crdb.Client, tokenHandler authr.JWTHandler) (*authService, error) {
	return &authService{
		authMethods: authenticatedMethods,
		logger:      logger,
		db:          dbClient,
		token:       tokenHandler,
	}, nil
}

func (s *authService) GetAuthenticatedMethods() map[string]bool {
	return s.authMethods
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
	return nil
}
