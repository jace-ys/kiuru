TARGET = main
SOURCE = main.go

.PHONY: all build run test format update clean proto

all: format build run

build:
	@echo "==> Building from source.."
	go build -o ${TARGET} ${SOURCE}

run:
	@echo "==> Executing binary.."
	./serve-dev.sh

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
	../../bin/protogen.sh -l go -t user -p ../../api/proto -g
