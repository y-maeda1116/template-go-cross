.PHONY: build build-cli build-desktop run-cli run-desktop clean test test-coverage test-race mocks lint fmt help

APP_NAME := myapp
BIN_DIR := bin

# CLI
CLI_MAIN := ./cmd/cli

# Detect OS
ifeq ($(OS),Windows_NT)
    EXE_EXT := .exe
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        EXE_EXT :=
    endif
    ifeq ($(UNAME_S),Darwin)
        EXE_EXT :=
    endif
endif

# --- ビルド ---

build-cli:
	@echo "Building CLI for current OS..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME)$(EXE_EXT) $(CLI_MAIN)

build-desktop:
	@echo "Building Desktop for current OS..."
	@test -d frontend && cd frontend && npm install && npm run build || echo "Frontend directory not found"
	@wails build

build-all: build-cli build-desktop

# --- 実行 ---

run-cli:
	@echo "Running CLI..."
	@go run $(CLI_MAIN) $(ARGS)

run-desktop:
	@echo "Running Desktop..."
	@wails dev

# --- テスト ---

test:
	@echo "Running tests..."
	@go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

test-race:
	@echo "Running tests with race detector..."
	@go test -race -v ./...

# --- モック生成 ---

mocks:
	@echo "Generating mocks..."
	@mkdir -p test/mocks
	@mockgen -source=internal/core/service.go -destination=test/mocks/service_mock.go

# --- Lint / Format ---

fmt:
	@go fmt ./...

lint:
	@echo "Installing golangci-lint if needed..."
	@which golangci-lint || (curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin latest)
	@echo "Running linter..."
	@golangci-lint run ./...

# --- クリーンアップ ---

clean:
	@rm -rf $(BIN_DIR) coverage.out coverage.html test/mocks
	@test -d frontend && cd frontend && rm -rf node_modules wailsjs dist || true

# --- ヘルプ ---

help:
	@echo "Available targets:"
	@echo "  build-cli       - Build CLI for current OS"
	@echo "  build-desktop    - Build Desktop for current OS"
	@echo "  build-all        - Build CLI and Desktop"
	@echo "  run-cli          - Run CLI (use ARGS=\"--help\" for options)"
	@echo "  run-desktop      - Run Desktop in dev mode"
	@echo "  test             - Run all tests"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  test-race        - Run tests with race detector"
	@echo "  mocks            - Generate mocks"
	@echo "  fmt              - Format Go code"
	@echo "  lint             - Run linter"
	@echo "  clean            - Remove build artifacts"
