package integration

import (
	"context"
	"testing"

	assert "github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	user "github.com/jace-ys/kru-travel/test/api/user"
)

func TestUserService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := NewUserServiceClient("127.0.0.1:5002")
	assert.NoError(t, err)

	t.Run("GET /users", func(t *testing.T) {
		req := &user.GetAllUsersRequest{}
		resp, err := client.GetAllUsers(ctx, req)
		code := status.Code(err)

		assert.Equal(t, codes.OK, code, "Should return codes.OK")
		assert.Equal(t, 1, len(resp.Users), "Should return a single user")
		assert.Equal(t, "jaceys", resp.Users[0].Username, "Should return the correct username")
		assert.Equal(t, "jaceys.tan@gmail.com", resp.Users[0].Email, "Should return the correct email")
		assert.Equal(t, "Jace Tan", resp.Users[0].Name, "Should return the correct name")
	})

	t.Run("GET /users/{userID}", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			req := &user.GetUserRequest{
				Id: "e8e37f43-9f11-4f0e-9bb6-f4b65cb10586",
			}
			resp, err := client.GetUser(ctx, req)
			code := status.Code(err)

			assert.Equal(t, codes.OK, code, "Should return codes.OK")
			assert.Equal(t, "jaceys", resp.User.Username, "Should return the correct username")
			assert.Equal(t, "jaceys.tan@gmail.com", resp.User.Email, "Should return the correct email")
			assert.Equal(t, "Jace Tan", resp.User.Name, "Should return the correct name")
		})

		t.Run("NotFound", func(t *testing.T) {
			req := &user.GetUserRequest{
				Id: "invalid",
			}
			resp, err := client.GetUser(ctx, req)
			code := status.Code(err)

			assert.Equal(t, codes.NotFound, code, "Should return codes.NotFound")
			assert.Nil(t, resp, "Should return a nil response")
		})
	})
}

func NewUserServiceClient(address string) (user.UserServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return user.NewUserServiceClient(conn), nil
}
