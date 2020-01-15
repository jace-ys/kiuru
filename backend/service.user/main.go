package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/kru-travel/airdrop-go/pkg/authr"
	"github.com/kru-travel/airdrop-go/pkg/crdb"
	"github.com/kru-travel/airdrop-go/pkg/gorpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/jace-ys/kru-travel/backend/service.user/pkg/server"
	"github.com/jace-ys/kru-travel/backend/service.user/pkg/user"
)

var logger log.Logger

func main() {
	c := parseCommand()

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
	gatewayProxy := server.NewGatewayProxy(c.gateway)

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
}

type config struct {
	server       server.GRPCServerConfig
	gateway      server.GatewayConfig
	database     crdb.Config
	jwtSecretKey string
}

func parseCommand() *config {
	c := config{}
	kingpin.Flag("port", "port for the gRPC server").Default("8080").IntVar(&c.server.Port)
	kingpin.Flag("host", "host for the gRPC server").Default("127.0.0.1").StringVar(&c.server.Host)
	kingpin.Flag("gateway-port", "port for the REST gateway proxy").Default("8081").IntVar(&c.gateway.Port)
	kingpin.Flag("crdb-host", "host for connecting to CockroachDB").Default("127.0.0.1").StringVar(&c.database.Host)
	kingpin.Flag("crdb-port", "port for connecting to CockroachDB").Default("26257").IntVar(&c.database.Port)
	kingpin.Flag("crdb-user", "user for connecting to CockroachDB").StringVar(&c.database.User)
	kingpin.Flag("crdb-dbname", "database name for connecting to CockroachDB").StringVar(&c.database.DBName)
	kingpin.Flag("jwt-secret", "secret key used to sign JWTs").StringVar(&c.jwtSecretKey)
	kingpin.Parse()

	c.gateway.Endpoint = fmt.Sprintf("%s:%d", c.server.Host, c.server.Port)
	return &c
}

func exit(err error) {
	level.Error(logger).Log("event", "service.fatal", "msg", err)
	os.Exit(1)
}
