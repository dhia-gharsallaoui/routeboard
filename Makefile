.PHONY: build run test lint clean dev web-dev web-build docker-build helm-lint helm-template

BINARY  := routeboard
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -X main.version=$(VERSION)
GOFLAGS := -trimpath -ldflags "$(LDFLAGS)"

## Build Go binary with embedded frontend
build: web-build embed
	go build $(GOFLAGS) -o bin/$(BINARY) ./cmd/routeboard

## Copy web dist into embed location
embed:
	rm -rf internal/server/dist
	cp -r web/dist internal/server/dist

## Run Go server locally (without embedded frontend — use web-dev proxy)
run:
	go run ./cmd/routeboard

## Run tests
test:
	go test ./... -race -cover -count=1

## Lint
lint:
	golangci-lint run ./...

## Build React frontend
web-build:
	cd web && bun install && bun run build

## Dev: run React dev server (proxies /api to Go backend on :8080)
web-dev:
	cd web && bun run dev

## Docker
docker-build:
	docker build -t $(BINARY):$(VERSION) .

## Helm
helm-lint:
	helm lint deploy/helm/routeboard

helm-template:
	helm template routeboard deploy/helm/routeboard

## Clean
clean:
	rm -rf bin/ web/dist/ web/node_modules/ internal/server/dist/

## Dev: run both Go backend and React frontend
dev:
	@echo "Run 'make run' and 'make web-dev' in separate terminals"
