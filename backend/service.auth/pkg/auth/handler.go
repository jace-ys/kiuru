package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log/level"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"github.com/kru-travel/airdrop-go/pkg/gorpc"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"

	pb "github.com/jace-ys/kru-travel/backend/service.auth/api/auth"
)

func (s *authService) GenerateAuthToken(ctx context.Context, req *pb.GenerateAuthTokenRequest) (*pb.GenerateAuthTokenResponse, error) {
	level.Info(s.logger).Log("event", "get_auth_token.started")
	defer level.Info(s.logger).Log("event", "get_auth_token.finished")

	err := s.validateLoginPayload(req)
	if err != nil {
		level.Error(s.logger).Log("event", "get_auth_token.failed", "msg", err)
		return nil, gorpc.Error(err)
	}

	userId, hashedPassword, err := s.getLoginUser(ctx, req.Username)
	if err != nil {
		level.Error(s.logger).Log("event", "get_auth_token.failed", "msg", err)
		return nil, gorpc.Error(err)
	}

	err = s.verifyLoginPassword(hashedPassword, req.Password)
	if err != nil {
		level.Error(s.logger).Log("event", "get_auth_token.failed", "msg", err)
		return nil, gorpc.Error(err)
	}

	jwt, err := s.generateJWT(userId, req.Username)
	if err != nil {
		level.Error(s.logger).Log("event", "get_auth_token.failed", "msg", err)
		return nil, gorpc.Error(err)
	}

	level.Info(s.logger).Log("event", "get_auth_token.success")
	return &pb.GenerateAuthTokenResponse{
		Token: jwt,
	}, nil
}

func (s *authService) validateLoginPayload(login *pb.GenerateAuthTokenRequest) error {
	switch {
	case login.Username == "":
		return gorpc.NewErr(codes.InvalidArgument, fmt.Errorf("%w: %s", ErrInvalidRequest, `missing "username" field`))
	case login.Password == "":
		return gorpc.NewErr(codes.InvalidArgument, fmt.Errorf("%w: %s", ErrInvalidRequest, `missing "password" field`))
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
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", "", gorpc.NewErr(codes.NotFound, ErrUserNotFound)
		default:
			return "", "", gorpc.NewErr(codes.Internal, err)
		}
	}
	return userId, hashedPassword, nil
}

func (s *authService) verifyLoginPassword(hashedPassword, loginPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginPassword)); err != nil {
		return gorpc.NewErr(codes.InvalidArgument, ErrIncorrectPassword)
	}
	return nil
}

func (s *authService) RefreshAuthToken(ctx context.Context, req *pb.RefreshAuthTokenRequest) (*pb.RefreshAuthTokenResponse, error) {
	level.Info(s.logger).Log("event", "refresh_auth_token.started")
	defer level.Info(s.logger).Log("event", "refresh_auth_token.finished")

	claims, err := s.validateToken(req.Token)
	if err != nil {
		level.Error(s.logger).Log("event", "refresh_auth_token.failed", "msg", err)
		return nil, gorpc.Error(err)
	}

	err = s.isRevoked(ctx, req.Token)
	if err != nil {
		level.Error(s.logger).Log("event", "refresh_auth_token.failed", "msg", err)
		return nil, gorpc.Error(err)
	}

	err = s.isRefreshable(claims)
	if err != nil {
		level.Error(s.logger).Log("event", "refresh_auth_token.failed", "msg", err)
		return nil, gorpc.Error(err)
	}

	jwt, err := s.generateJWT(claims.UserId, claims.Username)
	if err != nil {
		level.Error(s.logger).Log("event", "refresh_auth_token.failed", "msg", err)
		return nil, gorpc.Error(err)
	}

	level.Info(s.logger).Log("event", "refresh_auth_token.success")
	return &pb.RefreshAuthTokenResponse{
		Token: jwt,
	}, nil
}

func (s *authService) isRevoked(ctx context.Context, token string) error {
	err := s.redis.Transact(ctx, func(conn redis.Conn) error {
		reply, err := conn.Do("GET", token)
		if err != nil {
			return err
		}
		if reply != nil {
			return ErrTokenRevoked
		}
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, ErrTokenRevoked):
			return gorpc.NewErr(codes.InvalidArgument, err)
		default:
			return gorpc.NewErr(codes.Internal, err)
		}
	}
	return nil
}

func (s *authService) isRefreshable(claims *JWTClaims) error {
	refreshTime := time.Duration(float64(s.jwtConfig.TTL/time.Millisecond)*0.1) * time.Millisecond
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > refreshTime {
		return gorpc.NewErr(codes.AlreadyExists, ErrRefreshRateExceeded)
	}
	return nil
}

func (s *authService) RevokeAuthToken(ctx context.Context, req *pb.RevokeAuthTokenRequest) (*pb.RevokeAuthTokenResponse, error) {
	level.Info(s.logger).Log("event", "revoke_auth_token.started")
	defer level.Info(s.logger).Log("event", "revoke_auth_token.finished")

	_, err := s.validateToken(req.Token)
	if err != nil {
		level.Error(s.logger).Log("event", "revoke_auth_token.failed", "msg", err)
		return nil, gorpc.Error(err)
	}

	err = s.cacheRevokedToken(ctx, req.Token)
	if err != nil {
		level.Error(s.logger).Log("event", "revoke_auth_token.failed", "msg", err)
		return nil, gorpc.Error(err)
	}

	level.Info(s.logger).Log("event", "revoke_auth_token.success")
	return &pb.RevokeAuthTokenResponse{}, nil
}

func (s *authService) cacheRevokedToken(ctx context.Context, token string) error {
	expiryInSeconds := strconv.Itoa(int(s.jwtConfig.TTL / time.Second))
	err := s.redis.Transact(ctx, func(conn redis.Conn) error {
		_, err := conn.Do("SET", token, "revoked", "EX", expiryInSeconds)
		return err
	})
	if err != nil {
		return gorpc.NewErr(codes.Internal, err)
	}
	return nil
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
		return "", gorpc.NewErr(codes.Internal, fmt.Errorf("%w: %s", ErrGeneratingToken, err))
	}

	return tokenString, nil
}

func (s *authService) validateToken(token string) (*JWTClaims, error) {
	var claims JWTClaims
	jwt, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtConfig.SecretKey), nil
	})
	if err != nil || !jwt.Valid {
		return nil, gorpc.NewErr(codes.InvalidArgument, ErrInvalidToken)
	}
	return &claims, nil
}
