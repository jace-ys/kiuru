TARGET = service
SOURCE = main.go
PROTO = v1/user
DEPENDENCIES = db.cockroach db.cockroach.init db.redis

.PHONY: all build run test format update clean proto dependencies

build:
	@echo "==> Building from source.."
	go build -o ${TARGET} ${SOURCE}

run:
	@echo "==> Executing binary.."
	@go run main.go \
	--port 8080 \
	--gateway-port 8081 \
	--crdb-host 127.0.0.1 \
	--crdb-port 26257 \
	--crdb-user kru_service \
	--crdb-dbname kru \
	--redis-host 127.0.0.1 \
	--redis-port 6379 \
	--jwt-secret my-secret-key

test:
	@echo "==> Running tests.."
	go test ./... -v

format:
	@echo "==> Formatting code.."
	gofmt -w .

update:
	@echo "==> Updating dependencies.."
	go get -u
	go mod tidy

clean:
	@echo "==> Cleaning up.."
	go clean

proto:
	@echo "==> Generating proto.."
	../../bin/proto-gen -l go -t "${PROTO}" -p ../../api/proto -g

dependencies:
	@echo "==> Starting auxiliary containers.."
	docker-compose up -d ${DEPENDENCIES}
