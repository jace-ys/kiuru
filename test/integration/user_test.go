package integration

import (
	"context"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	assert "github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jace-ys/kiuru/test/api/auth"
	"github.com/jace-ys/kiuru/test/api/user"
)

func TestUserService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	authService, err := NewAuthServiceClient("127.0.0.1:5001")
	assert.NoError(t, err)
	userService, err := NewUserServiceClient("127.0.0.1:5002")
	assert.NoError(t, err)

	t.Run("ListUsers", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			resp, err := userService.ListUsers(ctx, &user.ListUsersRequest{})

			assert.Equal(t, codes.OK.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.NotNil(t, resp, "Should return a non-nil response")

			assert.Equal(t, 2, len(resp.Users), "Should return two users")
		})
	})

	t.Run("GetUser", func(t *testing.T) {
		t.Run("NotFound", func(t *testing.T) {
			resp, err := userService.GetUser(ctx, &user.GetUserRequest{
				Id: "invalid",
			})

			assert.Equal(t, codes.NotFound.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.Nil(t, resp, "Should return a nil response")
		})

		t.Run("OK", func(t *testing.T) {
			resp, err := userService.GetUser(ctx, &user.GetUserRequest{
				Id: UserOne.Id,
			})

			assert.Equal(t, codes.OK.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.NotNil(t, resp, "Should return a non-nil response")

			assert.Equal(t, UserOne.Username, resp.User.Username, "Should return the correct username")
			assert.Equal(t, UserOne.Name, resp.User.Name, "Should return the correct name")
			assert.Equal(t, UserOne.Email, resp.User.Email, "Should return the correct email")
		})
	})

	t.Run("CreateUser", func(t *testing.T) {
		t.Run("InvalidArgument", func(t *testing.T) {
			resp, err := userService.CreateUser(ctx, &user.CreateUserRequest{
				User: &user.User{
					Username: "username",
					Name:     "name",
					Email:    "email",
				},
			})

			assert.Equal(t, codes.InvalidArgument.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.Nil(t, resp, "Should return a nil response")
		})

		t.Run("OK", func(t *testing.T) {
			resp, err := userService.CreateUser(ctx, &user.CreateUserRequest{
				User: &user.User{
					Username: "username",
					Password: "password",
					Name:     "name",
					Email:    "email",
				},
			})

			assert.Equal(t, codes.OK.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.NotNil(t, resp, "Should return a non-nil response")

			id, err := uuid.FromString(resp.Id)
			assert.NoError(t, err)
			assert.Len(t, id.Bytes(), 16, "Should return a valid UUID")
		})

		t.Run("AlreadyExists", func(t *testing.T) {
			resp, err := userService.CreateUser(ctx, &user.CreateUserRequest{
				User: &user.User{
					Username: UserOne.Username,
					Password: UserOne.Password,
					Name:     UserOne.Name,
					Email:    UserOne.Email,
				},
			})

			assert.Equal(t, codes.AlreadyExists.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.Nil(t, resp, "Should return a nil response")
		})
	})

	t.Run("DeleteUser", func(t *testing.T) {
		t.Run("Unauthenticated", func(t *testing.T) {
			ctx = WithBearerAuthorization(ctx, "")
			resp, err := userService.DeleteUser(ctx, &user.DeleteUserRequest{
				Id: UserOne.Id,
			})

			assert.Equal(t, codes.Unauthenticated.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.Nil(t, resp, "Should return a nil response")
		})

		t.Run("PermissionDenied", func(t *testing.T) {
			token, err := GenerateToken(time.Minute, UserOne.Id, UserOne.Username)
			assert.NoError(t, err)

			ctx = WithBearerAuthorization(ctx, token)
			resp, err := userService.DeleteUser(ctx, &user.DeleteUserRequest{
				Id: UserTwo.Id,
			})

			assert.Equal(t, codes.PermissionDenied.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.Nil(t, resp, "Should return a nil response")
		})

		t.Run("NotFound", func(t *testing.T) {
			token, err := GenerateToken(time.Minute, "invalid", UserOne.Username)
			assert.NoError(t, err)

			ctx = WithBearerAuthorization(ctx, token)
			resp, err := userService.DeleteUser(ctx, &user.DeleteUserRequest{
				Id: "invalid",
			})

			assert.Equal(t, codes.NotFound.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.Nil(t, resp, "Should return a nil response")
		})

		t.Run("OK", func(t *testing.T) {
			token, err := GenerateToken(time.Minute, UserOne.Id, UserOne.Username)
			assert.NoError(t, err)

			ctx = WithBearerAuthorization(ctx, token)
			resp, err := userService.DeleteUser(ctx, &user.DeleteUserRequest{
				Id: UserOne.Id,
			})

			assert.Equal(t, codes.OK.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.NotNil(t, resp, "Should return a non-nil response")
		})

		t.Run("Unauthenticated/Revoked", func(t *testing.T) {
			token, err := GenerateToken(time.Minute, UserOne.Id, UserOne.Username)
			assert.NoError(t, err)

			_, err = authService.RevokeToken(ctx, &auth.RevokeTokenRequest{
				Token: token,
			})
			assert.NoError(t, err)

			ctx = WithBearerAuthorization(ctx, token)
			resp, err := userService.DeleteUser(ctx, &user.DeleteUserRequest{
				Id: UserOne.Id,
			})

			assert.Equal(t, codes.Unauthenticated.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.Nil(t, resp, "Should return a nil response")
		})
	})
}
