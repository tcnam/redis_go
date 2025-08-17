build:
	go build -o ./bin/cmd/main ./cmd/main.go

run:
	go build -o ./bin/cmd/main ./cmd/main.go
	./bin/cmd/main

run_pool:
	go build -o ./bin/exec/threadpool/main ./exec/threadpool/main.go
	./bin/exec/threadpool/main

run_asyncio:
	go build -o ./bin/exec/asyncio/main ./exec/asyncio/main.go
	./bin/exec/asyncio/main

test:
	go test -v ./..