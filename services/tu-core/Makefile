.PHONY: proto build run test clean

# Generate protobuf code
# proto:
# 	@echo "🔨 Generating protobuf code..."
# 	protoc \
# 		--go_out=. --go_opt=paths=source_relative \
# 		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
# 		proto/student/v1/student.proto
# 	@echo "✅ Protobuf code generated"
# Find all proto files
PROTO_FILES := $(shell find proto -name "*.proto")

# Generate protobuf code for all proto files
proto:
	@echo "🔨 Generating protobuf code..."
	protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	@echo "✅ Protobuf code generated"


# Build binary
build:
	@echo "🔨 Building tu-core..."
	go build -o bin/tu-core ./cmd/main.go
	@echo "✅ Build complete"

# Run service
run:
	go run ./cmd/main.go

# Run tests
test:
	go test ./... -v

# Run tests with coverage
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	go mod tidy
	go mod download

# Lint
lint:
	golangci-lint run

# All
all: deps proto build test



# ===========================================================================
# DATABASE MIGRATION
# ===========================================================================

# Database connection string (Mempertahankan format variabel Anda dengan tanda minus)
DB_USER ?= -smk_dev
DB_PASSWORD ?= -secret123
DB_NAME ?= -smk_platform
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Migrate up (apply all migrations)
# Migrate up (apply all migrations)
migrate-up:
	@echo "🗄️  Running database migrations (up)..."
	@echo "   DB URL: $(DB_URL)"
	@if ! command -v migrate >/dev/null 2>&1; then \
		echo "❌ Tool 'migrate' tidak ditemukan. Install dengan:"; \
		echo "   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; \
		exit 1; \
	fi
	@# PERBAIKAN: Gunakan 'file://migrations' agar path terbaca dengan benar di Windows/Git Bash
	@migrate -path migrations -database "$(DB_URL)" up
	@echo "✅ Migrations complete"

# Migrate down (rollback last N steps)
migrate-down:
	@echo "⏪ Rolling back database migrations..."
	@read -p "How many steps to rollback? (default: 1): " steps; \
	steps=$${steps:-1}; \
	migrate -path migrations -database "$(DB_URL)" down $$steps
	@echo "✅ Rollback complete"

# Migrate status (check current version)
migrate-status:
	@echo "📊 Migration Status:"
	@migrate -path migrations -database "$(DB_URL)" version

# Create new migration
migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name
	@echo "✅ Migration created"
