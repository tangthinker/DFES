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
	b, _ := os.ReadFile("/Users/tangyubin/Downloads/Win10_22H2_English_x64v1.iso")

	presp, err := mateClient.PushStream(context.Background())
	log.Println(presp)
	if err != nil {
		log.Println(err)
	}
	blen := len(b)
	fragmentSize := 24 * 1024 * 1024
	left := 0
	right := fragmentSize
	for right <= blen {
		_ = presp.Send(&proto2.PushRequest{
			Data: b[left:right],
		})
		log.Println("send data [", left, ":", right, "]")
		left = right
		right += fragmentSize
	}
	if left < blen && right != blen {
		_ = presp.Send(&proto2.PushRequest{
			Data: b[left:blen],
		})
		log.Println("send data [", left, ":", blen, "]")
	}
	pusresp, err := presp.CloseAndRecv()
	log.Println(err)
	log.Println(pusresp)
}
