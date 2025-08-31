//go:build darwin || freebsd || netbsd || openbsd

package io_multiplexing

import (
	"log"
	"syscall"
)

func CreateEvent(ke syscall.Kevent_t) Event {
	var operation int16
	switch ke.Filter {
	case syscall.EVFILT_READ:
		operation = syscall.EVFILT_READ
	case syscall.EVFILT_WRITE:
		operation = syscall.EVFILT_WRITE
	default:
		log.Fatalf("Kevent_t filter not support: %d\n", ke.Filter)
	}
	return Event{
		Fd:        ke.Ident,
		operation: operation,
	}
}

func (event Event) toNative(flags uint16) syscall.Kevent_t {
	var filter int16
	switch event.operation {
	case OpRead:
		filter = syscall.EVFILT_READ
	case OpWrite:
		filter = syscall.EVFILT_WRITE
	default:
		log.Fatalf("Operation not supported: %d\n", event.operation)
	}

	return syscall.Kevent_t{
		Ident:  event.Fd,
		Filter: filter,
		Flags:  flags,
	}
}
