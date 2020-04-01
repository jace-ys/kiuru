PROTO = v1/auth v1/user

.PHONY: all integration setup teardown proto migrations format

all:
	@make teardown
	@make setup
	@echo "==> Initializing.." && sleep 3
	@echo "==> Running integration tests.."
	go test -v ./integration

integration:
	@make teardown
	@make setup
	@echo "==> Initializing.." && sleep 3
	@echo "==> Running integration tests.."
	go test -v ./integration -run ${SUITE}

images:
	@echo "==> Building images.."
	docker-compose build

setup:
	@echo "==> Starting containers.."
	docker-compose -f ../docker-compose.yaml -f docker-compose.test.yaml up -d

teardown:
	@echo "==> Killing containers.."
	docker-compose down -v

proto:
	@echo "==> Generating proto.."
	../bin/protogen -t "${PROTO}" -p ../api/proto

migrations:
	@echo "==> Copying migrations.."
	rsync -a --ignore-existing ../databases/cockroach/migrations/ data/cockroach/migrations/

format:
	@echo "==> Formatting code.."
	gofmt -w .
