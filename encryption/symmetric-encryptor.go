package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

// initialization vector
var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

type SymmetricEncryptor struct {
}

func (se *SymmetricEncryptor) Encrypt(key []byte, data []byte, encryptType EncryptType) ([]byte, error) {
	switch encryptType {
	case AES:
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, fmt.Errorf("key is invalid")
		}
		cfb := cipher.NewCFBEncrypter(block, commonIV)
		result := make([]byte, len(data))
		cfb.XORKeyStream(result, data)
		return result, nil
	default:
		return nil, fmt.Errorf("encrypt type invalid")
	}
}

func (se *SymmetricEncryptor) Decrypt(key []byte, data []byte, encryptType EncryptType) ([]byte, error) {
	switch encryptType {
	case AES:
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, fmt.Errorf("key is invalid")
		}
		cfb := cipher.NewCFBDecrypter(block, commonIV)
		result := make([]byte, len(data))
		cfb.XORKeyStream(result, data)
		return result, nil
	default:
		return nil, fmt.Errorf("encrypt type invalid")
	}
}

const (
	AES = EncryptType("aes")
)
