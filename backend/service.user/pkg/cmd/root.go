package cmd

import (
	"context"
	"os"

	"github.com/kru-travel/airdrop-go/pkg/slogger"
	"github.com/spf13/cobra"

	"github.com/jace-ys/kru-travel/backend/service.user/pkg/server"
	"github.com/jace-ys/kru-travel/backend/service.user/pkg/user"
)

type config struct {
	server  serverConfig
	gateway gatewayConfig
}

type serverConfig struct {
	host string
	port int
}

type gatewayConfig struct {
	port int
}

func NewCmd() *cobra.Command {
	var c config

	rootCmd := &cobra.Command{
		Use:   "server",
		Short: "Start the server for the service",
		Run: func(cmd *cobra.Command, args []string) {
			userService := user.NewService()
			grpcServer := server.NewGrpcServer()
			gatewayProxy := server.NewGatewayProxy(c.server.host, c.server.port)

			errChan := make(chan error)
			ctx, cancelFunc := context.WithCancel(context.Background())
			defer cancelFunc()

			go func(errChan chan error) {
				errChan <- userService.ListenAndServe(ctx, grpcServer, c.server.port)
			}(errChan)

			go func(errChan chan error) {
				errChan <- userService.ListenAndServe(ctx, gatewayProxy, c.gateway.port)
			}(errChan)

			select {
			case err := <-errChan:
				slogger.Error().Log("event", "service.error", "msg", err.Error())
				os.Exit(1)
			case <-ctx.Done():
			}
		},
	}

	rootCmd.PersistentFlags().IntVar(&c.server.port, "port", 8080, "port for the gRPC server")
	rootCmd.PersistentFlags().IntVar(&c.gateway.port, "gateway-port", 8081, "port for the REST gateway proxy")
	rootCmd.PersistentFlags().StringVar(&c.server.host, "host", "localhost", "host for the gRPC server")

	return rootCmd
}
