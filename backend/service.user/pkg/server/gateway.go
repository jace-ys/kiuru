package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/kru-travel/airdrop-go/pkg/gorpc"

	"google.golang.org/grpc"

	gw "github.com/jace-ys/kru-travel/backend/service.user/api/user"
)

type GatewayProxyConfig struct {
	Port     int
	Endpoint string
}

type gatewayProxy struct {
	config *GatewayProxyConfig
	server *http.Server
}

func NewGatewayProxy(port int, endpoint string) *gatewayProxy {
	return &gatewayProxy{
		config: &GatewayProxyConfig{
			Port:     port,
			Endpoint: endpoint,
		},
		server: &http.Server{
			Handler: runtime.NewServeMux(),
		},
	}
}

func (g *gatewayProxy) Init(ctx context.Context, s gw.UserServiceServer) error {
	runtime.HTTPError = gorpc.HTTPError
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterUserServiceHandlerFromEndpoint(
		ctx,
		g.server.Handler.(*runtime.ServeMux),
		g.config.Endpoint,
		opts,
	)
	if err != nil {
		return fmt.Errorf("gateway proxy failed to initialize: %w", err)
	}
	return nil
}

func (g *gatewayProxy) Serve() error {
	g.server.Addr = fmt.Sprintf(":%d", g.config.Port)
	if err := g.server.ListenAndServe(); err != nil {
		return fmt.Errorf("gateway proxy failed to serve: %w", err)
	}
	return nil
}

func (g *gatewayProxy) Shutdown(ctx context.Context) error {
	if err := g.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("gateway proxy failed to shutdown: %w", err)
	}
	return nil
}
