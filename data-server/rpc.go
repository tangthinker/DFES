package data_server

import (
	"context"
	pb "github.com/shanliao420/DFES/data-server/proto"
)

type RpcServer struct {
	pb.UnimplementedDataServiceServer
}

func (RpcServer) Push(ctx context.Context, in *pb.PushRequest) (*pb.PushResponse, error) {
	fragmentId, err := dataService.Push(in.GetFragmentData())
	if err != nil {
		return nil, err
	}
	return &pb.PushResponse{
		FragmentId:  fragmentId,
		ServiceName: dataService.serverName,
		PushResult:  true,
	}, nil
}
func (RpcServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	data, err := dataService.Get(in.GetFragmentId())
	if err != nil {
		return nil, err
	}
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
