package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/jmoiron/sqlx"
	"github.com/kiuru-travel/airdrop-go/authr"
	"github.com/kiuru-travel/airdrop-go/gorpc"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"

	pb "github.com/jace-ys/kiuru/backend/service.auth/api/auth"
)

func (s *AuthService) GenerateToken(ctx context.Context, req *pb.GenerateTokenRequest) (*pb.GenerateTokenResponse, error) {
	level.Info(s.logger).Log("event", "generate_token.started")
	defer level.Info(s.logger).Log("event", "generate_token.finished")

	err := s.validateLoginPayload(req)
	if err != nil {
		level.Error(s.logger).Log("event", "generate_token.failure", "msg", err)
		return nil, gorpc.Error(codes.InvalidArgument, err)
	}

	userID, hashedPassword, err := s.getLoginUser(ctx, req.Username)
	if err != nil {
		level.Error(s.logger).Log("event", "generate_token.failure", "msg", err)
		switch {
		case errors.Is(err, ErrUserNotFound):
			return nil, gorpc.Error(codes.NotFound, err)
		default:
			return nil, gorpc.Error(codes.Internal, err)
		}
	}

	err = s.verifyLoginPassword(hashedPassword, req.Password)
	if err != nil {
		level.Error(s.logger).Log("event", "generate_token.failure", "msg", err)
		return nil, gorpc.Error(codes.Unauthenticated, err)
	}

	token, err := s.token.GenerateToken(ctx, userID, req.Username)
	if err != nil {
		level.Error(s.logger).Log("event", "generate_token.failure", "msg", err)
		return nil, gorpc.Error(codes.Internal, err)
	}

	level.Info(s.logger).Log("event", "generate_token.success")
	return &pb.GenerateTokenResponse{
		Token: token,
	}, nil
}

func (s *AuthService) validateLoginPayload(login *pb.GenerateTokenRequest) error {
	switch {
	case login.Username == "":
		return fmt.Errorf("invalid username")
	case login.Password == "":
		return fmt.Errorf("invalid password")
	}
	return nil
}

func (s *AuthService) getLoginUser(ctx context.Context, username string) (string, string, error) {
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

func (s *AuthService) verifyLoginPassword(hashedPassword, loginPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginPassword)); err != nil {
		return ErrPasswordIncorrect
	}
	return nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	level.Info(s.logger).Log("event", "refresh_token.started")
	defer level.Info(s.logger).Log("event", "refresh_token.finished")

	claims, err := s.token.ValidateToken(ctx, req.Token)
	if err != nil {
		level.Error(s.logger).Log("event", "refresh_token.failure", "msg", err)
		switch {
		case errors.Is(err, authr.ErrTokenInvalid):
			return nil, gorpc.Error(codes.InvalidArgument, err)
		default:
			return nil, gorpc.Error(codes.Internal, err)
		}
	}

	err = s.token.IsRevokedToken(ctx, req.Token)
	if err != nil {
		level.Error(s.logger).Log("event", "refresh_token.failure", "msg", err)
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
		level.Error(s.logger).Log("event", "refresh_token.failure", "msg", err)
		return nil, gorpc.Error(codes.ResourceExhausted, err)
	}

	token, err := s.token.GenerateToken(ctx, claims.UserMD.Id, claims.UserMD.Username)
	if err != nil {
		level.Error(s.logger).Log("event", "refresh_token.failure", "msg", err)
		return nil, gorpc.Error(codes.Internal, err)
	}

	err = s.token.RevokeToken(ctx, req.Token)
	if err != nil {
		level.Error(s.logger).Log("event", "refresh_token.failure", "msg", err)
		return nil, gorpc.Error(codes.Internal, err)
	}

	level.Info(s.logger).Log("event", "refresh_token.success")
	return &pb.RefreshTokenResponse{
		Token: token,
	}, nil
}

func (s *AuthService) isRefreshable(claims *authr.JWTClaims) error {
	ttl := time.Duration(claims.StandardClaims.ExpiresAt - claims.StandardClaims.IssuedAt)
	if time.Unix(claims.StandardClaims.ExpiresAt, 0).Sub(time.Now()) > ttl {
		return ErrRefreshRateExceeded
	}
	return nil
}

func (s *AuthService) RevokeToken(ctx context.Context, req *pb.RevokeTokenRequest) (*pb.RevokeTokenResponse, error) {
	level.Info(s.logger).Log("event", "revoke_token.started")
	defer level.Info(s.logger).Log("event", "revoke_token.finished")

	_, err := s.token.ValidateToken(ctx, req.Token)
	if err != nil {
		level.Error(s.logger).Log("event", "revoke_token.failure", "msg", err)
		switch {
		case errors.Is(err, authr.ErrTokenInvalid):
			return nil, gorpc.Error(codes.InvalidArgument, err)
		default:
			return nil, gorpc.Error(codes.Internal, err)
		}
	}

	err = s.token.RevokeToken(ctx, req.Token)
	if err != nil {
		level.Error(s.logger).Log("event", "revoke_token.failure", "msg", err)
		return nil, gorpc.Error(codes.Internal, err)
	}

	level.Info(s.logger).Log("event", "revoke_token.success")
	return &pb.RevokeTokenResponse{}, nil
}
