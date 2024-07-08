package utils

import (
	"context"
	pb "github.com/tangthinker/DFES/gateway/proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	ClientHeartbeatTime = 1 // 5s
	RetryTimes          = 3
	RetryTime           = 5 * time.Second
)

func RegisterServer(conn *grpc.ClientConn, info *pb.RegisterInfo) {
	client := pb.NewRegistryClient(conn)
	res, err := client.Register(context.Background(), info)
	if err != nil || !res.RegisterResult {
		log.Fatalln("register server", info.ServiceName, " error:", err)
	}
	timer := time.NewTimer(ClientHeartbeatTime * time.Second)
	go func() {
		for {
			select {
			case <-timer.C:
				resp, err := client.Heartbeat(context.Background(), &pb.HeartbeatInfo{
					ServiceName: info.ServiceName,
				})
				if err != nil || !resp.HeartBeatResult {
					for i := 0; i < RetryTimes; i++ {
						resp, err := client.Register(context.Background(), info)
						if err == nil && resp.RegisterResult {
							log.Println("ReRegister successful")
							break
						}
						time.Sleep(RetryTime)
					}
				}
				timer.Reset(ClientHeartbeatTime * time.Second)
			}
		}
	}()
	log.Println("register server", info.ServiceName, " successful")
}
