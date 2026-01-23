.PHONY: build run test clean docker

APP_NAME := servers-manager
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# Build the application
build:
	go build $(LDFLAGS) -o $(APP_NAME) ./app

# Run the application
run: build
	./$(APP_NAME) --debug

# Run with specific database
run-db: build
	./$(APP_NAME) --db=servers.db --debug

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -f $(APP_NAME)
	rm -f coverage.out coverage.html
	rm -f servers.db

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run ./...

# Generate enum code (requires go-pkgz/enum)
generate:
	go generate ./...

# Download dependencies
deps:
	go mod download
	go mod tidy

# Build Docker image
docker:
	docker build -t $(APP_NAME):$(VERSION) .

# Run Docker container
docker-run:
	docker run -p 8080:8080 -v $(PWD)/data:/data $(APP_NAME):$(VERSION)

# Development with auto-reload (requires air)
dev:
	air

# Install development tools
tools:
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

help:
	@echo "Available targets:"
	@echo "  build       - Build the application"
	@echo "  run         - Build and run the application"
	@echo "  run-db      - Build and run with servers.db"
	@echo "  test        - Run tests"
	@echo "  test-cover  - Run tests with coverage report"
	@echo "  clean       - Remove build artifacts"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  generate    - Generate enum code"
	@echo "  deps        - Download and tidy dependencies"
	@echo "  docker      - Build Docker image"
	@echo "  docker-run  - Run Docker container"
	@echo "  dev         - Run with auto-reload (requires air)"
	@echo "  tools       - Install development tools"
