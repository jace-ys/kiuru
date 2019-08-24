TARGET = main
SOURCE = main.go
PROTO_TARGETS = user

.PHONY: all build run test format clean proto

all: format build run

build:
	@echo "==> Building from source.."
	go build -o ${TARGET} ${SOURCE}

run:
	@echo "==> Executing binary.."
	@./${TARGET}

test:
	@echo "==> Running tests.."
	go test ./... -v

format:
	@echo "==> Formatting code.."
	gofmt -w .

clean:
	@echo "==> Cleaning up.."
	go mod tidy
	go clean

proto:
	../../bin/protogen.sh -l go -t ${PROTO_TARGETS} -p ../../api/proto
