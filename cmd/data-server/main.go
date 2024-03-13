package main

import (
	"flag"
	"google.golang.org/grpc"
	dataServer "tangthinker.work/DFES/data-server"
	dataServerPB "tangthinker.work/DFES/data-server/proto"
	"tangthinker.work/DFES/gateway"
	gatewayPB "tangthinker.work/DFES/gateway/proto"
	"tangthinker.work/DFES/utils"
)

var (
	host       = flag.String("host", "127.0.0.1", "the host to start server")
	port       = flag.String("port", "8001", "the port to start server")
	serverName = flag.String("server-name", "data-node-1", "Name Server")
)

func main() {
	flag.Parse()
	conn := utils.NewGrpcClient(gateway.DefaultRegistryServerAddr)
	utils.RegisterServer(conn, &gatewayPB.RegisterInfo{
		ServiceName: *serverName,
		ServiceAddress: &gatewayPB.ServiceAddress{
			Host: *host,
			Port: *port,
		},
		ServiceType:       gateway.DateService,
		ServiceInterfaces: make([]*gatewayPB.ServiceInterface, 0),
		HeartbeatAddress:  "",
	})
	dataServer.SetDataServerName(*serverName)
	dataServer.Init()
	utils.StartGrpcServer(*host+":"+*port, func(server *grpc.Server) {
		dataServerPB.RegisterDataServiceServer(server, dataServer.RpcServer{})
	})
}
