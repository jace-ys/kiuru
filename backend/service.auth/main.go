package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kru-travel/airdrop-go/pkg/authr"
	"github.com/kru-travel/airdrop-go/pkg/crdb"
	"github.com/kru-travel/airdrop-go/pkg/redis"
	"golang.org/x/sync/errgroup"

	"github.com/jace-ys/kru-travel/backend/service.auth/pkg/auth"
	"github.com/jace-ys/kru-travel/backend/service.auth/pkg/server"
)

var logger log.Logger

func main() {
	c := parseCommand()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
	gatewayProxy := server.NewGatewayProxy(c.gateway)

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
	gateway  server.GatewayConfig
	database crdb.Config
	redis    redis.Config
	jwt      authr.JWTConfig
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
	kingpin.Flag("redis-host", "host for connecting Redis").Default("127.0.0.1").StringVar(&c.redis.Host)
	kingpin.Flag("redis-port", "port for connecting to Redis").Default("6379").IntVar(&c.redis.Port)
	kingpin.Flag("jwt-secret", "secret key used to sign JWTs").StringVar(&c.jwt.SecretKey)
	kingpin.Flag("jwt-issuer", "issuer of generated JWTs").StringVar(&c.jwt.Issuer)
	kingpin.Flag("jwt-ttl", "time-to-live for generated JWTs").Default("15m").DurationVar(&c.jwt.TTL)
	kingpin.Parse()

	c.gateway.Endpoint = fmt.Sprintf("%s:%d", c.server.Host, c.server.Port)
	return &c
}

func exit(err error) {
	level.Error(logger).Log("event", "service.fatal", "msg", err)
	os.Exit(1)
}
