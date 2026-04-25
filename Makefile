.PHONY: build run test lint clean serve serve-dev serve-prod

BINARY_NAME=git-resume
BUILD_DIR=bin
PORT?=8080
HOST?=localhost

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

run:
	go run . $(ARGS)

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -rf $(BUILD_DIR)
	rm -f data/*.db

# Development
dev:
	go run . analyze --dry-run

# Install dependencies
deps:
	go mod download
	go mod tidy

# Create required directories
setup:
	mkdir -p $(BUILD_DIR) data output

# =====================
# Server Commands
# =====================

# Start server (default: localhost:8080)
serve: build
	./$(BUILD_DIR)/$(BINARY_NAME) serve --host=$(HOST) --port=$(PORT)

# Start server without rebuilding
serve-quick:
	./$(BUILD_DIR)/$(BINARY_NAME) serve --host=$(HOST) --port=$(PORT)

# Start server with hot reload (requires air: go install github.com/air-verse/air@latest)
serve-dev:
	air -c .air.toml 2>/dev/null || (echo "Installing air..." && go install github.com/air-verse/air@latest && air -c .air.toml)

# Start server for production (external access)
serve-prod: build
	./$(BUILD_DIR)/$(BINARY_NAME) serve --host=0.0.0.0 --port=$(PORT)

# Start server in background
serve-bg: build
	@echo "Starting server in background on $(HOST):$(PORT)..."
	@nohup ./$(BUILD_DIR)/$(BINARY_NAME) serve --host=$(HOST) --port=$(PORT) > logs/server.log 2>&1 & echo $$! > .server.pid
	@echo "Server PID: $$(cat .server.pid)"
	@echo "Logs: logs/server.log"

# Stop background server
serve-stop:
	@if [ -f .server.pid ]; then \
		kill $$(cat .server.pid) 2>/dev/null && echo "Server stopped" || echo "Server not running"; \
		rm -f .server.pid; \
	else \
		echo "No server PID file found"; \
	fi

# Check server status
serve-status:
	@if [ -f .server.pid ] && kill -0 $$(cat .server.pid) 2>/dev/null; then \
		echo "Server running (PID: $$(cat .server.pid))"; \
	else \
		echo "Server not running"; \
		rm -f .server.pid 2>/dev/null; \
	fi

# View server logs
serve-logs:
	@tail -f logs/server.log 2>/dev/null || echo "No logs found"

# =====================
# TUI Commands
# =====================

# Start interactive TUI
tui: build
	./$(BUILD_DIR)/$(BINARY_NAME) tui

# =====================
# Analysis Commands
# =====================

# Analyze current month
analyze: build
	./$(BUILD_DIR)/$(BINARY_NAME) analyze

# Dry run analysis
analyze-dry: build
	./$(BUILD_DIR)/$(BINARY_NAME) analyze --dry-run

# =====================
# Help
# =====================

help:
	@echo "Git Resume Analyzer - Available Commands"
	@echo ""
	@echo "Build & Run:"
	@echo "  make build        - Build binary"
	@echo "  make run ARGS=... - Run with arguments"
	@echo "  make test         - Run tests"
	@echo "  make lint         - Run linter"
	@echo "  make clean        - Clean build artifacts"
	@echo ""
	@echo "Server:"
	@echo "  make serve            - Build and start server (localhost:8080)"
	@echo "  make serve PORT=3000  - Start on custom port"
	@echo "  make serve-quick      - Start without rebuilding"
	@echo "  make serve-dev        - Start with hot reload"
	@echo "  make serve-prod       - Start for external access (0.0.0.0)"
	@echo "  make serve-bg         - Start in background"
	@echo "  make serve-stop       - Stop background server"
	@echo "  make serve-status     - Check server status"
	@echo "  make serve-logs       - View server logs"
	@echo ""
	@echo "Analysis:"
	@echo "  make tui              - Start interactive TUI"
	@echo "  make analyze          - Analyze current month"
	@echo "  make analyze-dry      - Dry run analysis"
