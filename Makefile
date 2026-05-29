.PHONY: help build run dev test lint clean docker-build docker-run

help:
	@echo "Service Order Bot — Make Commands"
	@echo ""
	@echo "  make build           Build binary"
	@echo "  make run             Build and run"
	@echo "  make dev             Run with hot reload (requires air)"
	@echo "  make lint            Run linters"
	@echo "  make clean           Remove build artifacts"
	@echo "  make docker-build    Build Docker image"
	@echo "  make docker-run      Run full stack (bot + api) via docker-compose"
	@echo "  make docker-stop     Stop docker-compose stack"
	@echo ""

build:
	@echo "Building..."
	@mkdir -p bin
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/service-order-bot ./cmd/bot
	@echo "✅ Built: bin/service-order-bot"

run: build
	./bin/service-order-bot

dev:
	@which air > /dev/null || go install github.com/cosmtrek/air@latest
	air -c .air.toml

lint:
	go vet ./...
	@which golangci-lint > /dev/null && golangci-lint run ./... || echo "golangci-lint not installed"

clean:
	rm -rf bin/
	go clean

deps:
	go mod download
	go mod tidy

docker-build:
	docker build -t service-order-bot:latest .

docker-run:
	docker-compose up --build

docker-stop:
	docker-compose down