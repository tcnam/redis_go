package asyncio

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

type Server struct {
	listenAddr     string
	listener       net.Listener
	quitChannel    chan struct{}
	messageChannel chan ClientMessage
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr:     listenAddr,
		quitChannel:    make(chan struct{}),
		messageChannel: make(chan ClientMessage),
	}
}

func (s *Server) start() error {
	var ln, err = net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}

	defer ln.Close()
	s.listener = ln

	go s.acceptLoop()

	<-s.quitChannel
	close(s.messageChannel)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		var conn, err = s.listener.Accept()
		if err != nil {
			log.Println("Accept error", err)
		}
		log.Println("new client connect to server with IP: ", conn.RemoteAddr().String())
		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
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
		var message = ClientMessage{
			clientAddr: conn.RemoteAddr().String(),
			payload:    buffer[:numBytes],
		}

		s.messageChannel <- message
		_, err = conn.Write([]byte(fmt.Sprintf("Received from %s: %s\n", conn.RemoteAddr().String(), string(buffer[:numBytes]))))
		if err != nil {
			log.Printf("Error writing to client %s: %v\n", conn.RemoteAddr().String(), err)
			return
		}

	}
}

func main() {
	var server = NewServer("localhost:3000")
	go func() {
		for message := range server.messageChannel {
			log.Println(message.clientAddr, string(message.payload))
		}
	}()
	server.start()
}
