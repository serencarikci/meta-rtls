.PHONY: help up down logs backend-run frontend-run migrate deps test lint fmt

help:
	@echo "MetaRTLS targets:"
	@echo "  make up            - Start Oracle, Redis, Mosquitto"
	@echo "  make down          - Stop infrastructure"
	@echo "  make logs          - Tail compose logs"
	@echo "  make deps          - Install Go and frontend dependencies"
	@echo "  make backend-run   - Run Go API (local)"
	@echo "  make frontend-run  - Run React app (local)"
	@echo "  make fmt           - Format Go (gofmt) and React (Prettier)"
	@echo "  make test          - Run backend tests"
	@echo "  make lint          - Run golangci-lint if available"

fmt:
	cd backend && go fmt ./...
	cd frontend && npm run format

up:
	docker compose --env-file .env.example up -d

down:
	docker compose down

logs:
	docker compose logs -f

deps:
	cd backend && go mod tidy
	cd frontend && npm install

backend-run:
	cd backend && go run ./cmd/api

frontend-run:
	cd frontend && npm run dev

test:
	cd backend && go test ./...

lint:
	cd backend && golangci-lint run ./... || true
