run:
	go run cmd/server/main.go

test:
	go test ./...

build:
	go build -o bin/gateway cmd/server/main.go

tidy:
	go mod tidy