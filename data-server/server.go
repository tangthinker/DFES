package data_server

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"tangthinker.work/DFES/encryption"
	idGenerator "tangthinker.work/DFES/id-generator"
	"tangthinker.work/DFES/utils"
)

const (
	DefaultDataStorePath = "./data/"
	RegistryAddr         = "127.0.0.1:60001"
)

type DataServer struct {
	privateKey          []byte
	publicKey           []byte
	serverName          string
	storePath           string
	asymmetricEncryptor *encryption.AsymmetricEncryptor
	symmetricEncryptor  *encryption.SymmetricEncryptor
	idGen               *idGenerator.SequenceIdGenerator
}

func NewDataServer(privateKey []byte, publicKey []byte, idPrefix string) *DataServer {
	return &DataServer{
		privateKey:          privateKey,
		publicKey:           publicKey,
		asymmetricEncryptor: &encryption.AsymmetricEncryptor{},
		symmetricEncryptor:  &encryption.SymmetricEncryptor{},
		idGen:               idGenerator.NewSequenceIdGenerator(idPrefix),
	}
}

func (ds *DataServer) Push(data []byte) string {
	fragmentKey := encryption.NextSymmetricKey()
	var fragment Fragment
	fragment.EncryptKey = ds.asymmetricEncryptor.Encrypt(ds.publicKey, fragmentKey, encryption.RSA)
	fragment.EncryptData = ds.symmetricEncryptor.Encrypt(fragmentKey, data, encryption.AES)
	fragment.FragmentId = ds.idGen.Next()
	store(ds.storePath, &fragment)
	return fragment.FragmentId
}

func (ds *DataServer) Get(id string) []byte {
	fragment := restore(ds.storePath, id)
	fragmentKey := ds.asymmetricEncryptor.Decrypt(ds.privateKey, fragment.EncryptKey, encryption.RSA)
	data := ds.symmetricEncryptor.Decrypt(fragmentKey, fragment.EncryptData, encryption.AES)
	return data
}

func (ds *DataServer) Delete(id string) bool {
	err := deleteById(ds.storePath, id)
	return err == nil
}

func store(path string, fragment *Fragment) {
	b, err := json.Marshal(fragment)
	if err != nil {
		log.Println("json marshal error:", err)
		return
	}
	filename := "data." + fragment.FragmentId
	_ = utils.CreateDirIfNotExist(filepath.Join(path))
	file, err := os.OpenFile(filepath.Join(path, filename), os.O_CREATE|os.O_WRONLY, 0700)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("close file error:", err)
		}
	}(file)
	if err != nil {
		log.Println("create file error:", err)
		return
	}
	_, err = file.Write(b)
	if err != nil {
		log.Println("write file error:", err)
		return
	}
}

func restore(path string, id string) *Fragment {
	filename := "data." + id
	b, err := os.ReadFile(filepath.Join(path, filename))
	if err != nil {
		log.Println("read file error:", err)
		return nil
	}
	var fragment Fragment
	err = json.Unmarshal(b, &fragment)
	if err != nil {
		log.Println("json unmarshal error:", err)
	}
	return &fragment
}

func deleteById(path string, id string) error {
	filename := "data." + id
	err := os.Remove(filepath.Join(path, filename))
	if err != nil {
		log.Println("remove file error:", err)
		return err
	}
	return nil
}
