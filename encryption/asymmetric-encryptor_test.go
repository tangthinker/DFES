package encryption

import (
	"fmt"
	"testing"
)

func TestAsymmetricEncryptor_Encrypt(t *testing.T) {
	ae := AsymmetricEncryptor{}
	pri, pub, err := ae.GenerateKey(RSA)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pri)
	fmt.Println(pub)
	cipher, err := ae.Encrypt(pub, text, RSA)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(cipher))
	plain, err := ae.Decrypt(pri, cipher, RSA)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(plain))
}
