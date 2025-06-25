.PHONY: help build run test clean docker-build docker-up docker-down docker-logs migrate-up migrate-down rebuild quick-rebuild

help:
	@echo "Available commands:"
	@echo "  make build          - Build the Go application"
	@echo "  make run            - Run the application locally"
	@echo "  make test           - Run tests"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-up      - Start all services with docker-compose"
	@echo "  make docker-down    - Stop all services"
	@echo "  make docker-logs    - View logs from all services"
	@echo "  make migrate-up     - Run database migrations"
	@echo "  make migrate-down   - Rollback database migrations"
	@echo "  make rebuild        - Full rebuild (stops, rebuilds, and restarts all services)"
	@echo "  make quick-rebuild  - Quick rebuild (only backend, keeps database)"

build:
	go build -o bin/server cmd/server/main.go

run:
	go run cmd/server/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

migrate-up:
	docker compose run --rm migrate

migrate-down:
	docker compose run --rm migrate -path /migrations -database "postgres://pixelgame:pixelgame123@postgres:5432/pixelgame_db?sslmode=disable" down 1

dev: docker-up migrate-up
	@echo "Development environment is ready!"
	@echo "Backend server: http://localhost:8080"
	@echo "PostgreSQL: localhost:5432"
	@echo "Redis: localhost:6379"

rebuild:
	@./scripts/rebuild.sh

quick-rebuild:
	@./scripts/quick-rebuild.sh