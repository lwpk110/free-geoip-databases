.PHONY: help build run clean test docker-build docker-run download-db

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the GeoIP API server
	@echo "Building GeoIP API server..."
	cd cmd/geoip-api && go build -o ../../geoip-api .
	@echo "Build complete: ./geoip-api"

run: build ## Build and run the server
	@echo "Starting GeoIP API server..."
	./geoip-api

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -f geoip-api
	rm -f cmd/geoip-api/geoip-api
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

install: build ## Install the binary to /usr/local/bin
	@echo "Installing geoip-api to /usr/local/bin..."
	sudo cp geoip-api /usr/local/bin/
	@echo "Installation complete"

dev: ## Run in development mode
	@echo "Starting development server..."
	cd cmd/geoip-api && GIN_MODE=debug go run main.go
