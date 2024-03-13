package encryption

import (
	"math/rand"
	"time"
)

func NextSymmetricKey() []byte {
	var key []byte
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < 16; i++ {
		key = append(key, uint8(r.Uint64()))
	}
	return key
}
