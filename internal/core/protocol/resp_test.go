package protocol_test

import (
	"bytes"
	"testing"

	"github.com/tcnam/redis_go/internal/core/protocol"
)

func TestEncodeSimpleString(t *testing.T) {
	testCases := map[string][]byte{
		"OK":       []byte("+OK\r\n"),
		"PONG":     []byte("+PONG\r\n"),
		"QUEUED":   []byte("+QUEUED\r\n"),
		"RESET":    []byte("+RESET\r\n"),
		"CONTINUE": []byte("+CONTINUE\r\n"),
		"NOAUTH":   []byte("+NOAUTH\r\n"),
	}

	for key, expectedVal := range testCases {
		realVal := protocol.Encode(key, true)
		// log.Printf("%q", realVal)
		// log.Printf("%q", expectedVal)
		if bytes.Equal(expectedVal, realVal) == false {
			t.Fail()
		}
	}
}

func TestEncodeBulkString(t *testing.T) {
	testCases := map[string][]byte{
		"Hello":       []byte("$5\r\nHello\r\n"),
		"Hello World": []byte("$11\r\nHello World\r\n"),
	}

	for key, expectedVal := range testCases {
		realVal := protocol.Encode(key, false)
		// log.Printf("%q", realVal)
		// log.Printf("%q", expectedVal)
		if bytes.Equal(expectedVal, realVal) == false {
			t.Fail()
		}
	}
}
