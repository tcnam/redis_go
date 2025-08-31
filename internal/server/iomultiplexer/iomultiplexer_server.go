package iomultiplexer

import (
	"fmt"
	"io"
	"log"
	"net"
	"redis_go/internal/config"
	"redis_go/internal/core/io_multiplexing"
	"syscall"
)

type IOMultiplexerServer struct {
	listenAddr  string
	protocol    string
	listener    net.Listener
	quitChannel chan struct{}
}

func NewIOMultiplexerServer(listenAddr string, protocol string) *IOMultiplexerServer {
	return &IOMultiplexerServer{
		listenAddr:  listenAddr,
		protocol:    protocol,
		quitChannel: make(chan struct{}),
	}
}

func (server *IOMultiplexerServer) Start() error {
	listener, err := net.Listen(config.Protocol, server.listenAddr)
	if err != nil {
		return err
	}

	defer listener.Close()

	server.listener = listener

	// Get the file descriptor from the listener
	tcpListerner, ok := server.listener.(*net.TCPListener)
	if !ok {
		log.Fatal("Listener is not a TCPListener")
	}

	listenerFile, err := tcpListerner.File()
	if err != nil {
		log.Fatal(err)
	}
	defer listenerFile.Close()

	serverFd := uint64(listenerFile.Fd())

	ioMuliplexer, err := io_multiplexing.CreateIOMultiplexer()
	if err != nil {
		log.Fatal(err)
	}
	defer ioMuliplexer.Close()

	err = ioMuliplexer.Monitor(
		io_multiplexing.NewEvent(serverFd, io_multiplexing.OpRead),
	)

	if err != nil {
		log.Fatal(err)
	}

	var events []io_multiplexing.Event = make([]io_multiplexing.Event, config.MaxConnection)
	// var lastActiveExpireExecTime = time.Now()

	for {
		// if time.Now().After(lastActiveExpireExecTime.Add(constant.ActiveExpireFrequency)) {
		// 	io_multiplexing
		// }
		events, err = ioMuliplexer.Wait()
		if err != nil {
			continue
		}

		for i := 0; i < len(events); i++ {
			if events[i].Fd == serverFd {
				log.Printf("New client is trying to connect")
				connFd, _, err := syscall.Accept(int(serverFd))
				if err != nil {
					log.Println("err: ", err)
				}
				log.Printf("Setup new connection")
				err = ioMuliplexer.Monitor(
					io_multiplexing.NewEvent(uint64(connFd), io_multiplexing.OpRead),
				)

				if err != nil {
					log.Fatal(err)
				}
			} else {
				err := server.readCommand(int(events[i].Fd))
				if err != nil {
					if err == io.EOF || err == syscall.ECONNRESET {
						log.Println("client disconnected")
						_ = syscall.Close(int(events[i].Fd))
						continue
					}
					log.Println("read error:", err)
					continue
				}
			}

		}
	}

	// Create an ioMultiplexer

	<-server.quitChannel

	return nil
}

func (server *IOMultiplexerServer) readCommand(fd int) error {
	var buf = make([]byte, 512)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		return err
	}

	if n == 0 {
		return io.EOF
	}
	msg := fmt.Sprintf("+Received from %s: %s\r\n", string(buf[:n]))
	resp := fmt.Sprintf("$%d\r\n%s\r\n", len(msg), msg)
	syscall.Write(fd, []byte(resp))
	return nil
}
