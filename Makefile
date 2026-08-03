.PHONY: help build run clean test docker-build docker-run download-db cli mcp compare

help: ## Show this help message
@echo 'Usage: make [target]'
@echo ''
@echo 'Available targets:'
@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build API and tools
@echo "Building binaries..."
cd cmd/geoip-api && go build -o ../../geoip-api .
cd cmd/geoip-cli && go build -o ../../geoip-cli .
cd cmd/geoip-mcp && go build -o ../../geoip-mcp .
cd cmd/geoip-compare && go build -o ../../geoip-compare .
@echo "Build complete"

run: ## Run API server
@echo "Starting GeoIP API server..."
cd cmd/geoip-api && go run .

cli: ## Run CLI help
cd cmd/geoip-cli && go run . -h

mcp: ## Run MCP server (stdio)
cd cmd/geoip-mcp && go run .

compare: ## Run compare tool help
cd cmd/geoip-compare && go run . -h

clean: ## Clean build artifacts
@echo "Cleaning build artifacts..."
rm -f geoip-api geoip-cli geoip-mcp geoip-compare
rm -f cmd/geoip-api/geoip-api cmd/geoip-cli/geoip-cli cmd/geoip-mcp/geoip-mcp cmd/geoip-compare/geoip-compare
@echo "Clean complete"

test: ## Run tests
@echo "Running tests..."
go test -v ./...

download-db: ## Download DB-IP database (no license required)
@echo "Downloading DB-IP database..."
./scripts/download_dbip.sh all
@echo "Database downloaded successfully"

docker-build: ## Build Docker image
@echo "Building Docker image..."
docker build -t geoip-api:latest .
@echo "Docker image built successfully"

docker-run: ## Run Docker container
@echo "Running Docker container..."
docker-compose up -d
@echo "Docker container started. Access at http://localhost:8080"

docker-stop: ## Stop Docker container
@echo "Stopping Docker container..."
docker-compose down
@echo "Docker container stopped"

dev: ## Run in development mode
@echo "Starting development server..."
cd cmd/geoip-api && GIN_MODE=debug go run .
