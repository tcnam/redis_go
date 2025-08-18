package main

import (
	"redis_go/internal/config"
	"redis_go/internal/server/threadpool"
)

func main() {
	server := threadpool.NewMultiThreadServer(config.Address, config.PoolSize)
	server.InitWorkerPool()
	server.Start()
}
