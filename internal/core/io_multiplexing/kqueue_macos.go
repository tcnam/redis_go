//go:build darwin || freebsd || netbsd || openbsd

package io_multiplexing

import (
	"log"
	"syscall"

	"github.com/tcnam/redis_go/internal/config"
)

type KQueue struct {
	Fd            int
	kqEvents      []syscall.Kevent_t
	genericEvents []Event
}

func CreateIOMultiplexer() (*KQueue, error) {
	kQueueFD, err := syscall.Kqueue()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &KQueue{
		Fd:            kQueueFD,
		kqEvents:      make([]syscall.Kevent_t, config.MaxConnection),
		genericEvents: make([]Event, config.MaxConnection),
	}, nil
}

func (kQueue *KQueue) Monitor(event Event) error {
	kqEvent := event.toNative(syscall.EV_ADD)
	// Add event.Fd to the monitoring list of kq.fd
	_, err := syscall.Kevent(kQueue.Fd, []syscall.Kevent_t{kqEvent}, nil, nil)
	return err
}

func (kQueue *KQueue) Wait() ([]Event, error) {
	n, err := syscall.Kevent(kQueue.Fd, nil, kQueue.kqEvents, nil)
	if err != nil {
		return nil, err
	}
	for i := 0; i < n; i++ {
		kQueue.genericEvents[i] = CreateEvent(kQueue.kqEvents[i])
	}
	return kQueue.genericEvents[:n], nil
}

func (kQueue *KQueue) Close() error {
	return syscall.Close(kQueue.Fd)
}
