package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/kru-travel/airdrop-go/pkg/gorpc"

	"google.golang.org/grpc"

	gw "github.com/jace-ys/kru-travel/backend/service.auth/api/auth"
)

type GatewayConfig struct {
	Port int
}

type gatewayProxy struct {
	server      *http.Server
	config      *GatewayConfig
	proxyConfig *proxyConfig
}

type proxyConfig struct {
	Host string
	Port int
}

func NewGatewayProxy(config GatewayConfig, proxyHost string, proxyPort int) *gatewayProxy {
	return &gatewayProxy{
		config: &config,
		server: &http.Server{
			Handler: runtime.NewServeMux(),
		},
		proxyConfig: &proxyConfig{
			Host: proxyHost,
			Port: proxyPort,
		},
	}
}

func (g *gatewayProxy) Init(ctx context.Context, s gw.AuthServiceServer) error {
	runtime.HTTPError = gorpc.GatewayHTTPError
	proxyAddr := fmt.Sprintf("%s:%d", g.proxyConfig.Host, g.proxyConfig.Port)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterAuthServiceHandlerFromEndpoint(
		ctx,
		g.server.Handler.(*runtime.ServeMux),
		proxyAddr,
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
