package main

import (
	"context"
	"fmt"
	"github.com/tangthinker/DFES/gateway"
	"github.com/tangthinker/DFES/gateway/proto"
	proto2 "github.com/tangthinker/DFES/mate-server/proto"
	"github.com/tangthinker/DFES/utils"
	"log"
	"os"
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
	b, _ := os.ReadFile("/Users/tangyubin/Downloads/ideaIU-2023.3.4.dmg")

	presp, err := mateClient.Push(context.Background(), &proto2.PushRequest{
		Data: b,
	})
	log.Println(presp)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(presp)
	fmt.Println(presp.LeaderMateServerAddr)
	if presp.GetCode() == proto2.MateCode_NotLeader { // 因为选举存在一定延迟，第二次请求也不一定为leader节点，可多次重试
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
