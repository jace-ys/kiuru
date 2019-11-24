TARGET = service
SOURCE = main.go
PROTO = user
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
	../../bin/protogen -l go -t ${PROTO} -p ../../api/proto -g

dependencies:
	docker-compose up -d ${DEPENDENCIES}
