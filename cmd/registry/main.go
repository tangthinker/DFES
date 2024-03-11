package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"tangthinker.work/DFES/gateway"
	pb "tangthinker.work/DFES/gateway/proto"
	"time"
)

func main() {
	lis, err := net.Listen("tcp", ":60001")
	if err != nil {
		log.Fatal(err)
	}
	srv := grpc.NewServer()
	pb.RegisterRegistryServer(srv, &gateway.RpcServer{})
	log.Printf("server listening at %v", lis.Addr())
	go func() {
		for {
			gateway.PrintOnlineServices()
			time.Sleep(4 * time.Second)
		}
	}()
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
