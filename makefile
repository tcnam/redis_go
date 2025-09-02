build:
	go build -o ./bin/main ./cmd/main.go

run:
	go build -o ./bin/main ./cmd/main.go
	./bin/main

test:
	go test -v ./..

test_protocol:
	go test -v ./internal/core/protocol