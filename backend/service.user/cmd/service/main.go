package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kiuru-travel/airdrop-go/pkg/authr"
	"github.com/kiuru-travel/airdrop-go/pkg/cache"
	"github.com/kiuru-travel/airdrop-go/pkg/crdb"
	"github.com/kiuru-travel/airdrop-go/pkg/gorpc"
	"github.com/kiuru-travel/airdrop-go/pkg/redis"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/jace-ys/kiuru/backend/service.user/pkg/server"
	"github.com/jace-ys/kiuru/backend/service.user/pkg/user"
)

var logger log.Logger

func main() {
	c := parseCommand()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "source", log.DefaultCaller)

	crdbClient, err := crdb.NewCRDBClient(c.database.Host, c.database.User, c.database.DBName)
	if err != nil {
		exit(err)
	}
	redisClient, err := redis.NewRedisClient(c.redis.Host)
	if err != nil {
		exit(err)
	}

	userService, err := user.NewService(logger, crdbClient)
	if err != nil {
		exit(err)
	}
	defer userService.Teardown()

	authInterceptor := gorpc.NewAuthInterceptor(
		authr.NewJWTAuthenticator(c.jwtSecretKey, cache.NewRedisCache(redisClient)),
		userService.GetAuthenticatedMethods(),
	)

	middleware := grpc.ChainUnaryInterceptor(
		authInterceptor.Authenticate(),
	)

	grpcServer := server.NewGRPCServer(c.server.Port, middleware)
	gatewayProxy := server.NewGatewayProxy(c.gateway.Port, c.gateway.Endpoint)

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
	gateway      server.GatewayProxyConfig
	database     crdb.Config
	redis        redis.Config
	jwtSecretKey string
}

func parseCommand() *config {
	var c config

	kingpin.Flag("port", "port for the gRPC server").Default("8080").IntVar(&c.server.Port)
	kingpin.Flag("gateway-port", "port for the REST gateway proxy").Default("8081").IntVar(&c.gateway.Port)
	kingpin.Flag("crdb-host", "host for connecting to CockroachDB").Default("127.0.0.1:26257").StringVar(&c.database.Host)
	kingpin.Flag("crdb-user", "user for connecting to CockroachDB").Default("default").StringVar(&c.database.User)
	kingpin.Flag("crdb-dbname", "database name for connecting to CockroachDB").Default("defaultdb").StringVar(&c.database.DBName)
	kingpin.Flag("redis-host", "host for connecting Redis").Default("127.0.0.1:6379").StringVar(&c.redis.Host)
	kingpin.Flag("jwt-secret", "secret key used to sign JWTs").Required().StringVar(&c.jwtSecretKey)
	kingpin.Parse()

	c.gateway.Endpoint = fmt.Sprintf(":%d", c.server.Port)
	return &c
}

func exit(err error) {
	level.Error(logger).Log("event", "service.fatal", "msg", err)
	os.Exit(1)
}
