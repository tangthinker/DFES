package main

import (
	"flag"
	"github.com/shanliao420/DFES/gateway"
	gatewayPB "github.com/shanliao420/DFES/gateway/proto"
	mateServer "github.com/shanliao420/DFES/mate-server"
	mateServerPB "github.com/shanliao420/DFES/mate-server/proto"
	"github.com/shanliao420/DFES/utils"
	"google.golang.org/grpc"
	"log"
)

var (
	host       = flag.String("host", "127.0.0.1", "the host to start server")
	port       = flag.String("port", "7001", "the port to start server")
	serverName = flag.String("server-name", "mate-node-1", "Name Server")
	leaderAddr = flag.String("leader-addr", "", "leader addr to connect")
	raftAddr   = flag.String("raft-addr", "127.0.0.1:9001", "host to communicat raft")
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
		ServiceType:       gateway.MateService,
		ServiceInterfaces: make([]*gatewayPB.ServiceInterface, 0),
		HeartbeatAddress:  "",
	})
	localAddr := *host + ":" + *port
	mateServer.Init()
	mateServer.SetServerName(*serverName)
	mateServer.SetLocalRpcAddr(localAddr)
	mateServer.SetRaftAddr(*raftAddr)
	mateServer.InitRaft(leaderAddr == nil || *leaderAddr == "")
	if leaderAddr != nil && *leaderAddr != "" {
		if err := mateServer.Join(*leaderAddr); err != nil {
			log.Fatalln("join cluster error:", err)
		}
	}
	utils.StartGrpcServer(localAddr, func(server *grpc.Server) {
		mateServerPB.RegisterMateServiceServer(server, &mateServer.RpcServer{})
	})
}
