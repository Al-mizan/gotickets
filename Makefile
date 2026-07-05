.PHONY: all build run migrate-install migrate-up migrate-down migrate-create

# Load environment variables
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Default database URL for migrations
DB_URL ?= "postgres://postgres:postgres@localhost:5432/gotickets?sslmode=disable"

all: build

build:
	@echo "Building the application..."
	go build -o tmp/gotickets cmd/server.go

run:
	@echo "Running the application..."
	go run cmd/server.go

migrate-install:
	@echo "Installing golang-migrate..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir db/migrations -seq $$name

migrate-up:
	@echo "Running migrations up..."
	migrate -path db/migrations -database $(DB_URL) -verbose up

migrate-down:
	@echo "Running migrations down..."
	migrate -path db/migrations -database $(DB_URL) -verbose down

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go fmt ./...
