package data_server

type DataService interface {
	Push(data []byte) string
	Get(id string) []byte
	Delete(id string) bool
}

type Fragment struct {
	EncryptKey  []byte
	EncryptData []byte
	FragmentId  string
}
