package main

import (
	"context"
	"github.com/tangthinker/DFES/gateway"
	"github.com/tangthinker/DFES/gateway/proto"
	proto2 "github.com/tangthinker/DFES/mate-server/proto"
	"github.com/tangthinker/DFES/utils"
	"log"
	"os"
)

func main() {
	registryCenter := utils.NewRegistryClient(gateway.DefaultRegistryServerAddr)
	resp, err := registryCenter.GetProvideService(context.Background(), &proto.GetProvideInfo{
		ServiceType: gateway.MateService,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
	mateClient := utils.NewMateServerClient(
		resp.GetProvideService().ServiceAddress.Host + ":" + resp.GetProvideService().ServiceAddress.Port)
	//b, _ := os.ReadFile("./api/interface.go")
	gresp, err := mateClient.Get(context.Background(), &proto2.GetRequest{
		DataId: "1769661300763295744",
	})
	log.Println(err)
	log.Println(gresp.GetResult)
	if !gresp.GetResult && gresp.GetCode() == proto2.MateCode_FileNotExist {
		log.Println("get file not exist")
		return
	}

	err = os.WriteFile("./data/test.dmg", gresp.Data, 0700)
	if err != nil {
		log.Println(err)
	}
	log.Println(err)
	//log.Println(gresp)
}
