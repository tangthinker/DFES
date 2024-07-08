package main

import (
	"github.com/tangthinker/DFES/gateway"
	pb "github.com/tangthinker/DFES/gateway/proto"
	"github.com/tangthinker/DFES/utils"
	"google.golang.org/grpc"
	"time"
)

func main() {
	gateway.Init()
	utils.StartGrpcServer(gateway.DefaultRegistryServerAddr, func(server *grpc.Server) {
		pb.RegisterRegistryServer(server, &gateway.RpcServer{})
		go func() {
			for {
				gateway.PrintOnlineServices()
				time.Sleep(10 * time.Second)
			}
		}()
	})
}
