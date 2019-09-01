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
		Use:   "service",
		Short: "Start the service",
		Run: func(cmd *cobra.Command, args []string) {
			userService, err := user.NewService()
			if err != nil {
				exit(err)
			}
			grpcServer := server.NewGrpcServer()
			gatewayProxy := server.NewGatewayProxy(c.server.host, c.server.port)

			errChan := make(chan error)
			ctx, cancelFunc := context.WithCancel(context.Background())
			defer cancelFunc()

			go func(errChan chan error) {
				errChan <- userService.StartServer(ctx, grpcServer, c.server.port)
			}(errChan)

			go func(errChan chan error) {
				errChan <- userService.StartServer(ctx, gatewayProxy, c.gateway.port)
			}(errChan)

			select {
			case err := <-errChan:
				exit(err)
			case <-ctx.Done():
				exit(ctx.Err())
			}
		},
	}

	rootCmd.PersistentFlags().IntVar(&c.server.port, "port", 8080, "port for the gRPC server")
	rootCmd.PersistentFlags().IntVar(&c.gateway.port, "gateway-port", 8081, "port for the REST gateway proxy")
	rootCmd.PersistentFlags().StringVar(&c.server.host, "host", "localhost", "host for the gRPC server")

	return rootCmd
}

func exit(err error) {
	slogger.Error().Log("event", "service.fatal", "msg", err)
	os.Exit(1)
}
