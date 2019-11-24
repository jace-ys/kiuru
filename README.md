# Kru Travel

[![CircleCI](https://circleci.com/gh/jace-ys/kru-travel.svg?style=svg&circle-token=86078b9731d4274ee92fb405f89a2fa3e4cf6bc5)](https://circleci.com/gh/jace-ys/kru-travel)

## Prerequisites

#### Dependencies:

- docker, docker-compose
- golang
- node, npm
- protoc, protoc-gen-go, protoc-gen-grpc-gateway

#### Development Setup:

- Start all containers

```
docker-compose up
```

- Start database containers

```
docker-compose up -d db.cockroach db.cockroach.init db.redis
```

- Run individual services

```
* Go:
  make test
  make

* Node:
  npm test
  npm start
```

- Generate gRPC stubs from proto files

```
* Go:
  make proto
```
