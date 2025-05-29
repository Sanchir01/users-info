PHONY:
SILENT:
MIGRATION_NAME ?= new_migration

DB_CONN_PROD = "host=localhost user=postgres password=postgres port=5439 dbname=postgres sslmode=disable"

swag:
	swag init -g cmd/main/main.go

build:
	go build -o ./.bin/main ./cmd/main/main.go
run: build
	ENV_FILE=".env.prod" ./.bin/main

migrations-up:
	goose -dir migrations postgres "host=localhost user=postgres password=postgres port=5439 dbname=test sslmode=disable"  up

migrations-down:
	goose -dir migrations postgres  "host=localhost user=postgres password=postgres port=5439 dbname=test sslmode=disable"  down


migrations-status:
	goose -dir migrations postgres  "host=localhost user=postgres password=postgres port=5439 dbname=test sslmode=disable" status

migrations-new:
	goose -dir migrations create $(MIGRATION_NAME) sql

migrations-up-prod:
	goose -dir migrations postgres $(DB_CONN_PROD) up

migrations-down-prod:
	goose -dir migrations postgres $(DB_CONN_PROD) down

migrations-status-prod:
	goose -dir migrations postgres $(DB_CONN_PROD) status

docker-build:
	docker build -t users-info .

docker:
	docker-compose  up -d

docker-app: docker-build docker

compose-prod:
	docker compose -f docker-compose.prod.yaml up --build -d
