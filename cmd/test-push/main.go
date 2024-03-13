package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"tangthinker.work/DFES/gateway"
	"tangthinker.work/DFES/gateway/proto"
	proto2 "tangthinker.work/DFES/mate-server/proto"
	"tangthinker.work/DFES/utils"
)

func main() {
	//conn := utils.NewGrpcClient(gateway.DefaultRegistryServerAddr)
	//utils.RegisterServer(conn, &pb.RegisterInfo{
	//	ServiceName: "shanliao",
	//	ServiceType: gateway.MateService,
	//	ServiceAddress: &pb.ServiceAddress{
	//		Host: "127.0.0.1",
	//		Port: "8080",
	//	},
	//	ServiceInterfaces: make([]*pb.ServiceInterface, 0),
	//	HeartbeatAddress:  "127.0.0.1:8080",
	//})
	//time.Sleep(20 * time.Second)

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
	b, _ := os.ReadFile("./api/interface.go")

	presp, err := mateClient.Push(context.Background(), &proto2.PushRequest{
		Data: b,
	})
	log.Println(presp)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(presp)
	if presp.GetCode() == proto2.PushCode_NotLeader {
		leaderClient := utils.NewMateServerClient(presp.LeaderMateServerAddr)
		ppresp, err := leaderClient.Push(context.Background(), &proto2.PushRequest{
			Data: b,
		})
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(ppresp)
	}
}
