package encryption

type Encryptor interface {
	Encrypt(key []byte, data []byte, encryptType EncryptType) []byte
	Decrypt(key []byte, data []byte, encryptType EncryptType) []byte
}

type EncryptType string
