package encryption

import (
	"fmt"
	"testing"
)

//var key = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

var text = []byte(string("shanliao"))

func TestSymmetricEncryptor_Encrypt(t *testing.T) {
	se := SymmetricEncryptor{}
	key := Next()
	fmt.Println(key)
	cipher := se.Encrypt(key, text, AES)
	fmt.Println(string(cipher))
	plain := se.Decrypt(key, cipher, AES)
	fmt.Println(string(plain))
}
