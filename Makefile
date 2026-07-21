.PHONY: help up down logs backend-run frontend-run deps test lint fmt ready

help:
	@echo "MetaRTLS targets:"
	@echo "  make up            - Start Oracle and Mosquitto"
	@echo "  make down          - Stop infrastructure"
	@echo "  make logs          - Tail compose logs"
	@echo "  make deps          - Install Go and frontend dependencies"
	@echo "  make backend-run   - Run Go API (local)"
	@echo "  make frontend-run  - Run React app (local)"
	@echo "  make ready         - Check API /ready (Oracle ping)"
	@echo "  make fmt           - Format Go (gofmt) and React (Prettier)"
	@echo "  make test          - Run backend tests"
	@echo "  make lint          - Run go vet + gofmt check"

fmt:
	cd backend && go fmt ./...
	cd frontend && npm run format

up:
	docker compose up -d

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

ready:
	curl -sf http://localhost:8090/ready | python3 -m json.tool

test:
	cd backend && go test ./...

lint:
	cd backend && test -z "$$(gofmt -l .)" && go vet ./...
