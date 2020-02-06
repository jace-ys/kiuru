package integration

import (
	"context"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	assert "github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jace-ys/kiuru/test/api/user"
)

func TestUserService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service, err := newUserServiceClient("127.0.0.1:5002")
	assert.NoError(t, err)

	t.Run("GetAllUsers", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			req := &user.GetAllUsersRequest{}

			resp, err := service.GetAllUsers(ctx, req)
			assert.Equal(t, codes.OK.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.NotNil(t, resp, "Should return a non-nil response")

			assert.Equal(t, 2, len(resp.Users), "Should return two users")
		})
	})

	t.Run("GetUser", func(t *testing.T) {
		t.Run("NotFound", func(t *testing.T) {
			req := &user.GetUserRequest{
				Id: "invalid",
			}

			resp, err := service.GetUser(ctx, req)
			assert.Equal(t, codes.NotFound.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.Nil(t, resp, "Should return a nil response")
		})

		t.Run("OK", func(t *testing.T) {
			req := &user.GetUserRequest{
				Id: UserOne.Id,
			}

			resp, err := service.GetUser(ctx, req)
			assert.Equal(t, codes.OK.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.NotNil(t, resp, "Should return a non-nil response")

			assert.Equal(t, UserOne.Username, resp.User.Username, "Should return the correct username")
			assert.Equal(t, UserOne.Name, resp.User.Name, "Should return the correct name")
			assert.Equal(t, UserOne.Email, resp.User.Email, "Should return the correct email")
		})
	})

	t.Run("CreateUser", func(t *testing.T) {
		t.Run("InvalidArgument", func(t *testing.T) {
			req := &user.CreateUserRequest{
				User: &user.User{
					Username: "username",
					Name:     "name",
					Email:    "email",
				},
			}

			resp, err := service.CreateUser(ctx, req)
			assert.Equal(t, codes.InvalidArgument.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.Nil(t, resp, "Should return a nil response")
		})

		t.Run("OK", func(t *testing.T) {
			req := &user.CreateUserRequest{
				User: &user.User{
					Username: "username",
					Password: "password",
					Name:     "name",
					Email:    "email",
				},
			}

			resp, err := service.CreateUser(ctx, req)
			assert.Equal(t, codes.OK.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.NotNil(t, resp, "Should return a non-nil response")

			id, err := uuid.FromString(resp.Id)
			assert.NoError(t, err)
			assert.Len(t, id.Bytes(), 16, "Should return a valid UUID")
		})

		t.Run("AlreadyExists", func(t *testing.T) {
			req := &user.CreateUserRequest{
				User: &user.User{
					Username: UserOne.Username,
					Password: UserOne.Password,
					Name:     UserOne.Name,
					Email:    UserOne.Email,
				},
			}

			resp, err := service.CreateUser(ctx, req)
			assert.Equal(t, codes.AlreadyExists.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.Nil(t, resp, "Should return a nil response")
		})
	})

	t.Run("DeleteUser", func(t *testing.T) {
		t.Run("Unauthenticated", func(t *testing.T) {
			ctx = withBearerAuthorization(ctx, "")
			req := &user.DeleteUserRequest{
				Id: UserOne.Id,
			}

			resp, err := service.DeleteUser(ctx, req)
			assert.Equal(t, codes.Unauthenticated.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.Nil(t, resp, "Should return a nil response")
		})

		t.Run("PermissionDenied", func(t *testing.T) {
			token, err := generateToken(time.Minute, UserOne.Id, UserOne.Username)
			assert.NoError(t, err)

			ctx = withBearerAuthorization(ctx, token)
			req := &user.DeleteUserRequest{
				Id: UserTwo.Id,
			}

			resp, err := service.DeleteUser(ctx, req)
			assert.Equal(t, codes.PermissionDenied.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.Nil(t, resp, "Should return a nil response")
		})

		t.Run("OK", func(t *testing.T) {
			token, err := generateToken(time.Minute, UserOne.Id, UserOne.Username)
			assert.NoError(t, err)

			ctx = withBearerAuthorization(ctx, token)
			req := &user.DeleteUserRequest{
				Id: UserOne.Id,
			}

			resp, err := service.DeleteUser(ctx, req)
			assert.Equal(t, codes.OK.String(), status.Code(err).String(), status.Convert(err).Message())
			assert.NotNil(t, resp, "Should return a non-nil response")
		})
	})
}
