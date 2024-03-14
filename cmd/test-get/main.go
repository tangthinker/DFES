package main

import (
	"context"
	"log"
	"tangthinker.work/DFES/gateway"
	"tangthinker.work/DFES/gateway/proto"
	proto2 "tangthinker.work/DFES/mate-server/proto"
	"tangthinker.work/DFES/utils"
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
		DataId: "mate-node-1.00000000000000000000",
	})

	log.Println(err)
	log.Println(gresp)
}
