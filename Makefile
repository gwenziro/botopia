.PHONY: build run test clean dev install-dev tidy

# Nama aplikasi dan path
APP_NAME=botopia
BUILD_DIR=./build
MAIN_PATH=./cmd/main.go

# Go commands
GO=go
GOBUILD=$(GO) build
GOTEST=$(GO) test
GOCLEAN=$(GO) clean
GOMOD=$(GO) mod
GOTIDY=$(GOMOD) tidy

# Build aplikasi
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

# Run aplikasi
run:
	@echo "Running $(APP_NAME)..."
	$(GO) run $(MAIN_PATH)

# Test aplikasi
test:
	@echo "Testing $(APP_NAME)..."
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -rf tmp
	@find . -name "*.db" -not -path "*/\.*" -delete
	@find . -name "*.db-*" -not -path "*/\.*" -delete

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GOTIDY)

# Development mode dengan hot-reload
dev:
	@echo "Running in development mode with hot-reload..."
	@air

# Install development dependencies
install-dev:
	@echo "Installing development dependencies..."
	$(GO) install github.com/air-verse/air@latest
	$(GOTIDY)

# Default target
default: build
