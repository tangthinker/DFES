package data_server

import (
	"log"
	"os"
	"path/filepath"
	"tangthinker.work/DFES/encryption"
	pb "tangthinker.work/DFES/gateway/proto"
	"tangthinker.work/DFES/utils"
)

const (
	DefaultKeyStorePath = "./data/key/"
	DefaultServerName   = "shanliao-data-node-1"
)

var (
	registerClient pb.RegistryClient
	ServerHost     string
	ServerPort     string
	dataService    *DataServer = NewDataServer(nil, nil, DefaultServerName)
)

func Init() {
	registerClient = pb.NewRegistryClient(utils.NewGrpcClient(RegistryAddr))
	pri, pub := getKey(DefaultKeyStorePath)
	dataService.privateKey = pri
	dataService.publicKey = pub
}

func SetDataServerName(serverName string) {
	dataService.serverName = serverName
	dataService.idGen.ResetPrefix(serverName)
	dataService.storePath = DefaultDataStorePath + serverName + "/"
}

func getKey(path string) ([]byte, []byte) {
	priName := dataService.serverName + ".private.key"
	pubName := dataService.serverName + ".public.key"
	pri, err := os.ReadFile(filepath.Join(path, priName))
	if err != nil {
		pri, pub := (&encryption.AsymmetricEncryptor{}).GenerateKey(encryption.RSA)
		_ = utils.CreateDirIfNotExist(path)
		err = os.WriteFile(filepath.Join(path, priName), pri, 0700)
		err = os.WriteFile(filepath.Join(path, pubName), pub, 0700)
		if err != nil {
			log.Println("key save error:", err)
		}
		return pri, pub
	}
	pub, err := os.ReadFile(filepath.Join(path, pubName))
	if err != nil {
		pri, pub := (&encryption.AsymmetricEncryptor{}).GenerateKey(encryption.RSA)
		_ = utils.CreateDirIfNotExist(path)
		err = os.WriteFile(priName, pri, 0700)
		err = os.WriteFile(pubName, pub, 0700)
		if err != nil {
			log.Println("key save error:", err)
		}
		return pri, pub
	}
	return pri, pub
}
