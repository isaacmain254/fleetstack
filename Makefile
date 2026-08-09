# .PHONY: dev backend web-app

dev-up:
	@echo "Starting docker containers..."
	docker compose -f docker-compose.yml up -d --build

dev-down:
	@echo "Stopping docker containers..."
	docker compose -f docker-compose.yml down

run:
	@echo "Starting backend server..."
	cd backend && go run ./cmd/api/main.go

build:
	@echo "Building backend server..."
	cd backend && go build -o ./bin/main ./cmd/api/main.go

worker:
	@echo "Starting deployment worker..."
	cd backend && go run ./cmd/worker/main.go

migrate-up:
	@echo "Running database migrations up..."
	cd backend && go run ./cmd/migrate/main.go up

migrate-down:
	@echo "Running database migrations down..."
	cd backend && go run ./cmd/migrate/main.go down

migrate-create:
	@echo "Creating new database migration..."
	cd backend && go run ./cmd/migrate/main.go create $(name)

