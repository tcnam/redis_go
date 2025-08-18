package iomultiplexing

import (
	"net"
)

type IOMultiplexingServer struct {
	listenAddr  string
	listener    net.Listener
	quitChannel chan struct{}
}
