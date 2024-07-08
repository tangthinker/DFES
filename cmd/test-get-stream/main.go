package main

import (
	"context"
	"github.com/tangthinker/DFES/gateway"
	"github.com/tangthinker/DFES/gateway/proto"
	proto2 "github.com/tangthinker/DFES/mate-server/proto"
	"github.com/tangthinker/DFES/utils"
	"io"
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
	grespStream, err := mateClient.GetStream(context.Background(), &proto2.GetRequest{
		DataId: "1769668805266595840",
	})
	if err != nil {
		log.Fatal(err)
	}
	var ret []byte
	for {
		data, err := grespStream.Recv()
		if err == io.EOF {
			log.Println("receive all")
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if !data.GetResult {
			log.Println(data.GetCode())
			return
		}
		ret = append(ret, data.GetData()...)
	}

	err = os.WriteFile("./data/test.iso", ret, 0700)
	if err != nil {
		log.Println(err)
	}
	log.Println(err)
}
