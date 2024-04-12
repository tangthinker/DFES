package main

import (
	"flag"
	dataServer "github.com/shanliao420/DFES/data-server"
	dataServerPB "github.com/shanliao420/DFES/data-server/proto"
	"github.com/shanliao420/DFES/gateway"
	gatewayPB "github.com/shanliao420/DFES/gateway/proto"
	"github.com/shanliao420/DFES/utils"
	"google.golang.org/grpc"
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
	dataServer.Init(*serverName)
	utils.StartGrpcServer(*host+":"+*port, func(server *grpc.Server) {
		dataServerPB.RegisterDataServiceServer(server, dataServer.RpcServer{})
	})
}
