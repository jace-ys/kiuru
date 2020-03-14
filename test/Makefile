PROTO = v1/auth v1/user

.PHONY: all integration setup teardown proto format clean

all:
	@make teardown
	@make setup
	go test -v ./integration

integration:
	@make teardown
	@make setup
	@echo "==> Running integration tests.."
	go test -v ./integration -run ${SUITE}

images:
	@echo "==> Building images.."
	docker-compose build

setup:
	@echo "==> Starting containers.."
	docker-compose -f ../docker-compose.yaml -f docker-compose.test.yaml up -d
	@echo "==> Initializing.."
	@sleep 3

teardown:
	@echo "==> Killing containers.."
	docker-compose down -v

proto:
	@echo "==> Generating proto.."
	../bin/protogen -t "${PROTO}" -p ../api/proto

format:
	@echo "==> Formatting code.."
	gofmt -w .

clean:
	@echo "==> Cleaning up.."
	go clean
