# Kru Travel

![.github/workflows/client.web.yml](https://github.com/jace-ys/kru-travel/workflows/.github/workflows/client.web.yml/badge.svg)
![.github/workflows/service.auth.yml](https://github.com/jace-ys/kru-travel/workflows/.github/workflows/service.auth.yml/badge.svg)
![.github/workflows/service.user.yml](https://github.com/jace-ys/kru-travel/workflows/.github/workflows/service.user.yml/badge.svg)
![.github/workflows/test.integration.yml](https://github.com/jace-ys/kru-travel/workflows/.github/workflows/test.integration.yml/badge.svg)

## Prerequisites

#### Dependencies:

- docker, docker-compose
- golang
- node, npm
- protoc, protoc-gen-go, protoc-gen-grpc-gateway

#### Development Setup:

- Generate gRPC stubs from proto files

```
* Go:
  make proto
```

- Start all containers

```
docker-compose up
```

- Start auxiliary containers

```
docker-compose up -d db.cockroach db.cockroach.init db.redis
```

- Start individual services

```
* Go:
  make

* Node:
  npm start
```

- Test individual services

```
* Go:
  make test

* Node:
  npm test
```
