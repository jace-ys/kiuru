package server

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	pb "github.com/jace-ys/kru-travel/backend/service.auth/api/auth"
)

type GRPCServerConfig struct {
	Host string
	Port int
}

type grpcServer struct {
	server *grpc.Server
	config *GRPCServerConfig
}

func NewGRPCServer(config GRPCServerConfig) *grpcServer {
	return &grpcServer{
		config: &config,
		server: grpc.NewServer(),
	}
}

func (g *grpcServer) Init(ctx context.Context, s pb.AuthServiceServer) error {
	pb.RegisterAuthServiceServer(g.server, s)
	return nil
}

func (g *grpcServer) Serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.config.Port))
	if err != nil {
		return fmt.Errorf("grpc server failed to serve: %w", err)
	}
	if err := g.server.Serve(lis); err != nil {
		return fmt.Errorf("grpc server failed to serve: %w", err)
	}
	return nil
}

func (g *grpcServer) Shutdown(ctx context.Context) error {
	g.server.GracefulStop()
	return nil
}
