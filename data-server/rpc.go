package data_server

import (
	"context"
	pb "tangthinker.work/DFES/data-server/proto"
)

type RpcServer struct {
	pb.UnimplementedDataServiceServer
}

func (RpcServer) Push(ctx context.Context, in *pb.PushRequest) (*pb.PushResponse, error) {
	fragmentId := dataService.Push(in.GetFragmentData())
	return &pb.PushResponse{
		FragmentId:  fragmentId,
		ServiceName: dataService.serverName,
		PushResult:  true,
	}, nil
}
func (RpcServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	data := dataService.Get(in.GetFragmentId())
	return &pb.GetResponse{
		FragmentData: data,
		GetResult:    true,
	}, nil
}
func (RpcServer) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	ret := dataService.Delete(in.GetFragmentId())
	return &pb.DeleteResponse{
		DeleteResult: ret,
	}, nil
}
