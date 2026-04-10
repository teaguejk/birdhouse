package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"
)

func RandomHash(filename string) string {
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	md5hash := md5.New()
	md5hash.Write([]byte(filename + timestamp))
	hash := hex.EncodeToString(md5hash.Sum(nil))
	return hash
}
