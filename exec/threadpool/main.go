package threadpool

import (
	"fmt"
	"io"
	"log"
	"net"
)

type ClientMessage struct {
	clientAddr string
	payload    []byte
}

// Worker represent a thread in thread pool
type Worker struct {
	id        int
	connQueue chan net.Conn
}

// Worker constructor
func NewWorker(id int, connQueue chan net.Conn) *Worker {
	return &Worker{
		id:        id,
		connQueue: connQueue,
	}
}

// Start a worker
func (worker *Worker) start() error {
	for conn := range worker.connQueue {
		log.Printf("Worker id %d serve conn %s", worker.id, conn.RemoteAddr().String())
		handleConnection(conn)
	}
	return nil
}

type Server struct {
	listenAddr  string
	listener    net.Listener
	quitChannel chan struct{}
	connQueue   chan net.Conn
	workerList  []*Worker
}

func NewServer(listenAddr string, numWorker int) *Server {
	return &Server{
		listenAddr:  listenAddr,
		quitChannel: make(chan struct{}),
		connQueue:   make(chan net.Conn),
		workerList:  make([]*Worker, numWorker),
	}
}

func (server *Server) addConn(conn net.Conn) {
	server.connQueue <- conn
}

func (server *Server) initWorkerPool() {
	for i := 0; i < len(server.workerList); i++ {
		worker := NewWorker(i, server.connQueue)
		server.workerList[i] = worker
		go worker.start()
	}
}

func (server *Server) start() error {
	var ln, err = net.Listen("tcp", server.listenAddr)
	if err != nil {
		return err
	}

	defer ln.Close()
	server.listener = ln

	server.acceptConnection()

	<-server.quitChannel

	return nil
}

func (server *Server) acceptConnection() {
	for {
		var conn, err = server.listener.Accept()
		if err != nil {
			log.Println("Accept error", err)
		}
		log.Println("New client connect to server with IP: ", conn.RemoteAddr().String())
		server.addConn(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	var buffer = make([]byte, 2048)
	for {
		numBytes, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Printf("Client disconnected (EOF): %s\n", conn.RemoteAddr().String())
			} else {
				log.Printf("Error reading from client %s: %v\n", conn.RemoteAddr().String(), err)
			}
			return // Exit the loop and handleConnection function
		}
		message := ClientMessage{
			clientAddr: conn.RemoteAddr().String(),
			payload:    buffer[:numBytes],
		}

		log.Printf("Server received from client: %s, with message: %s", message.clientAddr, string(message.payload))

		_, err = conn.Write([]byte(fmt.Sprintf("Received from %s: %s\n", conn.RemoteAddr().String(), string(buffer[:numBytes]))))
		if err != nil {
			log.Printf("Error writing to client %s: %v\n", conn.RemoteAddr().String(), err)
			return
		}
	}
}

func main() {
	server := NewServer("localhost:3000", 2)
	server.initWorkerPool()
	server.start()
}
