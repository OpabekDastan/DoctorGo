run:
	go run ./cmd/api

build:
	go build -o bin/doctorgo ./cmd/api

test:
	go test ./...
