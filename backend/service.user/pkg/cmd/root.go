package cmd

import (
	"context"
	"os"

	"github.com/kru-travel/airdrop-go/pkg/slogger"
	"github.com/spf13/cobra"

	"github.com/jace-ys/kru-travel/backend/service.user/pkg/crdb"
	"github.com/jace-ys/kru-travel/backend/service.user/pkg/server"
	"github.com/jace-ys/kru-travel/backend/service.user/pkg/user"
)

type config struct {
	server   server.GrpcServerConfig
	gateway  server.GatewayConfig
	database crdb.Config
}

func NewCmd() *cobra.Command {
	var c config

	rootCmd := &cobra.Command{
		Use:   "service",
		Short: "Start the service",
		Run: func(cmd *cobra.Command, args []string) {
			crdbClient, err := crdb.NewCrdbClient(c.database)
			if err != nil {
				exit(err)
			}

			userService := user.NewService()
			if err := userService.Init(crdbClient); err != nil {
				exit(err)
			}
			defer userService.Teardown()

			grpcServer := server.NewGrpcServer()
			gatewayProxy := server.NewGatewayProxy(c.server.Host, c.server.Port)

			errChan := make(chan error)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func(errChan chan error) {
				errChan <- userService.StartServer(ctx, grpcServer, c.server.Port)
			}(errChan)

			go func(errChan chan error) {
				errChan <- userService.StartServer(ctx, gatewayProxy, c.gateway.Port)
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
	rootCmd.PersistentFlags().StringVar(&c.server.Host, "host", "0.0.0.0", "host for the gRPC server")
	rootCmd.PersistentFlags().IntVar(&c.gateway.Port, "gateway-port", 8081, "port for the REST gateway proxy")
	rootCmd.PersistentFlags().StringVar(&c.database.Host, "crdb-host", "0.0.0.0", "host for the CockroachDB cluster")
	rootCmd.PersistentFlags().IntVar(&c.database.Port, "crdb-port", 26257, "port for the CockroachDB cluster")
	rootCmd.PersistentFlags().StringVar(&c.database.User, "crdb-user", "", "user for the CockroachDB cluster")
	rootCmd.PersistentFlags().StringVar(&c.database.DbName, "crdb-dbname", "", "database name for the CockroachDB cluster")
	rootCmd.PersistentFlags().IntVar(&c.database.Retry, "crdb-retry", 10, "retry interval for connecting to the CockroachDB cluster")
	rootCmd.PersistentFlags().BoolVar(&c.database.Insecure, "crdb-insecure", false, "enable insecure mode for the CockroachDB cluster")

	return rootCmd
}

func exit(err error) {
	slogger.Error().Log("event", "service.fatal", "msg", err)
	os.Exit(1)
}
