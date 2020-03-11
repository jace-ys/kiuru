package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/jmoiron/sqlx"
	"github.com/kiuru-travel/airdrop-go/pkg/authr"
	"github.com/kiuru-travel/airdrop-go/pkg/gorpc"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"

	pb "github.com/jace-ys/kiuru/backend/service.auth/api/auth"
)

func (s *authService) GenerateAuthToken(ctx context.Context, req *pb.GenerateAuthTokenRequest) (*pb.GenerateAuthTokenResponse, error) {
	level.Info(s.logger).Log("event", "get_auth_token.started")
	defer level.Info(s.logger).Log("event", "get_auth_token.finished")

	err := s.validateLoginPayload(req)
	if err != nil {
		level.Error(s.logger).Log("event", "get_auth_token.failed", "msg", err)
		return nil, gorpc.Error(codes.InvalidArgument, err)
	}

	userID, hashedPassword, err := s.getLoginUser(ctx, req.Username)
	if err != nil {
		level.Error(s.logger).Log("event", "get_auth_token.failed", "msg", err)
		switch {
		case errors.Is(err, ErrUserNotFound):
			return nil, gorpc.Error(codes.NotFound, err)
		default:
			return nil, gorpc.Error(codes.Internal, err)
		}
	}

	err = s.verifyLoginPassword(hashedPassword, req.Password)
	if err != nil {
		level.Error(s.logger).Log("event", "get_auth_token.failed", "msg", err)
		return nil, gorpc.Error(codes.Unauthenticated, err)
	}

	jwt, err := s.token.GenerateToken(ctx, userID, req.Username)
	if err != nil {
		level.Error(s.logger).Log("event", "get_auth_token.failed", "msg", err)
		return nil, gorpc.Error(codes.Internal, err)
	}

	level.Info(s.logger).Log("event", "get_auth_token.success")
	return &pb.GenerateAuthTokenResponse{
		Token: jwt,
	}, nil
}

func (s *authService) validateLoginPayload(login *pb.GenerateAuthTokenRequest) error {
	switch {
	case login.Username == "":
		return fmt.Errorf("missing \"username\"")
	case login.Password == "":
		return fmt.Errorf("missing \"password\"")
	}
	return nil
}

func (s *authService) getLoginUser(ctx context.Context, username string) (string, string, error) {
	var userID, hashedPassword string
	err := s.db.Transact(ctx, func(tx *sqlx.Tx) error {
		query := `
		SELECT u.id, u.password
		FROM users as u
		WHERE username=$1
		`
		row := tx.QueryRowxContext(ctx, query, username)
		return row.Scan(&userID, &hashedPassword)
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", "", ErrUserNotFound
		default:
			return "", "", err
		}
	}
	return userID, hashedPassword, nil
}

func (s *authService) verifyLoginPassword(hashedPassword, loginPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginPassword)); err != nil {
		return ErrIncorrectPassword
	}
	return nil
}

func (s *authService) RefreshAuthToken(ctx context.Context, req *pb.RefreshAuthTokenRequest) (*pb.RefreshAuthTokenResponse, error) {
	level.Info(s.logger).Log("event", "refresh_auth_token.started")
	defer level.Info(s.logger).Log("event", "refresh_auth_token.finished")

	claims, err := s.token.ValidateToken(ctx, req.Token)
	if err != nil {
		level.Error(s.logger).Log("event", "refresh_auth_token.failed", "msg", err)
		switch {
		case errors.Is(err, authr.ErrTokenInvalid):
			return nil, gorpc.Error(codes.InvalidArgument, err)
		default:
			return nil, gorpc.Error(codes.Internal, err)
		}
	}

	err = s.token.IsRevokedToken(ctx, req.Token)
	if err != nil {
		level.Error(s.logger).Log("event", "refresh_auth_token.failed", "msg", err)
		switch {
		case errors.Is(err, authr.ErrTokenInvalid):
			return nil, gorpc.Error(codes.InvalidArgument, err)
		case errors.Is(err, authr.ErrTokenRevoked):
			return nil, gorpc.Error(codes.Unauthenticated, err)
		default:
			return nil, gorpc.Error(codes.Internal, err)
		}
	}

	err = s.isRefreshable(claims)
	if err != nil {
		level.Error(s.logger).Log("event", "refresh_auth_token.failed", "msg", err)
		return nil, gorpc.Error(codes.ResourceExhausted, err)
	}

	jwt, err := s.token.GenerateToken(ctx, claims.UserMD.Id, claims.UserMD.Username)
	if err != nil {
		level.Error(s.logger).Log("event", "refresh_auth_token.failed", "msg", err)
		return nil, gorpc.Error(codes.Internal, err)
	}

	err = s.token.RevokeToken(ctx, req.Token)
	if err != nil {
		level.Error(s.logger).Log("event", "refresh_auth_token.failed", "msg", err)
		return nil, gorpc.Error(codes.Internal, err)
	}

	level.Info(s.logger).Log("event", "refresh_auth_token.success")
	return &pb.RefreshAuthTokenResponse{
		Token: jwt,
	}, nil
}

func (s *authService) isRefreshable(claims *authr.JWTClaims) error {
	ttl := claims.StandardClaims.ExpiresAt - claims.StandardClaims.IssuedAt
	refreshTime := time.Duration(ttl)
	if time.Unix(claims.StandardClaims.ExpiresAt, 0).Sub(time.Now()) > refreshTime {
		return ErrRefreshRateExceeded
	}
	return nil
}

func (s *authService) RevokeAuthToken(ctx context.Context, req *pb.RevokeAuthTokenRequest) (*pb.RevokeAuthTokenResponse, error) {
	level.Info(s.logger).Log("event", "revoke_auth_token.started")
	defer level.Info(s.logger).Log("event", "revoke_auth_token.finished")

	_, err := s.token.ValidateToken(ctx, req.Token)
	if err != nil {
		level.Error(s.logger).Log("event", "revoke_auth_token.failed", "msg", err)
		switch {
		case errors.Is(err, authr.ErrTokenInvalid):
			return nil, gorpc.Error(codes.InvalidArgument, err)
		default:
			return nil, gorpc.Error(codes.Internal, err)
		}
	}

	err = s.token.RevokeToken(ctx, req.Token)
	if err != nil {
		level.Error(s.logger).Log("event", "revoke_auth_token.failed", "msg", err)
		return nil, gorpc.Error(codes.Internal, err)
	}

	level.Info(s.logger).Log("event", "revoke_auth_token.success")
	return &pb.RevokeAuthTokenResponse{}, nil
}
