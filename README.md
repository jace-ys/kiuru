[![client.web][client.web-badge]][client.web-workflow]
[![service.auth][service.auth-badge]][service.auth-workflow]
[![service.user][service.user-badge]][service.user-workflow]
[![test.integration][test.integration-badge]][test.integration-workflow]

# Kiuru

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
  make run

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

[client.web-badge]: https://github.com/jace-ys/kiuru/workflows/client.web/badge.svg
[client.web-workflow]: https://github.com/jace-ys/kiuru/actions?query=workflow%3Atest.integration
[service.auth-badge]: https://github.com/jace-ys/kiuru/workflows/service.auth/badge.svg
[service.auth-workflow]: https://github.com/jace-ys/kiuru/actions?query=workflow%3Aservice.auth
[service.user-badge]: https://github.com/jace-ys/kiuru/workflows/service.user/badge.svg
[service.user-workflow]: https://github.com/jace-ys/kiuru/actions?query=workflow%3Aservice.user
[test.integration-badge]: https://github.com/jace-ys/kiuru/workflows/test.integration/badge.svg
[test.integration-workflow]: https://github.com/jace-ys/kiuru/actions?query=workflow%3Atest.integration
