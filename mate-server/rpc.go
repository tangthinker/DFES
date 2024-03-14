package mate_server

import (
	"context"
	pb "github.com/shanliao420/DFES/mate-server/proto"
)

type RpcServer struct {
	pb.UnimplementedMateServiceServer
}

func (RpcServer) Join(ctx context.Context, in *pb.JoinRequest) (*pb.JoinResponse, error) {
	err := mateServer.Join(in.GetServerName(), in.GetServerAddr())
	if err != nil {
		return nil, err
	}
	return &pb.JoinResponse{
		JoinResult: true,
	}, nil
}
func (RpcServer) Push(ctx context.Context, in *pb.PushRequest) (*pb.PushResponse, error) {
	if !mateServer.IsLeader() {
		return &pb.PushResponse{
			PushResult:           false,
			Code:                 pb.MateCode_NotLeader,
			LeaderMateServerAddr: mateServer.leaderRpcAddr,
		}, nil
	}
	fileMateId, err := mateServer.Push(ctx, in.Data)
	if err != nil {
		return nil, err
	}
	return &pb.PushResponse{
		PushResult: true,
		Code:       pb.MateCode_Success,
		DataId:     fileMateId,
	}, nil
}

func (RpcServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	b, err := mateServer.Get(ctx, in.GetDataId())
	if err != nil {
		return nil, err
	}
	return &pb.GetResponse{
		Data: b,
	}, nil
}
