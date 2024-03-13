package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Hash(data []byte) string {
	c := md5.New()
	c.Write(data)
	bytes := c.Sum(nil)
	return hex.EncodeToString(bytes)
}
