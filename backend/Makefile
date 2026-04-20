.PHONY: run build test migrate-up migrate-down

run:
	go run cmd/server/main.go

build:
	go build -o bin/skroda cmd/server/main.go

test:
	go test ./... -v

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1
