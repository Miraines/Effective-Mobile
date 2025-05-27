.PHONY: build run up down migrate

GOFLAGS = -mod=vendor

build:
	go build -o bin/people ./cmd/server

run: build
	./bin/people

up:
	docker compose -f build/docker-compose.yml up -d --build

down:
	docker compose -f build/docker-compose.yml down

migrate:
	migrate -path migrations -database "postgres://people:people@localhost:5433/people_enricher?sslmode=disable" up

run-local:
	go run ./cmd/server
