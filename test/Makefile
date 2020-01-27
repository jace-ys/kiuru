PROTO = v1/auth v1/user

.PHONY: all setup integration teardown format update clean proto

all:
	@make teardown 2>/dev/null
	@make setup 2>/dev/null
	@echo "==> Initialising.."
	@sleep 3
	@make integration

setup:
	@echo "==> Starting containers.."
	docker-compose -f ../docker-compose.yml -f docker-compose.test.yml up -d

integration:
	@echo "==> Running integration tests.."
	go test ./integration -v

teardown:
	@echo "==> Killing containers.."
	docker-compose down -v

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
	../bin/protogen -l go -t "${PROTO}" -p ../api/proto
