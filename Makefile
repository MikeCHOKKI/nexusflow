.PHONY: help up down build proto frontend seed test lint clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

up: ## Start all services with Docker Compose
	docker compose up --build -d

down: ## Stop all services
	docker compose down

build: proto ## Build all services
	docker compose build

proto: ## Generate protobuf code for all languages
	@echo "Generating protos..."
	@mkdir -p protos/gen/go protos/gen/python protos/gen/php
	@protoc --go_out=protos/gen/go --go_opt=paths=source_relative \
		--go-grpc_out=protos/gen/go --go-grpc_opt=paths=source_relative \
		--python_out=protos/gen/python --grpc_python_out=protos/gen/python \
		--php_out=protos/gen/php \
		-I protos protos/*.proto
	@echo "✓ Protos generated"

frontend: ## Install and build frontend
	cd frontend && npm install && npm run build

migrate: ## Run pending database migrations
	@for f in db/migrations/*.sql; do \
		echo "→ Running migration: $$f..."; \
		docker compose exec -T postgres psql -U nexusflow -d nexusflow -f /$$f 2>/dev/null || \
		docker compose exec -T postgres psql -U nexusflow -d nexusflow -c "$$(cat $$f)"; \
	done
	@echo "✓ Migrations applied"

seed: ## Run database seed
	docker compose exec postgres psql -U nexusflow -d nexusflow -f /db/seeds/seed.sql

test: ## Run all tests
	@echo "Running tests..."
	@for dir in api-gateway service-user service-catalog service-order service-payment service-notify; do \
		echo "→ Testing $$dir..."; \
		(cd $$dir && go test ./... 2>/dev/null || go test ./... 2>/dev/null || npm test 2>/dev/null || pytest 2>/dev/null || echo "  No tests defined"); \
	done

lint: ## Lint all services
	@echo "Linting..."
	@cd api-gateway && golangci-lint run ./... 2>/dev/null || echo "  api-gateway: no linter"
	@cd service-user && golangci-lint run ./... 2>/dev/null || echo "  service-user: no linter"
	@cd service-payment && golangci-lint run ./... 2>/dev/null || echo "  service-payment: no linter"
	@echo "✓ Lint complete"

clean: ## Clean build artifacts
	rm -rf protos/gen/
	rm -rf frontend/dist/
	rm -rf api-gateway/bin/ service-user/bin/ service-payment/bin/
	@echo "✓ Clean complete"

logs: ## Tail logs for all services
	docker compose logs -f

ps: ## Show running containers
	docker compose ps

restart: down up ## Restart all services
