.PHONY: help run build test clean docker-up docker-down migrate-up migrate-down

# Variables
APP_NAME=gobank
MAIN_PATH=cmd/api/main.go
DOCKER_COMPOSE=docker-compose
MIGRATE_PATH=db/migrations

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

run: ## Run the application locally
	go run $(MAIN_PATH)

build: ## Build the application
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -cover ./...

clean: ## Clean build artifacts
	rm -rf bin/

docker-up: ## Start all services with docker-compose
	$(DOCKER_COMPOSE) up -d

docker-down: ## Stop all services
	$(DOCKER_COMPOSE) down

docker-build: ## Build docker image
	$(DOCKER_COMPOSE) build

docker-logs: ## View logs
	$(DOCKER_COMPOSE) logs -f

migrate-up: ## Run database migrations
	migrate -path $(MIGRATE_PATH) -database "postgresql://postgres:postgres@localhost:5432/gobank?sslmode=disable" up

migrate-down: ## Rollback database migrations
	migrate -path $(MIGRATE_PATH) -database "postgresql://postgres:postgres@localhost:5432/gobank?sslmode=disable" down

migrate-create: ## Create new migration (usage: make migrate-create name=create_table)
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

install-tools: ## Install development tools
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

mod-tidy: ## Tidy go modules
	go mod tidy

mod-vendor: ## Vendor dependencies
	go mod vendor

fmt: ## Format code
	go fmt ./...

lint: ## Run linter
	golangci-lint run

dev: ## Run with hot reload (requires air)
	air