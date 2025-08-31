package io_multiplexing

const OpRead int16 = 0
const OpWrite int16 = 1

type Event struct {
	Fd        uint64
	operation int16
}

func NewEvent(fd uint64, operation int16) Event {
	return Event{
		Fd:        fd,
		operation: operation,
	}
}

type IOMultiplexer interface {
	Monitor(event Event) error
	Wait() ([]Event, error)
	Close() error
}
