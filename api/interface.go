package api

import "io"

type Api interface {
	Push(data []byte) string
	PushStream(stream io.Reader) string
	Get(id string) []byte
	GetStream(id string) io.Reader
}
