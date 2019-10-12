package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/kru-travel/airdrop-go/pkg/gorpc"
	"github.com/kru-travel/airdrop-go/pkg/slogger"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"

	pb "github.com/jace-ys/kru-travel/backend/service.auth/api/auth"
)

func (s *authService) GenerateAuthToken(ctx context.Context, req *pb.GenerateAuthTokenRequest) (*pb.GenerateAuthTokenResponse, error) {
	slogger.Info().Log("event", "get_auth_token.started")
	defer slogger.Info().Log("event", "get_auth_token.finished")

	err := s.validateLoginPayload(req)
	if err != nil {
		return nil, gorpc.Error(codes.NotFound, err)
	}

	userId, hashedPassword, err := s.getLoginUser(ctx, req.Username)
	if err != nil {
		slogger.Error().Log("event", "get_auth_token.failed", "msg", err)
		switch {
		case errors.Is(err, ErrUserNotFound):
			return nil, gorpc.Error(codes.NotFound, err)
		default:
			return nil, gorpc.InternalError()
		}
	}

	err = s.verifyLoginPassword(hashedPassword, req.Password)
	if err != nil {
		slogger.Error().Log("event", "get_auth_token.failed", "msg", err)
		return nil, gorpc.Error(codes.NotFound, err)
	}

	jwt, err := s.generateJWT(userId, req.Username)
	if err != nil {
		slogger.Error().Log("event", "get_auth_token.failed", "msg", err)
		return nil, gorpc.InternalError()
	}

	slogger.Info().Log("event", "get_auth_token.success")
	return &pb.GenerateAuthTokenResponse{
		Token: jwt,
	}, nil
}

func (s *authService) validateLoginPayload(login *pb.GenerateAuthTokenRequest) error {
	switch {
	case login.Username == "":
		return ErrInvalidRequestCtx(`missing "username" field`)
	case login.Password == "":
		return ErrInvalidRequestCtx(`missing "password" field`)
	}
	return nil
}

func (s *authService) getLoginUser(ctx context.Context, username string) (string, string, error) {
	var userId, hashedPassword string
	err := s.db.Transact(ctx, func(tx *sqlx.Tx) error {
		query := `
		SELECT u.id, u.password
		FROM users as u
		WHERE username=$1
		`
		row := tx.QueryRowxContext(ctx, query, username)
		err := row.Scan(&userId, &hashedPassword)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrUserNotFound
		case err != nil:
			return err
		}
		return nil
	})
	if err != nil {
		return "", "", err
	}
	return userId, hashedPassword, nil
}

func (s *authService) verifyLoginPassword(hashedPassword, loginPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginPassword)); err != nil {
		return ErrIncorrectPassword
	}
	return nil
}

func (s *authService) RefreshAuthToken(ctx context.Context, req *pb.RefreshAuthTokenRequest) (*pb.RefreshAuthTokenResponse, error) {
	slogger.Info().Log("event", "refresh_auth_token.started")
	defer slogger.Info().Log("event", "refresh_auth_token.finished")

	claims, err := s.validateToken(req.Token)
	if err != nil {
		slogger.Error().Log("event", "refresh_auth_token.failed", "msg", err)
		return nil, gorpc.Error(codes.AlreadyExists, err)
	}

	jwt, err := s.generateJWT(claims.UserId, claims.Username)
	if err != nil {
		slogger.Error().Log("event", "refresh_auth_token.failed", "msg", err)
		return nil, gorpc.InternalError()
	}

	slogger.Info().Log("event", "refresh_auth_token.success")
	return &pb.RefreshAuthTokenResponse{
		Token: jwt,
	}, nil
}

func (s *authService) validateToken(token string) (*JWTClaims, error) {
	var claims JWTClaims
	jwt, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtConfig.SecretKey), nil
	})
	if err != nil || !jwt.Valid {
		return nil, ErrInvalidToken
	}

	refreshTime := time.Duration(float64(s.jwtConfig.TTL/time.Millisecond)*0.1) * time.Millisecond
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > refreshTime {
		return nil, ErrRefreshRateExceeded
	}

	return &claims, nil
}

func (s *authService) generateJWT(userId, username string) (string, error) {
	claims := &JWTClaims{
		UserId:   userId,
		Username: username,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.jwtConfig.TTL).Unix(),
			Issuer:    s.jwtConfig.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtConfig.SecretKey))
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrGeneratingToken, err)
	}

	return tokenString, nil
}
