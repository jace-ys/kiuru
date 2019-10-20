package cmd

import (
	"context"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kru-travel/airdrop-go/pkg/crdb"
	"github.com/kru-travel/airdrop-go/pkg/redis"
	"github.com/spf13/cobra"

	"github.com/jace-ys/kru-travel/backend/service.auth/pkg/auth"
	"github.com/jace-ys/kru-travel/backend/service.auth/pkg/server"
)

var logger log.Logger

type config struct {
	server   server.GRPCServerConfig
	gateway  server.GatewayConfig
	database crdb.Config
	redis    redis.Config
	jwt      auth.JWTConfig
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
			redisClient := redis.NewRedisClient(c.redis)

			authService, err := auth.NewService(logger, crdbClient, redisClient, c.jwt)
			if err != nil {
				exit(err)
			}

			if err := authService.Init(); err != nil {
				exit(err)
			}
			defer authService.Teardown()

			grpcServer := server.NewGRPCServer(c.server)
			gatewayProxy := server.NewGatewayProxy(c.gateway, c.server.Host, c.server.Port)

			errChan := make(chan error)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func(errChan chan error) {
				level.Info(logger).Log("event", "grpc_server.started", "port", c.server.Port)
				defer level.Info(logger).Log("event", "grpc_server.stopped")
				errChan <- authService.StartServer(ctx, grpcServer)
			}(errChan)

			go func(errChan chan error) {
				level.Info(logger).Log("event", "gateway_proxy.started", "port", c.gateway.Port)
				defer level.Info(logger).Log("event", "gateway_proxy.stopped")
				errChan <- authService.StartServer(ctx, gatewayProxy)
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
	rootCmd.PersistentFlags().StringVar(&c.database.Host, "crdb-host", "127.0.0.1", "host for connecting to CockroachDB")
	rootCmd.PersistentFlags().IntVar(&c.database.Port, "crdb-port", 26257, "port for connecting to CockroachDB")
	rootCmd.PersistentFlags().StringVar(&c.database.User, "crdb-user", "", "user for connecting to CockroachDB")
	rootCmd.PersistentFlags().StringVar(&c.database.DBName, "crdb-dbname", "", "database name for connecting to CockroachDB")
	rootCmd.PersistentFlags().StringVar(&c.redis.Host, "redis-host", "127.0.0.1", "host for connecting Redis")
	rootCmd.PersistentFlags().IntVar(&c.redis.Port, "redis-port", 6379, "port for connecting to Redis")
	rootCmd.PersistentFlags().StringVar(&c.jwt.SecretKey, "token-secret", "", "secret key used to sign JWTs")
	rootCmd.PersistentFlags().StringVar(&c.jwt.Issuer, "token-issuer", "", "issuer of generated JWTs")
	rootCmd.PersistentFlags().DurationVar(&c.jwt.TTL, "token-ttl", 15*time.Minute, "time-to-live for generated JWTs")

	return rootCmd
}

func exit(err error) {
	level.Error(logger).Log("event", "service.fatal", "msg", err)
	os.Exit(1)
}
