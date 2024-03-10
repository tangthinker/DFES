package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

// privateKey to decrypt, publicKey to encrypt

type AsymmetricEncryptor struct {
}

func (ae *AsymmetricEncryptor) GenerateKey(encryptType EncryptType) (privateKey []byte, publicKey []byte) {
	switch encryptType {
	case RSA:
		priKey, err := rsa.GenerateKey(rand.Reader, defaultKeyBit)
		if err != nil {
			panic(err)
		}
		privateKey = x509.MarshalPKCS1PrivateKey(priKey) // x509PriKey
		pubKey := priKey.PublicKey
		publicKey = x509.MarshalPKCS1PublicKey(&pubKey)
		return
	default:
		panic("encryption type invalid")
	}
}

func (ae *AsymmetricEncryptor) GenerateKeyWithSize(keyBit int, encryptType EncryptType) (privateKey []byte, publicKey []byte) {
	switch encryptType {
	case RSA:
		priKey, err := rsa.GenerateKey(rand.Reader, keyBit)
		if err != nil {
			panic(err)
		}
		privateKey = x509.MarshalPKCS1PrivateKey(priKey) // x509PriKey
		pubKey := priKey.PublicKey
		publicKey = x509.MarshalPKCS1PublicKey(&pubKey)
		return
	default:
		panic("encryption type invalid")
	}
}

func (ae *AsymmetricEncryptor) Encrypt(publicKey []byte, data []byte, encryptType EncryptType) []byte {
	switch encryptType {
	case RSA:
		pubKey, err := x509.ParsePKCS1PublicKey(publicKey)
		if err != nil {
			panic("key is invalid")
		}
		cipher, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
		if err != nil {
			panic("encryption error")
		}
		return cipher
	default:
		panic("encryption type invalid")
	}
}

func (ae *AsymmetricEncryptor) Decrypt(privateKey []byte, data []byte, encryptType EncryptType) []byte {
	switch encryptType {
	case RSA:
		priKey, err := x509.ParsePKCS1PrivateKey(privateKey)
		if err != nil {
			panic("key is invalid")
		}
		cipher, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, data)
		return cipher
	default:
		panic("encryption type invalid")
	}
}

const (
	RSA = EncryptType("rsa")

	defaultKeyBit = 2048
)
