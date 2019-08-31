package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/kru-travel/airdrop-go/pkg/slogger"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	gw "github.com/jace-ys/kru-travel/backend/service.user/api/user"
)

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

func (g *gatewayProxy) Init(s gw.UserServiceServer) error {
	return nil
}

func (g *gatewayProxy) Serve(ctx context.Context, port int) error {
	proxyAddr := fmt.Sprintf("%s:%d", g.ProxyOptions.Host, g.ProxyOptions.Port)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterUserServiceHandlerFromEndpoint(
		ctx,
		g.Server.Handler.(*runtime.ServeMux),
		proxyAddr,
		opts,
	)
	if err != nil {
		return errors.Wrap(err, "gateway proxy failed to serve")
	}

	g.Server.Addr = fmt.Sprintf(":%d", port)
	slogger.Info().Log("event", "gateway_proxy.started", "port", port)
	defer slogger.Info().Log("event", "gateway_proxy.stopped")
	return errors.Wrap(g.Server.ListenAndServe(), "gateway proxy failed to serve")
}

func (g *gatewayProxy) Shutdown(ctx context.Context) error {
	g.Server.Shutdown(ctx)
	return nil
}
