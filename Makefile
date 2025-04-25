.PHONY: build run clean test help

# Default target
.DEFAULT_GOAL := help

# Variables
APP_NAME := schemaless-api
BUILD_DIR := ./build
COVERAGE_REPORT_DIR := ./coverage

# Help target
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build     Build the application"
	@echo "  run       Run the application"
	@echo "  clean     Clean build artifacts"
	@echo "  test      Run tests"
	@echo "  help      Show this help message"

# Build target
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME)

# Run target
run:
	@echo "Running $(APP_NAME)..."
	@go run main.go

# Clean target
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(COVERAGE_REPORT_DIR)

# Test target
test:
	@echo "Running tests..."
	@go test -v ./...

# Test coverage
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_REPORT_DIR)
	@go test -v -coverprofile=$(COVERAGE_REPORT_DIR)/coverage.out ./...
	@go tool cover -html=$(COVERAGE_REPORT_DIR)/coverage.out -o $(COVERAGE_REPORT_DIR)/coverage.html
	@echo "Coverage report generated: $(COVERAGE_REPORT_DIR)/coverage.html"
