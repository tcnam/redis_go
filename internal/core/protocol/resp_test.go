package protocol_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/tcnam/redis_go/internal/core/protocol"
)

type TestCaseEncode struct {
	input  interface{}
	output []byte
}

type TestCaseDecode struct {
	input  []byte
	output interface{}
}

func TestEncodeSimpleString(t *testing.T) {
	var testCases []TestCaseEncode = make([]TestCaseEncode, 0, 6)

	var testCase1 TestCaseEncode = TestCaseEncode{
		input:  "OK",
		output: []byte("+OK\r\n"),
	}
	var testCase2 TestCaseEncode = TestCaseEncode{
		input:  "PONG",
		output: []byte("+PONG\r\n"),
	}
	var testCase3 TestCaseEncode = TestCaseEncode{
		input:  "QUEUED",
		output: []byte("+QUEUED\r\n"),
	}
	var testCase4 TestCaseEncode = TestCaseEncode{
		input:  "RESET",
		output: []byte("+RESET\r\n"),
	}
	var testCase5 TestCaseEncode = TestCaseEncode{
		input:  "CONTINUE",
		output: []byte("+CONTINUE\r\n"),
	}
	var testCase6 TestCaseEncode = TestCaseEncode{
		input:  "NOAUTH",
		output: []byte("+NOAUTH\r\n"),
	}

	testCases = append(testCases, testCase1, testCase2, testCase3, testCase4, testCase5, testCase6)

	for i := 0; i < len(testCases); i++ {
		realOutput := protocol.Encode(testCases[i].input, true)
		// log.Printf("%q", realOutput)
		// log.Printf("%q", testCases[i].output)
		if !bytes.Equal(realOutput, testCases[i].output) {
			t.Fail()
		}
	}
}

func TestEncodeBulkString(t *testing.T) {
	var testCases []TestCaseEncode = make([]TestCaseEncode, 0, 2)

	var testCase1 TestCaseEncode = TestCaseEncode{
		input:  "Hello",
		output: []byte("$5\r\nHello\r\n"),
	}
	var testCase2 TestCaseEncode = TestCaseEncode{
		input:  "Hello World",
		output: []byte("$11\r\nHello World\r\n"),
	}

	testCases = append(testCases, testCase1, testCase2)

	for i := 0; i < len(testCases); i++ {
		realOutput := protocol.Encode(testCases[i].input, false)
		// log.Printf("%q", realOutput)
		// log.Printf("%q", testCases[i].output)
		if !bytes.Equal(realOutput, testCases[i].output) {
			t.Fail()
		}
	}
}

func TestEncodeInteger(t *testing.T) {
	var testCases []TestCaseEncode = make([]TestCaseEncode, 0, 2)

	var testCase1 TestCaseEncode = TestCaseEncode{
		input:  0,
		output: []byte(":0\r\n"),
	}
	var testCase2 TestCaseEncode = TestCaseEncode{
		input:  123,
		output: []byte(":123\r\n"),
	}
	var testCase3 TestCaseEncode = TestCaseEncode{
		input:  -123,
		output: []byte(":-123\r\n"),
	}

	testCases = append(testCases, testCase1, testCase2, testCase3)

	for i := 0; i < len(testCases); i++ {
		realOutput := protocol.Encode(testCases[i].input, false)
		// log.Printf("%q", realOutput)
		// log.Printf("%q", testCases[i].output)
		if !bytes.Equal(realOutput, testCases[i].output) {
			t.Fail()
		}
	}
}

func TestEncodeStringArray(t *testing.T) {
	var testCases []TestCaseEncode = make([]TestCaseEncode, 0, 3)

	var testCase1 TestCaseEncode = TestCaseEncode{
		input:  []string{},
		output: []byte("*0\r\n"),
	}
	var testCase2 TestCaseEncode = TestCaseEncode{
		input:  []string{"hello"},
		output: []byte("*1\r\n$5\r\nhello\r\n"),
	}
	var testCase3 TestCaseEncode = TestCaseEncode{
		input:  []string{"foo", "bar"},
		output: []byte("*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"),
	}

	testCases = append(testCases, testCase1, testCase2, testCase3)

	for i := 0; i < len(testCases); i++ {
		realOutput := protocol.Encode(testCases[i].input, false)
		// log.Printf("%q", realOutput)
		// log.Printf("%q", testCases[i].output)
		if !bytes.Equal(realOutput, testCases[i].output) {
			t.Fail()
		}
	}
}

func TestEncodeInterfaceArray(t *testing.T) {
	var testCases []TestCaseEncode = make([]TestCaseEncode, 0, 3)

	var testCase1 TestCaseEncode = TestCaseEncode{
		input:  []interface{}{"foo", "bar"},
		output: []byte("*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"),
	}
	var testCase2 TestCaseEncode = TestCaseEncode{
		input:  []interface{}{-42, "hello"},
		output: []byte("*2\r\n:-42\r\n$5\r\nhello\r\n"),
	}
	var testCase3 TestCaseEncode = TestCaseEncode{
		input:  []interface{}{},
		output: []byte("*0\r\n"),
	}

	testCases = append(testCases, testCase1, testCase2, testCase3)

	for i := 0; i < len(testCases); i++ {
		realOutput := protocol.Encode(testCases[i].input, false)
		// log.Printf("%q", realOutput)
		// log.Printf("%q", testCases[i].output)
		if !bytes.Equal(realOutput, testCases[i].output) {
			t.Fail()
		}
	}
}

func TestDecodeSimpleString(t *testing.T) {
	var testCases []TestCaseDecode = make([]TestCaseDecode, 0, 6)

	var testCase1 TestCaseDecode = TestCaseDecode{
		input:  []byte("+OK\r\n"),
		output: "OK",
	}
	var testCase2 TestCaseDecode = TestCaseDecode{
		input:  []byte("+PONG\r\n"),
		output: "PONG",
	}
	var testCase3 TestCaseDecode = TestCaseDecode{
		input:  []byte("+QUEUED\r\n"),
		output: "QUEUED",
	}
	var testCase4 TestCaseDecode = TestCaseDecode{
		input:  []byte("+RESET\r\n"),
		output: "RESET",
	}
	var testCase5 TestCaseDecode = TestCaseDecode{
		input:  []byte("+CONTINUE\r\n"),
		output: "CONTINUE",
	}
	var testCase6 TestCaseDecode = TestCaseDecode{
		input:  []byte("+NOAUTH\r\n"),
		output: "NOAUTH",
	}

	testCases = append(testCases, testCase1, testCase2, testCase3, testCase4, testCase5, testCase6)
	for i := 0; i < len(testCases); i++ {
		realOutput, _, _ := protocol.Decode(testCases[i].input)
		if realOutput != testCases[i].output {
			t.Fail()
		}
	}
}

func TestError(t *testing.T) {
	cases := map[string]string{
		"-Error message\r\n": "Error message",
	}
	for k, v := range cases {
		value, _, _ := protocol.Decode([]byte(k))
		if v != value {
			t.Fail()
		}
	}
}

func TestDecodeInteger(t *testing.T) {
	var testCases []TestCaseDecode = make([]TestCaseDecode, 0, 4)
	var testCase1 TestCaseDecode = TestCaseDecode{
		input:  []byte(":0\r\n"),
		output: int64(0),
	}
	var testCase2 TestCaseDecode = TestCaseDecode{
		input:  []byte(":1000\r\n"),
		output: int64(1000),
	}
	var testCase3 TestCaseDecode = TestCaseDecode{
		input:  []byte(":+1000\r\n"),
		output: int64(1000),
	}
	var testCase4 TestCaseDecode = TestCaseDecode{
		input:  []byte(":-1000\r\n"),
		output: int64(-1000),
	}
	testCases = append(testCases, testCase1, testCase2, testCase3, testCase4)
	for i := 0; i < len(testCases); i++ {
		realOutput, _, _ := protocol.Decode(testCases[i].input)
		if realOutput != testCases[i].output {
			t.Fail()
		}
	}
}

func TestDecodeBulkString(t *testing.T) {
	var testCases []TestCaseDecode = make([]TestCaseDecode, 0, 2)
	var testCase1 TestCaseDecode = TestCaseDecode{
		input:  []byte("$5\r\nhello\r\n"),
		output: "hello",
	}
	var testCase2 TestCaseDecode = TestCaseDecode{
		input:  []byte("$0\r\n\r\n"),
		output: "",
	}
	testCases = append(testCases, testCase1, testCase2)

	for i := 0; i < len(testCases); i++ {
		realOutput, _, _ := protocol.Decode(testCases[i].input)
		if realOutput != testCases[i].output {
			t.Fail()
		}
	}
}

func TestDecodeArray(t *testing.T) {
	cases := map[string][]interface{}{
		"*0\r\n":                                        {},
		"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n":          {"hello", "world"},
		"*3\r\n:1\r\n:2\r\n:3\r\n":                      {int64(1), int64(2), int64(3)},
		"*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n": {int64(1), int64(2), int64(3), int64(4), "hello"},
		"*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Hello\r\n$5\r\nWorld\r\n": {[]int64{int64(1), int64(2), int64(3)}, []interface{}{"Hello", "World"}},
	}
	for k, v := range cases {
		value, _, _ := protocol.Decode([]byte(k))
		array := value.([]interface{})
		if len(array) != len(v) {
			t.Fail()
		}
		for i := range array {
			if fmt.Sprintf("%v", v[i]) != fmt.Sprintf("%v", array[i]) {
				t.Fail()
			}
		}
	}
}
