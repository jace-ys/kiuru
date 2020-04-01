package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kiuru-travel/airdrop-go/authr"
	"github.com/kiuru-travel/airdrop-go/cache"
	"github.com/kiuru-travel/airdrop-go/crdb"
	"github.com/kiuru-travel/airdrop-go/redis"
	"golang.org/x/sync/errgroup"

	"github.com/jace-ys/kiuru/backend/service.auth/pkg/auth"
	"github.com/jace-ys/kiuru/backend/service.auth/pkg/server"
)

var logger log.Logger

func main() {
	c := parseCommand()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	crdbClient, err := crdb.NewCRDBClient(c.crdb.Host, c.crdb.User, c.crdb.Database)
	if err != nil {
		exit(err)
	}
	redisClient, err := redis.NewRedisClient(c.redis.Host)
	if err != nil {
		exit(err)
	}
	jwtHandler := authr.NewJWTHandler(c.jwt.Issuer, c.jwt.SecretKey, c.jwt.TTL, cache.NewRedisCache(redisClient))

	authService, err := auth.NewService(logger, crdbClient, jwtHandler)
	if err != nil {
		exit(err)
	}
	defer authService.Teardown()

	grpcServer := server.NewGRPCServer(c.server.Port)
	gatewayProxy := server.NewGatewayProxy(c.proxy.Port, c.proxy.Endpoint)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		level.Info(logger).Log("event", "grpc_server.started", "port", c.server.Port)
		defer level.Info(logger).Log("event", "grpc_server.stopped")
		return authService.StartServer(ctx, grpcServer)
	})
	g.Go(func() error {
		level.Info(logger).Log("event", "gateway_proxy.started", "port", c.proxy.Port)
		defer level.Info(logger).Log("event", "gateway_proxy.stopped")
		return authService.StartServer(ctx, gatewayProxy)
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
	server server.GRPCServerConfig
	proxy  server.GatewayProxyConfig
	crdb   crdb.Config
	redis  redis.Config
	jwt    authr.JWTHandlerConfig
}

func parseCommand() *config {
	var c config

	kingpin.Flag("port", "port for the gRPC server").Default("8081").IntVar(&c.server.Port)
	kingpin.Flag("gateway-port", "port for the REST gateway proxy").Default("8080").IntVar(&c.proxy.Port)
	kingpin.Flag("crdb-host", "host for connecting to CockroachDB").Default("127.0.0.1:26257").StringVar(&c.crdb.Host)
	kingpin.Flag("crdb-user", "user for connecting to CockroachDB").Default("default").StringVar(&c.crdb.User)
	kingpin.Flag("crdb-database", "database name for connecting to CockroachDB").Default("defaultdb").StringVar(&c.crdb.Database)
	kingpin.Flag("redis-host", "host for connecting Redis").Default("127.0.0.1:6379").StringVar(&c.redis.Host)
	kingpin.Flag("jwt-secret", "secret key used to sign JWTs").Required().StringVar(&c.jwt.SecretKey)
	kingpin.Flag("jwt-issuer", "issuer of generated JWTs").Default("").StringVar(&c.jwt.Issuer)
	kingpin.Flag("jwt-ttl", "time to live for generated JWTs").Default("15m").DurationVar(&c.jwt.TTL)
	kingpin.Parse()

	c.proxy.Endpoint = fmt.Sprintf(":%d", c.server.Port)
	return &c
}

func exit(err error) {
	level.Error(logger).Log("event", "service.fatal", "msg", err)
	os.Exit(1)
}
