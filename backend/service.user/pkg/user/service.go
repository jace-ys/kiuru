package user

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"github.com/kru-travel/airdrop-go/pkg/slogger"

	pb "github.com/jace-ys/kru-travel/backend/service.user/api/user"
)

type userService struct {
	grpcServer *grpc.Server
}

func NewService() *userService {
	opts := []grpc.ServerOption{}
	return &userService{grpc.NewServer(opts...)}
}

func (s *userService) Init() error {
	pb.RegisterUserServiceServer(s.grpcServer, s)
	return nil
}

func (s *userService) Serve(lis net.Listener) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, ":8080", opts)
	if err != nil {
		return err
	}

	go http.ListenAndServe(":8081", mux)

	return s.grpcServer.Serve(lis)
}

func (s *userService) Shutdown() {
	s.grpcServer.Stop()
}

func (s *userService) GetUser(ctx context.Context, r *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	slogger.Info().Log("event", "request.get_user", "user_id", r.Id)
	return &pb.GetUserResponse{
		Id: r.Id,
	}, nil
}
