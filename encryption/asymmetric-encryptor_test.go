package encryption

import (
	"fmt"
	"testing"
)

func TestAsymmetricEncryptor_Encrypt(t *testing.T) {
	ae := AsymmetricEncryptor{}
	pri, pub := ae.GenerateKey(RSA)
	fmt.Println(pri)
	fmt.Println(pub)
	cipher := ae.Encrypt(pub, text, RSA)
	fmt.Println(string(cipher))
	plain := ae.Decrypt(pri, cipher, RSA)
	fmt.Println(string(plain))
}
