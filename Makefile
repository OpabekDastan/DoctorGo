run:
	go run ./cmd/api

build:
	go build -o bin/doctor_go ./cmd/api

test:
	go test ./...
