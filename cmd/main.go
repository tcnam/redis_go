package main

import "redis_go/internal/server/threadpool"

func main() {
	server := threadpool.NewServer("localhost:3000", 200)
	server.InitWorkerPool()
	server.Start()
}
