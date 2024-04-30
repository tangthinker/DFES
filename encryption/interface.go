package encryption

type Encryptor interface {
	Encrypt(key []byte, data []byte, encryptType EncryptType) ([]byte, error)
	Decrypt(key []byte, data []byte, encryptType EncryptType) ([]byte, error)
}

type EncryptType string
