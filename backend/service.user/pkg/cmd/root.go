package cmd

import (
	"context"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/kru-travel/airdrop-go/pkg/authr"
	"github.com/kru-travel/airdrop-go/pkg/crdb"
	"github.com/kru-travel/airdrop-go/pkg/gorpc"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/jace-ys/kru-travel/backend/service.user/pkg/server"
	"github.com/jace-ys/kru-travel/backend/service.user/pkg/user"
)

var logger log.Logger

type config struct {
	server       server.GRPCServerConfig
	gateway      server.GatewayConfig
	database     crdb.Config
	jwtSecretKey string
}

func NewRootCmd() *cobra.Command {
	var c config

	rootCmd := &cobra.Command{
		Use:   "service",
		Short: "Start the service",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

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

			authInterceptor := gorpc.NewAuthInterceptor(
				authr.NewJWTAuthenticator(c.jwtSecretKey),
				userService.GetAuthenticatedMethods(),
			)

			grpcServer := server.NewGRPCServer(c.server, grpc.UnaryInterceptor(
				middleware.ChainUnaryServer(
					authInterceptor.AuthenticateHeader(),
				),
			))
			gatewayProxy := server.NewGatewayProxy(c.gateway, c.server.Host, c.server.Port)

			g, ctx := errgroup.WithContext(ctx)
			g.Go(func() error {
				level.Info(logger).Log("event", "grpc_server.started", "port", c.server.Port)
				defer level.Info(logger).Log("event", "grpc_server.stopped")
				return userService.StartServer(ctx, grpcServer)
			})
			g.Go(func() error {
				level.Info(logger).Log("event", "gateway_proxy.started", "port", c.gateway.Port)
				defer level.Info(logger).Log("event", "gateway_proxy.stopped")
				return userService.StartServer(ctx, gatewayProxy)
			})
			g.Go(func() error {
				select {
				case <-ctx.Done():
					grpcServer.Shutdown(ctx)
					gatewayProxy.Shutdown(ctx)
					return ctx.Err()
				}
			})

			if err := g.Wait(); err != nil {
				exit(err)
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
	rootCmd.PersistentFlags().StringVar(&c.jwtSecretKey, "jwt-secret", "", "secret key used to sign JWTs")

	return rootCmd
}

func exit(err error) {
	level.Error(logger).Log("event", "service.fatal", "msg", err)
	os.Exit(1)
}
