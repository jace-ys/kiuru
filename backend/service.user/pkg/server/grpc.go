package server

import (
	"context"
	"fmt"
	"net"

	"github.com/kru-travel/airdrop-go/pkg/slogger"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	pb "github.com/jace-ys/kru-travel/backend/service.user/api/user"
)

type GrpcServerConfig struct {
	Host string
	Port int
}

type grpcServer struct {
	Server *grpc.Server
}

func NewGrpcServer() *grpcServer {
	return &grpcServer{
		Server: grpc.NewServer(),
	}
}

func (g *grpcServer) Init(ctx context.Context, s pb.UserServiceServer) error {
	pb.RegisterUserServiceServer(g.Server, s)
	return nil
}

func (g *grpcServer) Serve(port int) error {
	slogger.Info().Log("event", "grpc_server.started", "port", port)
	defer slogger.Info().Log("event", "grpc_server.stopped")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return errors.Wrap(err, "grpc server failed to serve")
	}
	return errors.Wrap(g.Server.Serve(lis), "grpc server failed to serve")
}

func (g *grpcServer) Shutdown(ctx context.Context) error {
	g.Server.Stop()
	return nil
}
