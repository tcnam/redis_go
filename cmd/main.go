package main

import (
	"redis_go/internal/config"
	"redis_go/internal/server/iomultiplexer"
)

func main() {
	server := iomultiplexer.NewIOMultiplexerServer(config.Address, config.Protocol)
	server.Start()
	// server := threadpool.NewMultiThreadServer(config.Address, config.PoolSize)
	// server.InitWorkerPool()
	// server.Start()
}
