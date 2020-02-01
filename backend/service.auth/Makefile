TARGET = service
SOURCE = main.go
PROTO = v1/auth
DEPENDENCIES = db.cockroach db.cockroach.init db.redis

.PHONY: all build run test format update clean proto dependencies

all: format build run

build:
	@echo "==> Building from source.."
	go build -o ${TARGET} ${SOURCE}

run:
	@echo "==> Executing binary.."
	./serve-dev

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
