package utils

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	dataServer "tangthinker.work/DFES/data-server/proto"
	gateway "tangthinker.work/DFES/gateway/proto"
	mateServerPB "tangthinker.work/DFES/mate-server/proto"
	"time"
)

const (
	DialTimeout = 30 // 30s
)

func NewGrpcClient(addr string) *grpc.ClientConn {
	ctx, cal := context.WithTimeoutCause(context.Background(), DialTimeout*time.Second, errors.New("GrpcDial timeout"))
	defer cal()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	cnt := 0
	for ; err != nil && cnt < 10; cnt++ {
		conn, err = grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	if err != nil {
		panic("new grpc client error, retry 10 times error")
	}
	return conn
}

func NewRegistryClient(addr string) gateway.RegistryClient {
	return gateway.NewRegistryClient(NewGrpcClient(addr))
}

func NewDataServerClient(addr string) dataServer.DataServiceClient {
	return dataServer.NewDataServiceClient(NewGrpcClient(addr))
}

func NewMateServerClient(addr string) mateServerPB.MateServiceClient {
	return mateServerPB.NewMateServiceClient(NewGrpcClient(addr))
}

func StartGrpcServer(addr string, registerFunc func(server *grpc.Server)) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	srv := grpc.NewServer()
	registerFunc(srv)
	log.Println("start grpc server at ", addr)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
