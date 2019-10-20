package cmd

import (
	"context"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kru-travel/airdrop-go/pkg/crdb"
	"github.com/spf13/cobra"

	"github.com/jace-ys/kru-travel/backend/service.user/pkg/server"
	"github.com/jace-ys/kru-travel/backend/service.user/pkg/user"
)

var logger log.Logger

type config struct {
	server   server.GRPCServerConfig
	gateway  server.GatewayConfig
	database crdb.Config
}

func NewRootCmd() *cobra.Command {
	var c config

	rootCmd := &cobra.Command{
		Use:   "service",
		Short: "Start the service",
		Run: func(cmd *cobra.Command, args []string) {
			logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
			logger = log.With(logger, "ts", log.DefaultTimestampUTC, "source", log.DefaultCaller)

			crdbClient := crdb.NewCRDBClient(c.database)

			userService, err := user.NewService(logger, crdbClient)
			if err != nil {
				exit(err)
			}

			if err := userService.Init(); err != nil {
				exit(err)
			}
			defer userService.Teardown()

			grpcServer := server.NewGRPCServer(c.server)
			gatewayProxy := server.NewGatewayProxy(c.gateway, c.server.Host, c.server.Port)

			errChan := make(chan error)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func(errChan chan error) {
				level.Info(logger).Log("event", "grpc_server.started", "port", c.server.Port)
				defer level.Info(logger).Log("event", "grpc_server.stopped")
				errChan <- userService.StartServer(ctx, grpcServer)
			}(errChan)

			go func(errChan chan error) {
				level.Info(logger).Log("event", "gateway_proxy.started", "port", c.gateway.Port)
				defer level.Info(logger).Log("event", "gateway_proxy.stopped")
				errChan <- userService.StartServer(ctx, gatewayProxy)
			}(errChan)

			select {
			case err := <-errChan:
				exit(err)
			case <-ctx.Done():
				exit(ctx.Err())
			}
		},
	}

	rootCmd.PersistentFlags().IntVar(&c.server.Port, "port", 8080, "port for the gRPC server")
	rootCmd.PersistentFlags().StringVar(&c.server.Host, "host", "127.0.0.1", "host for the gRPC server")
	rootCmd.PersistentFlags().IntVar(&c.gateway.Port, "gateway-port", 8081, "port for the REST gateway proxy")
	rootCmd.PersistentFlags().StringVar(&c.database.Host, "crdb-host", "127.0.0.1", "host for the CockroachDB cluster")
	rootCmd.PersistentFlags().IntVar(&c.database.Port, "crdb-port", 26257, "port for the CockroachDB cluster")
	rootCmd.PersistentFlags().StringVar(&c.database.User, "crdb-user", "", "user for the CockroachDB cluster")
	rootCmd.PersistentFlags().StringVar(&c.database.DBName, "crdb-dbname", "", "database name for the CockroachDB cluster")

	return rootCmd
}

func exit(err error) {
	level.Error(logger).Log("event", "service.fatal", "msg", err)
	os.Exit(1)
}
