package data_server

import (
	"context"
	"github.com/tangthinker/DFES/encryption"
	pb "github.com/tangthinker/DFES/gateway/proto"
	idGenerator "github.com/tangthinker/DFES/id-generator"
	"github.com/tangthinker/DFES/utils"
	"log"
	"os"
	"path/filepath"
)

const (
	DefaultKeyStorePath = "./data/key/"
)

var (
	registerClient pb.RegistryClient
	dataService    *DataServer
)

func Init(serverName string) {
	dataService = NewDataServer(nil, nil, serverName)
	dataService.serverName = serverName
	dataService.storePath = DefaultDataStorePath + serverName + "/"
	registerClient = pb.NewRegistryClient(utils.NewGrpcClient(RegistryAddr))
	pri, pub := getKey(DefaultKeyStorePath)
	dataService.privateKey = pri
	dataService.publicKey = pub
	serviceCnt, err := registerClient.GetHistoryAllServiceCnt(context.Background(), nil)
	if err != nil {
		log.Fatalln("init data service in id node err:", err)
	}
	log.Println("init mate server id node -> ", serviceCnt.GetServiceCnt())
	dataService.idGen = idGenerator.NewSnowflakeIdGenerator(serviceCnt.ServiceCnt)
}

func SetDataServerName(serverName string) {
	dataService.serverName = serverName
	dataService.storePath = DefaultDataStorePath + serverName + "/"
}

func getKey(path string) ([]byte, []byte) {
	priName := dataService.serverName + ".private.key"
	pubName := dataService.serverName + ".public.key"
	pri, err := os.ReadFile(filepath.Join(path, priName))
	if err != nil {
		pri, pub, _ := (&encryption.AsymmetricEncryptor{}).GenerateKey(encryption.RSA)
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
		pri, pub, _ := (&encryption.AsymmetricEncryptor{}).GenerateKey(encryption.RSA)
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
