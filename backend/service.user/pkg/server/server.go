package server

import (
	"context"

	pb "github.com/jace-ys/kru-travel/backend/service.user/api/user"
)

type Server interface {
	Init(ctx context.Context, s pb.UserServiceServer) error
	Serve(port int) error
	Shutdown(ctx context.Context) error
}
