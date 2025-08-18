build:
	go build -o ./bin/main ./cmd/main.go

run:
	go build -o ./bin/main ./cmd/main.go
	./bin/main

run_pool:
	go build -o ./bin/threadpool/main ./threadpool/main.go
	./bin/threadpool/main

test:
	go test -v ./..