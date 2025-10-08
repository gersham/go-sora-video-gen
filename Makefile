# Video Generator - Build Configuration
.PHONY: all build build-all clean dist help

# Binary name
BINARY_NAME=video-gen
VERSION?=1.0.0

# Build directories
DIST_DIR=./dist
RELEASE_DIR=./releases

# Build flags
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION)"
BUILD_FLAGS=-trimpath

# Platforms
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build optimized binary for current platform
	@echo "Building $(BINARY_NAME) for current platform..."
	@mkdir -p $(DIST_DIR)
	go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME) .
	@echo "✓ Binary created at $(DIST_DIR)/$(BINARY_NAME)"

build-all: build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-linux-arm64 build-windows-amd64 ## Build for all platforms

build-darwin-amd64: ## Build for macOS Intel
	@echo "Building for macOS (Intel)..."
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@echo "✓ Binary created at $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64"

build-darwin-arm64: ## Build for macOS Apple Silicon
	@echo "Building for macOS (Apple Silicon)..."
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@echo "✓ Binary created at $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64"

build-linux-amd64: ## Build for Linux amd64
	@echo "Building for Linux (amd64)..."
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .
	@echo "✓ Binary created at $(DIST_DIR)/$(BINARY_NAME)-linux-amd64"

build-linux-arm64: ## Build for Linux ARM64
	@echo "Building for Linux (ARM64)..."
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 .
	@echo "✓ Binary created at $(DIST_DIR)/$(BINARY_NAME)-linux-arm64"

build-windows-amd64: ## Build for Windows amd64
	@echo "Building for Windows (amd64)..."
	@mkdir -p $(DIST_DIR)
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "✓ Binary created at $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe"

dist: dist-darwin-amd64 dist-darwin-arm64 dist-linux-amd64 dist-linux-arm64 dist-windows-amd64 ## Create distribution archives for all platforms
	@echo ""
	@echo "✓ All distribution archives created in $(RELEASE_DIR)/"
	@ls -lh $(RELEASE_DIR)/

dist-darwin-amd64: build-darwin-amd64 ## Create macOS Intel distribution
	@echo "Creating distribution for macOS (Intel)..."
	@mkdir -p $(RELEASE_DIR)/tmp/$(BINARY_NAME)
	@cp $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(RELEASE_DIR)/tmp/$(BINARY_NAME)/$(BINARY_NAME)
	@cp README.md $(RELEASE_DIR)/tmp/$(BINARY_NAME)/
	@if [ -f LICENSE ]; then cp LICENSE $(RELEASE_DIR)/tmp/$(BINARY_NAME)/; fi
	@cd $(RELEASE_DIR)/tmp && tar -czf ../$(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)
	@rm -rf $(RELEASE_DIR)/tmp
	@echo "✓ Archive created: $(RELEASE_DIR)/$(BINARY_NAME)-darwin-amd64.tar.gz"

dist-darwin-arm64: build-darwin-arm64 ## Create macOS Apple Silicon distribution
	@echo "Creating distribution for macOS (Apple Silicon)..."
	@mkdir -p $(RELEASE_DIR)/tmp/$(BINARY_NAME)
	@cp $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(RELEASE_DIR)/tmp/$(BINARY_NAME)/$(BINARY_NAME)
	@cp README.md $(RELEASE_DIR)/tmp/$(BINARY_NAME)/
	@if [ -f LICENSE ]; then cp LICENSE $(RELEASE_DIR)/tmp/$(BINARY_NAME)/; fi
	@cd $(RELEASE_DIR)/tmp && tar -czf ../$(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)
	@rm -rf $(RELEASE_DIR)/tmp
	@echo "✓ Archive created: $(RELEASE_DIR)/$(BINARY_NAME)-darwin-arm64.tar.gz"

dist-linux-amd64: build-linux-amd64 ## Create Linux amd64 distribution
	@echo "Creating distribution for Linux (amd64)..."
	@mkdir -p $(RELEASE_DIR)/tmp/$(BINARY_NAME)
	@cp $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(RELEASE_DIR)/tmp/$(BINARY_NAME)/$(BINARY_NAME)
	@cp README.md $(RELEASE_DIR)/tmp/$(BINARY_NAME)/
	@if [ -f LICENSE ]; then cp LICENSE $(RELEASE_DIR)/tmp/$(BINARY_NAME)/; fi
	@cd $(RELEASE_DIR)/tmp && tar -czf ../$(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)
	@rm -rf $(RELEASE_DIR)/tmp
	@echo "✓ Archive created: $(RELEASE_DIR)/$(BINARY_NAME)-linux-amd64.tar.gz"

dist-linux-arm64: build-linux-arm64 ## Create Linux ARM64 distribution
	@echo "Creating distribution for Linux (ARM64)..."
	@mkdir -p $(RELEASE_DIR)/tmp/$(BINARY_NAME)
	@cp $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 $(RELEASE_DIR)/tmp/$(BINARY_NAME)/$(BINARY_NAME)
	@cp README.md $(RELEASE_DIR)/tmp/$(BINARY_NAME)/
	@if [ -f LICENSE ]; then cp LICENSE $(RELEASE_DIR)/tmp/$(BINARY_NAME)/; fi
	@cd $(RELEASE_DIR)/tmp && tar -czf ../$(BINARY_NAME)-linux-arm64.tar.gz $(BINARY_NAME)
	@rm -rf $(RELEASE_DIR)/tmp
	@echo "✓ Archive created: $(RELEASE_DIR)/$(BINARY_NAME)-linux-arm64.tar.gz"

dist-windows-amd64: build-windows-amd64 ## Create Windows amd64 distribution
	@echo "Creating distribution for Windows (amd64)..."
	@mkdir -p $(RELEASE_DIR)/tmp/$(BINARY_NAME)
	@cp $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe $(RELEASE_DIR)/tmp/$(BINARY_NAME)/$(BINARY_NAME).exe
	@cp README.md $(RELEASE_DIR)/tmp/$(BINARY_NAME)/
	@if [ -f LICENSE ]; then cp LICENSE $(RELEASE_DIR)/tmp/$(BINARY_NAME)/; fi
	@cd $(RELEASE_DIR)/tmp && zip -r ../$(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)
	@rm -rf $(RELEASE_DIR)/tmp
	@echo "✓ Archive created: $(RELEASE_DIR)/$(BINARY_NAME)-windows-amd64.zip"

clean: ## Remove build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(DIST_DIR) $(RELEASE_DIR)
	@rm -f $(BINARY_NAME)
	@echo "✓ Clean complete"

test: ## Run tests
	go test -v -race ./...

fmt: ## Format code
	go fmt ./...

lint: ## Run linter
	golangci-lint run

deps: ## Download dependencies
	go mod download
	go mod tidy

.DEFAULT_GOAL := help
