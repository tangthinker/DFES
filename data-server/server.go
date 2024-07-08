package data_server

import (
	"bytes"
	"encoding/gob"
	"github.com/tangthinker/DFES/encryption"
	idGenerator "github.com/tangthinker/DFES/id-generator"
	"github.com/tangthinker/DFES/utils"
	"log"
	"os"
	"path/filepath"
)

const (
	DefaultDataStorePath = "./data/"
	RegistryAddr         = "127.0.0.1:6001"
)

type DataServer struct {
	privateKey          []byte
	publicKey           []byte
	serverName          string
	storePath           string
	asymmetricEncryptor *encryption.AsymmetricEncryptor
	symmetricEncryptor  *encryption.SymmetricEncryptor
	idGen               idGenerator.IdGenerator
}

func NewDataServer(privateKey []byte, publicKey []byte, idPrefix string) *DataServer {
	return &DataServer{
		privateKey:          privateKey,
		publicKey:           publicKey,
		asymmetricEncryptor: &encryption.AsymmetricEncryptor{},
		symmetricEncryptor:  &encryption.SymmetricEncryptor{},
	}
}

func (ds *DataServer) Push(data []byte) (string, error) {
	fragmentKey := encryption.NextSymmetricKey()
	var fragment Fragment
	encryptKey, err := ds.asymmetricEncryptor.Encrypt(ds.publicKey, fragmentKey, encryption.RSA)
	if err != nil {
		log.Println("encrypt fragment key error:", err)
		return "", err
	}
	encryptData, err := ds.symmetricEncryptor.Encrypt(fragmentKey, data, encryption.AES)
	if err != nil {
		log.Println("encrypt fragment data error:", err)
		return "", err
	}
	fragment.EncryptKey = encryptKey
	fragment.EncryptData = encryptData
	fragment.FragmentId = ds.idGen.Next()
	store(ds.storePath, &fragment)
	return fragment.FragmentId, nil
}

func (ds *DataServer) Get(id string) ([]byte, error) {
	fragment := restore(ds.storePath, id)
	fragmentKey, err := ds.asymmetricEncryptor.Decrypt(ds.privateKey, fragment.EncryptKey, encryption.RSA)
	if err != nil {
		log.Println("decrypt fragment key error:", err)
		return nil, err
	}
	data, err := ds.symmetricEncryptor.Decrypt(fragmentKey, fragment.EncryptData, encryption.AES)
	if err != nil {
		log.Println("decrypt fragment data error:", err)
		return nil, err
	}
	return data, nil
}

func (ds *DataServer) Delete(id string) bool {
	err := deleteById(ds.storePath, id)
	return err == nil
}

func store(path string, fragment *Fragment) {
	b, err := serialize(fragment)
	if err != nil {
		log.Println("serialize to binary err:", err)
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
	fragment, err := deserialize(b)
	if err != nil {
		log.Println("deserialize to struct err:", err)
		return nil
	}
	return fragment
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

func serialize(fragment *Fragment) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(fragment)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func deserialize(b []byte) (*Fragment, error) {
	var fragment Fragment
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&fragment)
	if err != nil {
		return nil, err
	}
	return &fragment, nil
}
