package main

import (
	"context"
	"github.com/shanliao420/DFES/gateway"
	"github.com/shanliao420/DFES/gateway/proto"
	proto2 "github.com/shanliao420/DFES/mate-server/proto"
	"github.com/shanliao420/DFES/utils"
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
		DataId: "mate-node-1.00000000000000000001",
	})
	log.Println(err)

	err = os.WriteFile("./data/test.iso", gresp.Data, 0700)
	if err != nil {
		log.Println(err)
	}
	log.Println(err)
	//log.Println(gresp)
}
