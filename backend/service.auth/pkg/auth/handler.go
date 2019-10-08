package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/kru-travel/airdrop-go/pkg/gorpc"
	"github.com/kru-travel/airdrop-go/pkg/slogger"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"

	pb "github.com/jace-ys/kru-travel/backend/service.auth/api/auth"
)

func (s *authService) GetAuthToken(ctx context.Context, req *pb.GetAuthTokenRequest) (*pb.GetAuthTokenResponse, error) {
	slogger.Info().Log("event", "get_auth_token.started", "username", req.Username)
	defer slogger.Info().Log("event", "get_auth_token.finished", "username", req.Username)

	hashedPassword, err := s.getLoginPassword(ctx, req.Username)
	if err != nil {
		slogger.Error().Log("event", "get_auth_token.failed", "username", req.Username, "msg", err)
		switch {
		case errors.Is(err, ErrUserNotFound):
			return nil, gorpc.Error(codes.NotFound, err)
		default:
			return nil, gorpc.InternalError()
		}
	}

	err = s.verifyLoginPassword(hashedPassword, req.Password)
	if err != nil {
		slogger.Error().Log("event", "get_auth_token.failed", "username", req.Username, "msg", err)
		return nil, gorpc.Error(codes.NotFound, err)
	}

	return &pb.GetAuthTokenResponse{}, nil
}

func (s *authService) getLoginPassword(ctx context.Context, username string) (string, error) {
	var hashedPassword string
	err := s.db.Transact(ctx, func(tx *sqlx.Tx) error {
		query := `
		SELECT u.password
		FROM users as u
		WHERE username=$1
		`
		row := tx.QueryRowxContext(ctx, query, username)
		err := row.Scan(&hashedPassword)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrUserNotFound
		case err != nil:
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

func (s *authService) verifyLoginPassword(hashedPassword, loginPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginPassword)); err != nil {
		return ErrIncorrectPassword
	}
	return nil
}

func (s *authService) RefreshAuthToken(ctx context.Context, req *pb.RefreshAuthTokenRequest) (*pb.RefreshAuthTokenResponse, error) {
	return &pb.RefreshAuthTokenResponse{}, nil
}
