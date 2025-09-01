package protocol

import (
	"bytes"
	"fmt"

	"github.com/tcnam/redis_go/internal/constant"
)

const CRLF string = "\r\n"

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
		return encodeError(v)
	case []string:
		return encodeStringArray(v)
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

func encodeError(value error) []byte {
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
