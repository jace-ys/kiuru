package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/kiuru-travel/airdrop-go/gorpc"
	"google.golang.org/grpc"

	gw "github.com/jace-ys/kiuru/backend/service.auth/api/auth"
)

type GatewayProxyConfig struct {
	Port     int
	Endpoint string
}

type GatewayProxy struct {
	config *GatewayProxyConfig
	server *http.Server
}

func NewGatewayProxy(port int, endpoint string) *GatewayProxy {
	return &GatewayProxy{
		config: &GatewayProxyConfig{
			Port:     port,
			Endpoint: endpoint,
		},
		server: &http.Server{
			Handler: runtime.NewServeMux(
				runtime.WithProtoErrorHandler(gorpc.HTTPError),
				runtime.WithForwardResponseOption(gorpc.ForwardReponseHTTPStatus),
				runtime.WithOutgoingHeaderMatcher(gorpc.OutgoingHeaderMatcher),
				runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: false}),
			),
			Addr: fmt.Sprintf(":%d", port),
		},
	}
}

func (g *GatewayProxy) Init(ctx context.Context, s gw.AuthServiceServer) error {
	err := gw.RegisterAuthServiceHandlerFromEndpoint(
		ctx,
		g.server.Handler.(*runtime.ServeMux),
		g.config.Endpoint,
		[]grpc.DialOption{grpc.WithInsecure()},
	)
	if err != nil {
		return fmt.Errorf("gateway proxy failed to initialize: %w", err)
	}
	return nil
}

func (g *GatewayProxy) Serve() error {
	if err := g.server.ListenAndServe(); err != nil {
		return fmt.Errorf("gateway proxy failed to serve: %w", err)
	}
	return nil
}

func (g *GatewayProxy) Shutdown(ctx context.Context) error {
	if err := g.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("gateway proxy failed to shutdown: %w", err)
	}
	return nil
}
