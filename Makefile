.PHONY: build build-win build-mac run clean help

APP_NAME := app
BIN_DIR := bin
MAIN := ./cmd/app

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

help:
	@echo "Available targets:"
	@echo "  build       - Build for current OS"
	@echo "  build-win   - Build for Windows (amd64)"
	@echo "  build-mac   - Build for Mac (arm64/Apple Silicon)"
	@echo "  run         - Run the application with go run"
	@echo "  clean       - Remove bin directory"

build:
	@echo "Building for current OS..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME)$(EXE_EXT) $(MAIN)

build-win:
	@echo "Building for Windows (amd64)..."
	@mkdir -p $(BIN_DIR)
	@GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME).exe $(MAIN)

build-mac:
	@echo "Building for Mac (arm64)..."
	@mkdir -p $(BIN_DIR)
	@GOOS=darwin GOARCH=arm64 go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN)

run:
	@go run $(MAIN)

clean:
	@if exist $(BIN_DIR) (rmdir /s /q $(BIN_DIR)) else (rm -rf $(BIN_DIR))
