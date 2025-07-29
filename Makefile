# Project variables
APP_NAME = goreportx
PKG = ./...
BUILD_DIR = bin
MAIN = ./cmd/goreportx/main.go
OUT = $(BUILD_DIR)/$(APP_NAME)

# Go tools
GO = go
GOTEST = $(GO) test
GOBUILD = $(GO) build
GOLINT = golangci-lint

# Default target
.PHONY: all
all: build

# Build the CLI binary
.PHONY: build
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(OUT) $(MAIN)

# Run all tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v $(PKG)

# Run tests with coverage
.PHONY: cover
cover:
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out $(shell go list ./... | grep -vE "/examples|/internal/template")
	@go tool cover -html=coverage.out

# Lint the codebase
.PHONY: lint
lint:
	@echo "Linting..."
	$(GOLINT) run

# Format code using gofmt
.PHONY: fmt
fmt:
	@echo "Formatting..."
	go fmt $(PKG)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR) coverage.out

# Install the CLI locally (into your $GOBIN)
.PHONY: install
install:
	@echo "Installing..."
	$(GO) install $(MAIN)

# Run the CLI with optional args
.PHONY: run
run:
	@echo "Running goreportx..."
	$(GO) run $(MAIN) --input examples/data.json --template examples/template.html --format json

