package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/kru-travel/airdrop-go/pkg/gorpc"
	"github.com/kru-travel/airdrop-go/pkg/slogger"

	"google.golang.org/grpc"

	gw "github.com/jace-ys/kru-travel/backend/service.user/api/user"
)

type GatewayConfig struct {
	Port int
}

type gatewayProxy struct {
	Server       *http.Server
	ProxyOptions *proxyOptions
}

type proxyOptions struct {
	Host string
	Port int
}

func NewGatewayProxy(proxyHost string, proxyPort int) *gatewayProxy {
	return &gatewayProxy{
		Server: &http.Server{
			Handler: runtime.NewServeMux(),
		},
		ProxyOptions: &proxyOptions{
			Host: proxyHost,
			Port: proxyPort,
		},
	}
}

func (g *gatewayProxy) Init(ctx context.Context, s gw.UserServiceServer) error {
	runtime.HTTPError = gorpc.GatewayHTTPError
	proxyAddr := fmt.Sprintf("%s:%d", g.ProxyOptions.Host, g.ProxyOptions.Port)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	return gw.RegisterUserServiceHandlerFromEndpoint(
		ctx,
		g.Server.Handler.(*runtime.ServeMux),
		proxyAddr,
		opts,
	)
}

func (g *gatewayProxy) Serve(port int) error {
	g.Server.Addr = fmt.Sprintf(":%d", port)
	slogger.Info().Log("event", "gateway_proxy.started", "port", port)
	defer slogger.Info().Log("event", "gateway_proxy.stopped")
	return fmt.Errorf("gateway proxy failed to serve: %w", g.Server.ListenAndServe())
}

func (g *gatewayProxy) Shutdown(ctx context.Context) error {
	return g.Server.Shutdown(ctx)
}
