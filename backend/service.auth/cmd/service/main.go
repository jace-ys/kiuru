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
	"github.com/kiuru-travel/airdrop-go/pkg/redis"
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
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "source", log.DefaultCaller)

	crdbClient, err := crdb.NewCRDBClient(c.database.Host, c.database.Port, c.database.User, c.database.DBName)
	if err != nil {
		exit(err)
	}
	redisClient, err := redis.NewRedisClient(c.redis.Host, c.redis.Port)
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
	gatewayProxy := server.NewGatewayProxy(c.gateway.Port, c.gateway.Endpoint)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		level.Info(logger).Log("event", "grpc_server.started", "port", c.server.Port)
		defer level.Info(logger).Log("event", "grpc_server.stopped")
		return authService.StartServer(ctx, grpcServer)
	})
	g.Go(func() error {
		level.Info(logger).Log("event", "gateway_proxy.started", "port", c.gateway.Port)
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
	server   server.GRPCServerConfig
	gateway  server.GatewayProxyConfig
	database crdb.Config
	redis    redis.Config
	jwt      authr.JWTHandlerConfig
}

func parseCommand() *config {
	var c config

	kingpin.Flag("port", "port for the gRPC server").Default("8080").IntVar(&c.server.Port)
	kingpin.Flag("gateway-port", "port for the REST gateway proxy").Default("8081").IntVar(&c.gateway.Port)
	kingpin.Flag("crdb-host", "host for connecting to CockroachDB").Default("127.0.0.1").StringVar(&c.database.Host)
	kingpin.Flag("crdb-port", "port for connecting to CockroachDB").Default("26257").IntVar(&c.database.Port)
	kingpin.Flag("crdb-user", "user for connecting to CockroachDB").Default("default").StringVar(&c.database.User)
	kingpin.Flag("crdb-dbname", "database name for connecting to CockroachDB").Default("defaultdb").StringVar(&c.database.DBName)
	kingpin.Flag("redis-host", "host for connecting Redis").Default("127.0.0.1").StringVar(&c.redis.Host)
	kingpin.Flag("redis-port", "port for connecting to Redis").Default("6379").IntVar(&c.redis.Port)
	kingpin.Flag("jwt-secret", "secret key used to sign JWTs").Required().StringVar(&c.jwt.SecretKey)
	kingpin.Flag("jwt-issuer", "issuer of generated JWTs").Default("").StringVar(&c.jwt.Issuer)
	kingpin.Flag("jwt-ttl", "time-to-live for generated JWTs").Default("15m").DurationVar(&c.jwt.TTL)
	kingpin.Parse()

	c.gateway.Endpoint = fmt.Sprintf(":%d", c.server.Port)
	return &c
}

func exit(err error) {
	level.Error(logger).Log("event", "service.fatal", "msg", err)
	os.Exit(1)
}
