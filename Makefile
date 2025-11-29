PROJECT_NAME := simple-cli
PROJECT_ROOT := $(shell pwd)
GO := go
GOFMT := gofmt
GOIMPORTS := goimports
WIRE := wire
BUILD_DIR := build

.PHONY: all
all: tidy wire imports build

.PHONY: tidy
tidy:
	@echo "Tidying up dependencies..."
	@$(GO) mod tidy

.PHONY: wire
wire:
	@echo "Generating Wire dependency injection code..."
	@cd $(PROJECT_ROOT)/internal/handler/order && $(WIRE)

.PHONY: imports
imports:
	@echo "Formatting imports..."
	@find $(PROJECT_ROOT) -name "*.go" -exec $(GOIMPORTS) -w {} \;

.PHONY: build
build:
	@echo "Building the project..."
	@mkdir -p $(BUILD_DIR)/
	@$(GO) build -o $(BUILD_DIR)/$(PROJECT_NAME) $(PROJECT_ROOT)/cmd/main.go

.PHONY: run
run:
	@echo "Running the project..."
	@$(BUILD_DIR)/$(PROJECT_NAME)

.PHONY: test
test:
	@echo "Running the tests..."
	@$(GO) test $(PROJECT_ROOT)/...

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)/

.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  all - Run tidy, wire, imports, and build"
	@echo "  tidy - Run go mod tidy"
	@echo "  wire - Generate Wire dependency injection code"
	@echo "  imports - Format imports"
	@echo "  build - Build the project"
	@echo "  run - Run the project"
	@echo "  test - Run the tests"
	@echo "  clean - Clean the project"
	@echo "  help - Show this help message"