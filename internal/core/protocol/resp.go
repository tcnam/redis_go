package protocol

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/tcnam/redis_go/internal/constant"
)

func Encode(value interface{}, isSimpleString bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimpleString {
			return encodeSimpleString(v)
		} else {
			return encodeBulkString(v)
		}
	case int64, int32, int16, int8, int:
		return encodeInteger(v)
	case error:
		return encodeSimpleError(v)
	case []string:
		return encodeStringArray(v)
	case [][]string:
		return encodeStringMatrix(v)
	case []interface{}:
		return encodeInterfaceArray(v)
	default:
		return constant.RespNil
	}
}

func encodeSimpleString(value string) []byte {
	return []byte(fmt.Sprintf("+%s%s", value, constant.CRLF))
}

func encodeBulkString(value string) []byte {
	return []byte(fmt.Sprintf("$%d%s%s%s", len(value), constant.CRLF, value, constant.CRLF))
}

func encodeInteger(value interface{}) []byte {
	return []byte(fmt.Sprintf(":%d%s", value, constant.CRLF))
}

func encodeSimpleError(value error) []byte {
	return []byte(fmt.Sprintf("-%s%s", value, constant.CRLF))
}

func encodeStringArray(sa []string) []byte {
	var b []byte
	buf := bytes.NewBuffer(b)
	for _, s := range sa {
		buf.Write(encodeBulkString(s))
	}
	return []byte(fmt.Sprintf("*%d%s%s", len(sa), constant.CRLF, buf.Bytes()))
}

func encodeStringMatrix(sm [][]string) []byte {
	var b []byte
	buf := bytes.NewBuffer(b)
	for _, sa := range sm {
		buf.Write(encodeStringArray(sa))
	}
	return []byte(fmt.Sprintf("*%d%s%s", len(sm), constant.CRLF, buf.Bytes()))
}

func encodeInterfaceArray(ia []interface{}) []byte {
	var b []byte
	buf := bytes.NewBuffer(b)
	for _, x := range ia {
		buf.Write(Encode(x, false))
	}
	return []byte(fmt.Sprintf("*%d%s%s", len(ia), constant.CRLF, buf.Bytes()))
}

func Decode(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}
	switch data[0] {
	case '+':
		return decodeSimpleString(data)
	case '-':
		return decodeSimpleError(data)
	case ':':
		return decodeInteger(data)
	case '$':
		return decodeBulkString(data)
	case '*':
		return decodeArray(data)
	}
	return nil, 0, nil
}

func decodeSimpleString(data []byte) (string, int, error) {
	pos := 1
	for data[pos] != '\r' {
		pos++
	}
	return string(data[1:pos]), pos + 2, nil
}

func decodeInteger(data []byte) (int64, int, error) {
	var res int64 = 0
	var pos int = 1
	var sign int64 = 1

	switch data[pos] {
	case '-':
		sign = -1
		pos++
	case '+':
		pos++
	}

	for data[pos] != '\r' {
		// substract ANSI digit with ANSI of '0' digit to get its value in integer
		res = res*10 + int64(data[pos]-'0')
		pos++
	}
	return sign * res, pos + 2, nil
}

func decodeSimpleError(data []byte) (string, int, error) {
	return decodeSimpleString(data)
}

func decodeBulkString(data []byte) (string, int, error) {
	var startPos int = 1
	for data[startPos] != '\n' {
		startPos++
	}
	startPos++

	var endPos int = startPos
	for data[endPos] != '\r' {
		endPos++
	}

	return string(data[startPos:endPos]), endPos + 2, nil
}

func decodeArray(data []byte) ([]interface{}, int, error) {
	arrSize, ind, _ := decodeInteger(data)
	var arr []interface{} = make([]interface{}, 0, arrSize)
	// log.Printf("%q", data)
	var elementCount int64 = 0
	for elementCount < arrSize && ind < len(data) {
		tempVal, tempInd, _ := Decode(data[ind:])
		// log.Printf("%v", tempVal)
		arr = append(arr, tempVal)
		elementCount++
		ind += tempInd
	}
	// log.Printf("%v", arr)
	return arr, ind, nil
}
