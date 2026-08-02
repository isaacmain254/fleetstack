# .PHONY: dev backend web-app

dev:
	docker compose -f docker-compose.yml up --build

run:
	@echo "Starting backend server..."
	cd backend && go run ./cmd/api/main.go

build:
	@echo "Building backend server..."
	cd backend && go build -o ./bin/main ./cmd/api/main.go

migrate-up:
	@echo "Running database migrations up..."
	cd backend && go run ./cmd/migrate/main.go up

migrate-down:
	@echo "Running database migrations down..."
	cd backend && go run ./cmd/migrate/main.go down

migrate-create:
	@echo "Creating new database migration..."
	cd backend && go run ./cmd/migrate/main.go create $(name)

