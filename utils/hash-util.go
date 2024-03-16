package utils

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
)

func Hash(data []byte) string {
	c := md5.New()
	c.Write(data)
	bytes := c.Sum(nil)
	return hex.EncodeToString(bytes)
}

type HashCoder struct {
	c hash.Hash
}

func NewHashCoder() *HashCoder {
	return &HashCoder{
		c: md5.New(),
	}
}

func (hc *HashCoder) Join(data []byte) {
	hc.c.Write(data)
}

func (hc *HashCoder) Get() string {
	bytes := hc.c.Sum(nil)
	return hex.EncodeToString(bytes)
}
