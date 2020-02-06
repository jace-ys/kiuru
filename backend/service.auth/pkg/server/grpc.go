package server

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	pb "github.com/jace-ys/kiuru/backend/service.auth/api/auth"
)

type GRPCServerConfig struct {
	Port int
}

type grpcServer struct {
	config *GRPCServerConfig
	server *grpc.Server
}

func NewGRPCServer(port int, opt ...grpc.ServerOption) *grpcServer {
	return &grpcServer{
		config: &GRPCServerConfig{
			Port: port,
		},
		server: grpc.NewServer(opt...),
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
