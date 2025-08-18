package threadpool

import (
	"log"
	"net"
	"redis_go/internal/config"
)

type MultiThreadServer struct {
	listenAddr  string
	listener    net.Listener
	quitChannel chan struct{}
	connQueue   chan net.Conn
	workerList  []*Worker
}

func NewMultiThreadServer(listenAddr string, numWorker int) *MultiThreadServer {
	return &MultiThreadServer{
		listenAddr:  listenAddr,
		quitChannel: make(chan struct{}),
		connQueue:   make(chan net.Conn),
		workerList:  make([]*Worker, numWorker),
	}
}

func (server *MultiThreadServer) addConn(conn net.Conn) {
	server.connQueue <- conn
}

func (server *MultiThreadServer) InitWorkerPool() {
	for i := 0; i < len(server.workerList); i++ {
		worker := NewWorker(i, server.connQueue)
		server.workerList[i] = worker
		go worker.Start()
	}
}

func (server *MultiThreadServer) Start() error {
	var ln, err = net.Listen(config.Protocol, server.listenAddr)
	if err != nil {
		return err
	}

	defer ln.Close()
	server.listener = ln

	server.acceptConnection()

	<-server.quitChannel

	return nil
}

func (server *MultiThreadServer) acceptConnection() {
	for {
		var conn, err = server.listener.Accept()
		if err != nil {
			log.Println("Accept error", err)
		}
		log.Println("New client connect to MultiThreadServer with IP: ", conn.RemoteAddr().String())
		server.addConn(conn)
	}
}
