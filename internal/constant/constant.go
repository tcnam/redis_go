package constant

import "time"

var RespNil []byte = []byte("$-1\r\n")
var RespOk []byte = []byte("+OK\r\n")
var TtlKeyNotExist = []byte(":-2\r\n")
var TtlKeyExistNoExpire = []byte(":-1\r\n")
var ActiveExpireFrequency = 100 * time.Millisecond
var ActiveExpireSampleSize = 20
var ActiveExpireThreshold = 0.1
