package gateway

import (
	"context"
	"errors"
	pb "github.com/shanliao420/DFES/gateway/proto"
	"log"
)

type RpcServer struct {
	pb.UnimplementedRegistryServer
}

func (s *RpcServer) Register(ctx context.Context, in *pb.RegisterInfo) (*pb.RegisterRes, error) {
	err := registryStore.Register(transRpcInfo2RegisterInfo(in))
	if err != nil {
		return nil, err
	}
	return &pb.RegisterRes{
		RegisterMessage: "Register successful",
		RegisterResult:  true,
	}, nil
}
func (s *RpcServer) UnRegister(ctx context.Context, in *pb.UnRegisterInfo) (*pb.UnRegisterRes, error) {
	err := registryStore.UnRegister(ServiceName(in.ServiceName))
	if err != nil {
		return nil, err
	}
	return &pb.UnRegisterRes{
		UnRegisterMessage: "UnRegister successful",
		UnRegisterResult:  true,
	}, nil
}
func (s *RpcServer) Heartbeat(ctx context.Context, in *pb.HeartbeatInfo) (*pb.HeartbeatResp, error) {
	err := registryStore.Heartbeat(ServiceName(in.ServiceName))
	if err != nil {
		return &pb.HeartbeatResp{
			HeartBeatResult: false,
		}, err
	}
	return &pb.HeartbeatResp{
		HeartBeatResult: true,
	}, nil
}

func (s *RpcServer) GetProvideServices(ctx context.Context, in *pb.GetProvideInfo) (*pb.GetProvidesResp, error) {
	providers := GetProvideServices(ServiceType(in.GetServiceType()))
	return &pb.GetProvidesResp{
		GetResult:       true,
		ProvideServices: transRegisterInfos2RpcInfos(providers),
	}, nil

}

func (s *RpcServer) GetProvideService(ctx context.Context, in *pb.GetProvideInfo) (*pb.GetProvideResp, error) {
	provider, err := GetProvideService(ServiceType(in.ServiceType))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.GetProvideResp{
		GetResult:      true,
		ProvideService: transRegisterInfo2RpcInfo(provider),
	}, nil
}

//func (s *RpcServer) GetLeaderMate(ctx context.Context, in *pb.Empty) (*pb.GetProvideResp, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method GetLeaderMate not implemented")
//}

func (s *RpcServer) GetProvideByName(ctx context.Context, in *pb.GetByNameInfo) (*pb.GetProvideResp, error) {
	provider, ok := GetByServiceName(ServiceName(in.ServiceName))
	if !ok {
		return nil, errors.New("get provider empty")
	}
	return &pb.GetProvideResp{
		GetResult:      true,
		ProvideService: transRegisterInfo2RpcInfo(provider),
	}, nil
}

func transRpcInfo2RegisterInfo(rpcInfo *pb.RegisterInfo) RegisterInfo {
	return RegisterInfo{
		ServiceName: ServiceName(rpcInfo.ServiceName),
		ServiceType: ServiceType(rpcInfo.ServiceType),
		ServiceAddress: ServiceAddress{
			Host: rpcInfo.ServiceAddress.Host,
			Port: rpcInfo.ServiceAddress.Port,
		},
		ServiceInterfaces: transRpcInterfaces2RegisterInterfaces(rpcInfo.ServiceInterfaces),
		HeartbeatAddress:  rpcInfo.HeartbeatAddress,
	}
}

func transRegisterInfos2RpcInfos(info []RegisterInfo) []*pb.RegisterInfo {
	result := make([]*pb.RegisterInfo, 0)
	for _, item := range info {
		result = append(result, transRegisterInfo2RpcInfo(item))
	}
	return result
}

func transRegisterInfo2RpcInfo(info RegisterInfo) *pb.RegisterInfo {
	return &pb.RegisterInfo{
		ServiceName: string(info.ServiceName),
		ServiceType: string(info.ServiceType),
		ServiceAddress: &pb.ServiceAddress{
			Host: info.ServiceAddress.Host,
			Port: info.ServiceAddress.Port,
		},
		ServiceInterfaces: transRegisterInterfaces2RpcInterfaces(info.ServiceInterfaces),
		HeartbeatAddress:  info.HeartbeatAddress,
	}
}

func transRegisterInterfaces2RpcInterfaces(serviceInterfaces []ServiceInterface) []*pb.ServiceInterface {
	result := make([]*pb.ServiceInterface, 0)
	for _, item := range serviceInterfaces {
		result = append(result, &pb.ServiceInterface{
			Path:     item.Path,
			Protocol: string(item.Protocol),
		})
	}
	return result
}

func transRpcInterfaces2RegisterInterfaces(rpcInterfaces []*pb.ServiceInterface) []ServiceInterface {
	result := make([]ServiceInterface, 0)
	for _, item := range rpcInterfaces {
		result = append(result, ServiceInterface{
			Path:     item.Path,
			Protocol: Protocol(item.Protocol),
		})
	}
	return result
}
