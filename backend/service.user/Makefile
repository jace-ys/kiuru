TARGET = service
SOURCE = cmd/service/main.go
PROTO = v1/user
DEPENDENCIES = db.cockroach db.cockroach.init db.redis

.PHONY: all build run proto dependencies test format

all: format run

build:
	@echo "==> Compiling code.."
	go build -o ${TARGET} ${SOURCE}

run:
	@echo "==> Executing code.."
	@go run ${SOURCE} \
		--port 8081 \
		--gateway-port 8080 \
		--crdb-host 127.0.0.1:26257 \
		--crdb-user kiuru \
		--crdb-database kiuru \
		--redis-host 127.0.0.1:6379 \
		--jwt-secret my-secret-key

proto:
	@echo "==> Generating proto.."
	../../bin/protogen -t "${PROTO}" -p ../../api/proto -g

dependencies:
	@echo "==> Starting auxiliary containers.."
	docker-compose up -d ${DEPENDENCIES}

test:
	@echo "==> Running tests.."
	go test -v ./...

format:
	@echo "==> Formatting code.."
	gofmt -w .
