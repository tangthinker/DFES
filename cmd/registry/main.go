package main

import (
	"github.com/shanliao420/DFES/gateway"
	pb "github.com/shanliao420/DFES/gateway/proto"
	"github.com/shanliao420/DFES/utils"
	"google.golang.org/grpc"
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
