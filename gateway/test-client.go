package gateway

import (
	"context"
	"google.golang.org/grpc"
	pb "tangthinker.work/DFES/gateway/proto"
)

func RegisterClient(conn *grpc.ClientConn) {
	client := pb.NewRegistryClient(conn)
	client.Register(context.Background(), &pb.RegisterInfo{
		ServiceName: "shanliao",
		ServiceType: MateService,
		ServiceAddress: &pb.ServiceAddress{
			Host: "127.0.0.1",
			Port: "8080",
		},
		ServiceInterfaces: make([]*pb.ServiceInterface, 0),
		HeartbeatAddress:  "127.0.0.1:8080",
	})

}
