package threadpool

import (
	"fmt"
	"io"
	"log"
	"net"
)

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
func (worker *Worker) Start() error {
	for conn := range worker.connQueue {
		log.Printf("Worker id %d serve conn %s", worker.id, conn.RemoteAddr().String())
		worker.handleConnection(conn)
	}
	return nil
}

func (worker *Worker) handleConnection(conn net.Conn) {
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

		log.Printf("MultiThreadServer received from client: %s, with message: %s", message.clientAddr, string(message.payload))

		msg := fmt.Sprintf("+Received from %s: %s\r\n", conn.RemoteAddr().String(), string(buffer[:numBytes]))
		resp := fmt.Sprintf("$%d\r\n%s\r\n", len(msg), msg)
		_, err = conn.Write([]byte(resp))

		if err != nil {
			log.Printf("Error writing to client %s: %v\n", conn.RemoteAddr().String(), err)
			return
		}
	}
}
