package integration

import (
	"context"
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jace-ys/kiuru/test/api/auth"
)

func TestAuthService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	authService, err := NewAuthServiceClient("127.0.0.1:5001")
	assert.NoError(t, err)

	t.Run("GenerateAuthToken", func(t *testing.T) {
		t.Run("InvalidArgument", func(t *testing.T) {
			resp, err := authService.GenerateAuthToken(ctx, &auth.GenerateAuthTokenRequest{})

			assert.Equal(t, codes.InvalidArgument.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.Nil(t, resp, "Should return a nil response")
		})

		t.Run("NotFound", func(t *testing.T) {
			resp, err := authService.GenerateAuthToken(ctx, &auth.GenerateAuthTokenRequest{
				Username: "invalid",
				Password: "password",
			})

			assert.Equal(t, codes.NotFound.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.Nil(t, resp, "Should return a nil response")
		})

		t.Run("Unauthenticated", func(t *testing.T) {
			resp, err := authService.GenerateAuthToken(ctx, &auth.GenerateAuthTokenRequest{
				Username: UserOne.Username,
				Password: "password",
			})

			assert.Equal(t, codes.Unauthenticated.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.Nil(t, resp, "Should return a nil response")
		})

		t.Run("OK", func(t *testing.T) {
			resp, err := authService.GenerateAuthToken(ctx, &auth.GenerateAuthTokenRequest{
				Username: UserOne.Username,
				Password: UserOne.Password,
			})

			assert.Equal(t, codes.OK.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.NotNil(t, resp.Token, "Should return an authentication token")
		})
	})

	t.Run("RefreshAuthToken", func(t *testing.T) {
		t.Run("InvalidArgument", func(t *testing.T) {
			resp, err := authService.RefreshAuthToken(ctx, &auth.RefreshAuthTokenRequest{
				Token: "",
			})

			assert.Equal(t, codes.InvalidArgument.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.Nil(t, resp, "Should return a nil response")
		})

		t.Run("ResourceExhausted", func(t *testing.T) {
			token, err := GenerateToken(time.Minute, "userID", "username")
			assert.NoError(t, err)

			resp, err := authService.RefreshAuthToken(ctx, &auth.RefreshAuthTokenRequest{
				Token: token,
			})

			assert.Equal(t, codes.ResourceExhausted.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.Nil(t, resp, "Should return a nil response")
		})

		t.Run("OK", func(t *testing.T) {
			token, err := GenerateToken(time.Millisecond, "userID", "username")
			assert.NoError(t, err)

			resp, err := authService.RefreshAuthToken(ctx, &auth.RefreshAuthTokenRequest{
				Token: token,
			})

			assert.Equal(t, codes.OK.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.NotNil(t, resp, "Should return an authentication token")
		})
	})

	t.Run("RevokeAuthToken", func(t *testing.T) {
		t.Run("InvalidArgument", func(t *testing.T) {
			resp, err := authService.RevokeAuthToken(ctx, &auth.RevokeAuthTokenRequest{
				Token: "",
			})

			assert.Equal(t, codes.InvalidArgument.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.Nil(t, resp, "Should return a nil response")
		})

		t.Run("OK", func(t *testing.T) {
			token, err := GenerateToken(time.Minute, "userID", "username")
			assert.NoError(t, err)

			resp, err := authService.RevokeAuthToken(ctx, &auth.RevokeAuthTokenRequest{
				Token: token,
			})

			assert.Equal(t, codes.OK.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.NotNil(t, resp, "Should return a non-nil response")
		})

		t.Run("RefreshAuthToken/Unauthenticated", func(t *testing.T) {
			token, err := GenerateToken(time.Minute, "userID", "username")
			assert.NoError(t, err)

			resp, err := authService.RefreshAuthToken(ctx, &auth.RefreshAuthTokenRequest{

				Token: token,
			})

			assert.Equal(t, codes.Unauthenticated.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.Nil(t, resp, "Should return a nil response")
		})
	})
}
