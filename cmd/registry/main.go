package main

import (
	"google.golang.org/grpc"
	"tangthinker.work/DFES/gateway"
	pb "tangthinker.work/DFES/gateway/proto"
	"tangthinker.work/DFES/utils"
	"time"
)

func main() {
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
