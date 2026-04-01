.PHONY: tidy sqlc migrate-up migrate-down swag build-api build-worker build docker-up docker-down run-api run-worker test test-coverage

tidy:
	go mod tidy

sqlc:
	sqlc generate

migrate-up:
	migrate -path db/migration -database "$${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/gox?sslmode=disable}" -verbose up

migrate-down:
	migrate -path db/migration -database "$${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/gox?sslmode=disable}" -verbose down

swag:
	swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal

build-api: sqlc
	go build -o bin/api ./cmd/api

build-worker: sqlc
	go build -o bin/worker ./cmd/worker

build: build-api build-worker

docker-up:
	docker compose up -d

docker-down:
	docker compose down

run-api: sqlc
	go run ./cmd/api -config configs/config.yaml

run-worker: sqlc
	go run ./cmd/worker -config configs/config.yaml

test: sqlc
	go test ./... -count=1

test-coverage: sqlc
	go test ./... -count=1 -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
