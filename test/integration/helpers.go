package integration

import (
	"context"
	"fmt"
	"time"

	"github.com/kiuru-travel/airdrop-go/pkg/authr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/jace-ys/kiuru/test/api/auth"
	"github.com/jace-ys/kiuru/test/api/user"
)

var (
	UserOne = user.User{
		Id:       "e8e37f43-9f11-4f0e-9bb6-f4b65cb10586",
		Username: "jaceys",
		Password: "my-secret-password",
		Name:     "Jace Tan",
		Email:    "jaceys.tan@gmail.com",
	}

	UserTwo = user.User{
		Id:       "26706953-415f-41ac-b720-8919fe8c611d",
		Username: "lowchiaying",
		Password: "my-secret-password",
		Name:     "Low Chia Ying",
		Email:    "chiaying.low@gmail.com",
	}
)

func generateToken(ttl time.Duration, userID, username string) (string, error) {
	claims := authr.NewJWTClaims("issuer", ttl, userID, username)
	return authr.GenerateToken(claims, "my-secret-key")
}

func withBearerAuthorization(ctx context.Context, token string) context.Context {
	md := metadata.New(map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	})
	return metadata.NewOutgoingContext(ctx, md)
}

func newAuthServiceClient(address string) (auth.AuthServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return auth.NewAuthServiceClient(conn), nil
}

func newUserServiceClient(address string) (user.UserServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return user.NewUserServiceClient(conn), nil
}
