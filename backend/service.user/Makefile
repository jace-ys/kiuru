TARGET = main
SOURCE = main.go

.PHONY: all build run test format clean

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
