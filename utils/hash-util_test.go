package utils

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {

	data := []byte("shanliao")
	fmt.Println(Hash(data))

}
