BINARY=postpanda-backend
CMD=./cmd/server

.PHONY: build run dev test tidy

build:
	go build -o $(BINARY) $(CMD)

run:
	go run $(CMD)/main.go

dev:
	DATABASE_URL=postgresql://postgres:postgres@localhost:5432/postpanda \
	REDIS_URL=redis://localhost:6379 \
	JWT_SECRET=dev-secret \
	FRONTEND_URL=http://localhost:4200 \
	go run $(CMD)/main.go

test:
	go test ./...

tidy:
	go mod tidy

lint:
	golangci-lint run ./...
