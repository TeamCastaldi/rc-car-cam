.PHONY: dev-backend dev-frontend build test lint

dev-backend:
	cd backend && go run ./cmd/server

dev-frontend:
	cd frontend && npm run dev -- --open

build:
	cd backend && go build -o bin/server ./cmd/server
	cd frontend && npm run build

test:
	cd backend && go test ./...
	cd frontend && npm test

lint:
	cd backend && go vet ./... && golangci-lint run
	cd frontend && npm run lint
