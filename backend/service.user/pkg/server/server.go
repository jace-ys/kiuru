package server

import (
	"context"

	pb "github.com/jace-ys/kru-travel/backend/service.user/api/user"
)

type Server interface {
	Init(s pb.UserServiceServer) error
	Serve(ctx context.Context, port int) error
	Shutdown(ctx context.Context) error
}
