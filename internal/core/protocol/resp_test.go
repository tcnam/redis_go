package protocol_test

import (
	"bytes"
	"testing"

	"github.com/tcnam/redis_go/internal/core/protocol"
)

type TestCaseEncode struct {
	input          interface{}
	expectedOutput []byte
}

func TestEncodeSimpleString(t *testing.T) {
	var testCases []TestCaseEncode = make([]TestCaseEncode, 0, 6)

	var testCase1 TestCaseEncode = TestCaseEncode{
		input:          "OK",
		expectedOutput: []byte("+OK\r\n"),
	}
	var testCase2 TestCaseEncode = TestCaseEncode{
		input:          "PONG",
		expectedOutput: []byte("+PONG\r\n"),
	}
	var testCase3 TestCaseEncode = TestCaseEncode{
		input:          "QUEUED",
		expectedOutput: []byte("+QUEUED\r\n"),
	}
	var testCase4 TestCaseEncode = TestCaseEncode{
		input:          "RESET",
		expectedOutput: []byte("+RESET\r\n"),
	}
	var testCase5 TestCaseEncode = TestCaseEncode{
		input:          "CONTINUE",
		expectedOutput: []byte("+CONTINUE\r\n"),
	}
	var testCase6 TestCaseEncode = TestCaseEncode{
		input:          "NOAUTH",
		expectedOutput: []byte("+NOAUTH\r\n"),
	}

	testCases = append(testCases, testCase1, testCase2, testCase3, testCase4, testCase5, testCase6)

	for i := 0; i < len(testCases); i++ {
		realOutput := protocol.Encode(testCases[i].input, true)
		// log.Printf("%q", realOutput)
		// log.Printf("%q", testCases[i].expectedOutput)
		if !bytes.Equal(realOutput, testCases[i].expectedOutput) {
			t.Fail()
		}
	}
}

func TestEncodeBulkString(t *testing.T) {
	var testCases []TestCaseEncode = make([]TestCaseEncode, 0, 2)

	var testCase1 TestCaseEncode = TestCaseEncode{
		input:          "Hello",
		expectedOutput: []byte("$5\r\nHello\r\n"),
	}
	var testCase2 TestCaseEncode = TestCaseEncode{
		input:          "Hello World",
		expectedOutput: []byte("$11\r\nHello World\r\n"),
	}

	testCases = append(testCases, testCase1, testCase2)

	for i := 0; i < len(testCases); i++ {
		realOutput := protocol.Encode(testCases[i].input, false)
		// log.Printf("%q", realOutput)
		// log.Printf("%q", testCases[i].expectedOutput)
		if !bytes.Equal(realOutput, testCases[i].expectedOutput) {
			t.Fail()
		}
	}
}

func TestEncodeInteger(t *testing.T) {
	var testCases []TestCaseEncode = make([]TestCaseEncode, 0, 2)

	var testCase1 TestCaseEncode = TestCaseEncode{
		input:          0,
		expectedOutput: []byte(":0\r\n"),
	}
	var testCase2 TestCaseEncode = TestCaseEncode{
		input:          123,
		expectedOutput: []byte(":123\r\n"),
	}
	var testCase3 TestCaseEncode = TestCaseEncode{
		input:          -123,
		expectedOutput: []byte(":-123\r\n"),
	}

	testCases = append(testCases, testCase1, testCase2, testCase3)

	for i := 0; i < len(testCases); i++ {
		realOutput := protocol.Encode(testCases[i].input, false)
		// log.Printf("%q", realOutput)
		// log.Printf("%q", testCases[i].expectedOutput)
		if !bytes.Equal(realOutput, testCases[i].expectedOutput) {
			t.Fail()
		}
	}
}

func TestEncodeStringArray(t *testing.T) {
	var testCases []TestCaseEncode = make([]TestCaseEncode, 0, 3)

	var testCase1 TestCaseEncode = TestCaseEncode{
		input:          []string{},
		expectedOutput: []byte("*0\r\n"),
	}
	var testCase2 TestCaseEncode = TestCaseEncode{
		input:          []string{"hello"},
		expectedOutput: []byte("*1\r\n$5\r\nhello\r\n"),
	}
	var testCase3 TestCaseEncode = TestCaseEncode{
		input:          []string{"foo", "bar"},
		expectedOutput: []byte("*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"),
	}

	testCases = append(testCases, testCase1, testCase2, testCase3)

	for i := 0; i < len(testCases); i++ {
		realOutput := protocol.Encode(testCases[i].input, false)
		// log.Printf("%q", realOutput)
		// log.Printf("%q", testCases[i].expectedOutput)
		if !bytes.Equal(realOutput, testCases[i].expectedOutput) {
			t.Fail()
		}
	}
}

func TestEncodeInterfaceArray(t *testing.T) {
	var testCases []TestCaseEncode = make([]TestCaseEncode, 0, 3)

	var testCase1 TestCaseEncode = TestCaseEncode{
		input:          []interface{}{"foo", "bar"},
		expectedOutput: []byte("*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"),
	}
	var testCase2 TestCaseEncode = TestCaseEncode{
		input:          []interface{}{-42, "hello"},
		expectedOutput: []byte("*2\r\n:-42\r\n$5\r\nhello\r\n"),
	}
	var testCase3 TestCaseEncode = TestCaseEncode{
		input:          []interface{}{},
		expectedOutput: []byte("*0\r\n"),
	}

	testCases = append(testCases, testCase1, testCase2, testCase3)

	for i := 0; i < len(testCases); i++ {
		realOutput := protocol.Encode(testCases[i].input, false)
		// log.Printf("%q", realOutput)
		// log.Printf("%q", testCases[i].expectedOutput)
		if !bytes.Equal(realOutput, testCases[i].expectedOutput) {
			t.Fail()
		}
	}
}
