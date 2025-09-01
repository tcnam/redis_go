package constant

import "time"

const CRLF string = "\r\n"

var RespNil []byte = []byte("$-1\r\n")
var RespOk []byte = []byte("+OK\r\n")
var TtlKeyNotExist []byte = []byte(":-2\r\n")
var TtlKeyExistNoExpire []byte = []byte(":-1\r\n")
var ActiveExpireFrequency time.Duration = 100 * time.Millisecond
var ActiveExpireSampleSize int = 20
var ActiveExpireThreshold float64 = 0.1
