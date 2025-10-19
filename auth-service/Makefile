# -------------------------------
# Environment variables
# -------------------------------
APP_NAME := auth-service
BINARY := $(APP_NAME)
MAIN := ./cmd/main.go
MIGRATION_DIR := ./migrations
PROTO_SRC := ./internal/delivery/grpc/proto_gen
GO_PACKAGES := ./...

# -------- Load .env --------
ifneq (,$(wildcard ./.env))
	include .env
	export
endif

# -------------------------------
# Basic commands
# -------------------------------

.PHONY: run
run: 
	@echo "Running $(APP_NAME)..."
	@go run $(MAIN)

.PHONY: build
build:
	@echo "Building binary..."
	@go build -o $(BINARY) $(MAIN)
	@echo "Build complete: ./$(BINARY)"

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY)
	@echo "Cleaned."

# -------------------------------
# Docker commands
# -------------------------------
.PHONY: docker-up
docker-up:
	@echo "Starting Docker Compose..."
	@docker compose up -d --build

.PHONY: docker-down
docker-down:
	@echo "Stopping Docker Compose..."
	@docker compose down

# -------------------------------
# Migration commands
# -------------------------------
.PHONY: migrate
migrate:
	@echo "Running database migrations..."
	@migrate -path $(MIGRATION_DIR) \
		-database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)" up
	@echo "Migration completed!"

.PHONY: migrate-down
migrate-down:
	@echo "Rolling back last migration..."
	@migrate -path $(MIGRATION_DIR) \
		-database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)" down 1
	@echo "Rolled back one migration."

.PHONY: migrate-create
migrate-create: ## Create a new migration (usage: make migrate-create name=create_user_table)
	@if [ -z "$(name)" ]; then \
		@echo "Please provide a migration name. Example: make migrate-create name=add_refresh_tokens"; \
		exit 1; \
	@fi
	@migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(name)
	@echo "Created new migration: $(name)"

# ==============================================
# Proto Generation (gRPC)
# ==============================================
.PHONY: proto_gen
proto_gen: ## Generate gRPC code from proto files
	@echo "Generating gRPC stubs..."
	@proto --go_out=. --go-grpc_out=. $(PROTO_SRC)/*.proto
	@echo "Proto files generated successfully!"

# -------------------------------
# Testing / Linting / Formatting
# -------------------------------
.PHONY: test
test:
	@echo "Running tests..."
	@go test $(GO_PACKAGES) -v
	@echo "Tests completed."

.PHONY: lint
lint: ## Run static code analysis
	@echo "Linting Go code..."
	@golangci-lint run --timeout 2m
	@echo "Lint passed!"

.PHONY: format
format: ## Format codebase with gofmt
	@echo "Formatting code..."
	@gofmt -s -w .
	@echo "Code formatting."
