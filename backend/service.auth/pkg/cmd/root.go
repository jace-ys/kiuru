package cmd

import (
	"context"
	"os"
	"time"

	"github.com/kru-travel/airdrop-go/pkg/crdb"
	"github.com/kru-travel/airdrop-go/pkg/slogger"
	"github.com/spf13/cobra"

	"github.com/jace-ys/kru-travel/backend/service.auth/pkg/auth"
	"github.com/jace-ys/kru-travel/backend/service.auth/pkg/server"
)

type config struct {
	server   server.GRPCServerConfig
	gateway  server.GatewayConfig
	database crdb.Config
	jwt      auth.JWTConfig
}

func NewRootCmd() *cobra.Command {
	var c config

	rootCmd := &cobra.Command{
		Use:   "service",
		Short: "Start the service",
		Run: func(cmd *cobra.Command, args []string) {
			crdbClient := crdb.NewCRDBClient(c.database)

			authService, err := auth.NewService(crdbClient, c.jwt)
			if err != nil {
				exit(err)
			}

			if err := authService.Init(); err != nil {
				exit(err)
			}
			defer authService.Teardown()

			grpcServer := server.NewGRPCServer()
			gatewayProxy := server.NewGatewayProxy(c.server.Host, c.server.Port)

			errChan := make(chan error)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func(errChan chan error) {
				errChan <- authService.StartServer(ctx, grpcServer, c.server.Port)
			}(errChan)

			go func(errChan chan error) {
				errChan <- authService.StartServer(ctx, gatewayProxy, c.gateway.Port)
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
	rootCmd.PersistentFlags().DurationVar(&c.database.RetryInterval, "crdb-retry-interval", 15*time.Second, "retry interval for connecting to the CockroachDB cluster")
	rootCmd.PersistentFlags().IntVar(&c.database.RetryCount, "crdb-retry-count", 10, "max number of retries for connecting to the CockroachDB cluster")
	rootCmd.PersistentFlags().BoolVar(&c.database.Insecure, "crdb-insecure", false, "enable insecure mode for the CockroachDB cluster")
	rootCmd.PersistentFlags().StringVar(&c.jwt.SecretKey, "token-secret", "", "secret key used to sign JWTs")
	rootCmd.PersistentFlags().StringVar(&c.jwt.Issuer, "token-issuer", "", "issuer of generated JWTs")
	rootCmd.PersistentFlags().DurationVar(&c.jwt.TTL, "token-ttl", 15*time.Minute, "time-to-live for generated JWTs")

	return rootCmd
}

func exit(err error) {
	slogger.Error().Log("event", "service.fatal", "msg", err)
	os.Exit(1)
}
