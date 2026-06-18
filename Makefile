BINARY   := bin/nagiosql
MODULE   := go-nagiosql
VERSION  := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILDDATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS  := -s -w \
	-X main.Version=$(VERSION) \
	-X main.BuildDate=$(BUILDDATE)

.PHONY: all build clean tidy fmt vet lint swagger test test-integration db-start db-stop db-reset server-start server-stop test-api check docker-up docker-down docker-logs nagioscore-up nagioscore-down

all: build

## build: compile the nagiosql binary into bin/
build:
	@mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

## tidy: tidy and verify Go modules
tidy:
	go mod tidy
	go mod verify

## fmt: format all Go files
fmt:
	gofmt -w -s ./...

## vet: run go vet
vet:
	go vet ./...

## lint: run golangci-lint (install separately)
lint:
	golangci-lint run ./...

## test: run unit tests (no external dependencies required)
test:
	go test -race -cover ./...

## test-integration: run integration tests against a live MariaDB (requires make db-start first)
test-integration:
	go test -tags integration -v -count=1 ./internal/integration/...

## swagger: generate OpenAPI docs from code annotations
swagger:
	swag init --generalInfo main.go --output docs --parseDependency --parseInternal

## clean: remove build artefacts
clean:
	rm -rf bin/
	rm -rf docs/swagger.json docs/swagger.yaml docs/docs.go

## run: build and start the server (requires config.toml)
run: build
	./$(BINARY) serve

## db-start: start the test MariaDB container on port 3307
db-start:
	docker compose -f docker/test/docker-compose.db.yml up -d --wait
	@echo "MariaDB test ready on :3307"

## db-stop: stop test MariaDB container (keeps volume; use db-reset to wipe)
db-stop:
	docker compose -f docker/test/docker-compose.db.yml down

## db-reset: wipe and restart test DB
db-reset: db-stop db-start

## server-start: build, migrate (with sample), and start the server in background
server-start: build
	@cp -n config.toml.example config.toml 2>/dev/null || true
	./$(BINARY) migrate --admin-password admin123 --sample --config config.toml 2>/dev/null || true
	./$(BINARY) serve --config config.toml &
	@echo $$! > .server.pid
	@until curl -sf http://localhost:8081/healthz > /dev/null 2>&1; do sleep 1; done
	@echo "server ready"

## server-stop: stop the background server
server-stop:
	@if [ -f .server.pid ]; then kill $$(cat .server.pid) 2>/dev/null; rm .server.pid; fi

## test-api: run all bash smoke tests against a local server
test-api: build
	@echo "Starting server for API tests..."
	$(MAKE) server-start || true
	BASE_URL=http://localhost:8081 bash test/api/smoke.sh; CODE=$$?; $(MAKE) server-stop; exit $$CODE

## check: vet + build + test (CI entry point)
check: vet build test

## docker-up: build and start the go-nagiosql stack (nagios4 + API + MariaDB)
docker-up:
	docker compose -f docker/go-nagiosql/docker-compose.yml up -d --build

## docker-down: stop the go-nagiosql stack
docker-down:
	docker compose -f docker/go-nagiosql/docker-compose.yml down

## docker-logs: tail logs from the go-nagiosql stack
docker-logs:
	docker compose -f docker/go-nagiosql/docker-compose.yml logs -f

## nagioscore-up: start the reference nagios-core stack (DOCUMENTS/)
nagioscore-up:
	docker compose -f DOCUMENTS/docker/nagios-core/docker-compose.yml up -d

## nagioscore-down: stop the reference nagios-core stack
nagioscore-down:
	docker compose -f DOCUMENTS/docker/nagios-core/docker-compose.yml down
