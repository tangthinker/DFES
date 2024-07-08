package main

import (
	"context"
	"fmt"
	"github.com/tangthinker/DFES/gateway"
	"github.com/tangthinker/DFES/gateway/proto"
	proto2 "github.com/tangthinker/DFES/mate-server/proto"
	"github.com/tangthinker/DFES/utils"
	"log"
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

	dresp, err := mateClient.Delete(context.Background(), &proto2.DeleteRequest{
		DataId: "1769661300763295744",
	})
	log.Println(dresp)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(dresp)
	fmt.Println(dresp.LeaderMateServerAddr)
	if dresp.GetCode() == proto2.MateCode_NotLeader { // 因为选举存在一定延迟，第二次请求也不一定为leader节点，可多次重试
		leaderClient := utils.NewMateServerClient(dresp.LeaderMateServerAddr)
		ddresp, err := leaderClient.Delete(context.Background(), &proto2.DeleteRequest{
			DataId: "mate-node-1.00000000000000000000",
		})
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(ddresp)
	}
}
