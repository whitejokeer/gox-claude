.PHONY: help build test lint clean install dev setup
.DEFAULT_GOAL := help

# Variables
BINARY_NAME=gox
BUILD_DIR=dist
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X github.com/gox-framework/gox/pkg/version.Version=$(VERSION)"

## help: Display this help message
help:
	@echo "GOX Framework - Available commands:"
	@echo ""
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
	@echo ""

## setup: Install development dependencies
setup:
	@echo "Installing development dependencies..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/goreleaser/goreleaser@latest
	@go mod download
	@echo "Setup complete!"

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/gox
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## install: Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS) ./cmd/gox
	@echo "Installed $(BINARY_NAME) to $(shell go env GOPATH)/bin"

## test: Run all tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## test-short: Run tests without integration tests
test-short:
	@echo "Running short tests..."
	@go test -short -v ./...

## lint: Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

## lint-fix: Run linter with auto-fix
lint-fix:
	@echo "Running linter with auto-fix..."
	@golangci-lint run --fix

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@go clean -cache
	@echo "Clean complete!"

## dev: Start development mode
dev:
	@echo "Starting development mode..."
	@go run ./cmd/gox dev

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

## mod-tidy: Tidy Go modules
mod-tidy:
	@echo "Tidying modules..."
	@go mod tidy

## check: Run all checks (lint, test, build)
check: lint test build
	@echo "All checks passed!"

## release: Create a release (requires goreleaser)
release:
	@echo "Creating release..."
	@goreleaser release --rm-dist

## release-snapshot: Create a snapshot release
release-snapshot:
	@echo "Creating snapshot release..."
	@goreleaser release --snapshot --rm-dist