package encryption

import (
	"fmt"
	"testing"
)

//var key = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

var text = []byte(string("shanliao"))

func TestSymmetricEncryptor_Encrypt(t *testing.T) {
	se := SymmetricEncryptor{}
	key := NextSymmetricKey()
	fmt.Println(key)
	cipher, err := se.Encrypt(key, text, AES)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(cipher))
	plain, err := se.Decrypt(key, cipher, AES)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(plain))
}
