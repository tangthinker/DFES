package mate_server

import (
	"context"
	pb "github.com/shanliao420/DFES/mate-server/proto"
	"io"
	"log"
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
	if b == nil {
		return &pb.GetResponse{
			GetResult: false,
			Code:      pb.MateCode_FileNotExist,
		}, nil
	}
	return &pb.GetResponse{
		Data:      b,
		GetResult: true,
		Code:      pb.MateCode_Success,
	}, nil
}

func (RpcServer) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if !mateServer.IsLeader() {
		return &pb.DeleteResponse{
			DeleteResult:         false,
			Code:                 pb.MateCode_NotLeader,
			LeaderMateServerAddr: mateServer.leaderRpcAddr,
		}, nil
	}
	ret, err := mateServer.Delete(ctx, in.GetDataId())
	if err != nil {
		return nil, err
	}
	if !ret {
		return &pb.DeleteResponse{
			DeleteResult: false,
			Code:         pb.MateCode_FileNotExist,
		}, nil
	}
	return &pb.DeleteResponse{
		DeleteResult: ret,
		Code:         pb.MateCode_Success,
	}, nil
}

func (RpcServer) PushStream(stream pb.MateService_PushStreamServer) error {
	if !mateServer.IsLeader() {
		_ = stream.SendAndClose(&pb.PushResponse{
			PushResult:           false,
			Code:                 pb.MateCode_NotLeader,
			LeaderMateServerAddr: mateServer.leaderRpcAddr,
		})
		return nil
	}
	r, w := io.Pipe()
	idCh := make(chan string)
	errCh := make(chan error)
	go func() {
		id, err := mateServer.PushStream(stream.Context(), r)
		if err != nil {
			errCh <- err
		}
		idCh <- id
	}()
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			_ = w.Close()
			select {
			case id := <-idCh:
				_ = stream.SendAndClose(&pb.PushResponse{
					PushResult: true,
					DataId:     id,
					Code:       pb.MateCode_Success,
				})
				return nil
			case err := <-errCh:
				return err
			}
		}
		if err != nil {
			return err
		}
		log.Println("receive data")
		_, _ = w.Write(req.Data)
	}
}

func (RpcServer) GetStream(in *pb.GetRequest, stream pb.MateService_GetStreamServer) error {
	r, err := mateServer.GetStream(stream.Context(), in.DataId)
	if err != nil {
		return err
	}
	if r == nil {
		_ = stream.Send(&pb.GetResponse{
			GetResult: false,
			Code:      pb.MateCode_FileNotExist,
		})
		return nil
	}
	for {
		var buff []byte
		n, err := r.Read(buff)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			_ = stream.Send(&pb.GetResponse{
				GetResult: false,
				Code:      pb.MateCode_Fail,
			})
			return err
		}
		_ = stream.Send(&pb.GetResponse{
			GetResult: true,
			Data:      buff[:n],
			Code:      pb.MateCode_Success,
		})
	}
}
