# Variables
APP_NAME=devboard
MAIN_PATH=./cmd/api
MIGRATE_PATH=./migrations
DB_URL=postgresql://postgres:password@localhost:5432/devboard?sslmode=disable
REDIS_ADDR=localhost:6379

.PHONY: run build test lint migrate-create migrate-up migrate-down generate tidy help docs

## run: Run app
run:
	GO_ENV=development DATABASE_URL=${DB_URL} REDIS_ADDR=${REDIS_ADDR} go run $(MAIN_PATH)/main.go

## build: Compile binary
build:
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

## test: Run all tests wth race detector
test:
	go test -race -cover ./...

## lint: Analyze code with golangci-lint
lint:
	golangci-lint run ./...

## migrate-create: Create a new migration. Example use: make migrate-create name=create_users
migrate-create:
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

## migrate-up: Apply all pending migrations
migrate-up:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

## migrate-down: Revert last migration
migrate-down:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down 1

## generate: Run go generate in the entire project
generate:
	sqlc generate
	go generate ./...

## tidy: Clean and verify dependencies
tidy:
	go mod tidy
	go mod verify

## help: Show this menu
help:
	@grep -E '^##' Makefile | sed 's/## //'

## docker-up: Start development services
docker-up:
	docker compose -f docker-compose.dev.yml up -d

## docker-down: Stop development services
docker-down: 
	docker compose -f docker-compose.dev.yml down

## docker-logs: Watch logs from all services
docker-logs:
	docker compose logs -f

## docs: Generate OpenAPI docs with swaggo
docs:
	swag init -g cmd/api/main.go -o docs