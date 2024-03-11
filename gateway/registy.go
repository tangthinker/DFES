package gateway

import (
	"context"
	"fmt"
	"sync"
	pb "tangthinker.work/DFES/gateway/proto"
	"time"
)

const (
	heartbeatTime = 10 // 10s
)

type registry struct {
	onlineService []RegisterInfo
	mutex         *sync.RWMutex
}

var registryStore = &registry{
	onlineService: make([]RegisterInfo, 0),
	mutex:         &sync.RWMutex{},
}

func GetProvideService(serviceType ServiceType) *RegisterInfo {
	provideServices := getProvideServices(serviceType)
	if len(provideServices) == 0 {
		return nil
	}
	return provideServices[0]
}

func PrintOnlineServices() {
	registryStore.mutex.RLock()
	defer registryStore.mutex.RUnlock()
	for _, item := range registryStore.onlineService {
		fmt.Printf("%s %s %s:%s\n", item.ServiceName, item.ServiceType, item.ServiceAddress.Host, &item.ServiceAddress.Port)
	}
}

func getProvideServices(serviceType ServiceType) []*RegisterInfo {
	result := make([]*RegisterInfo, 1)
	for _, item := range registryStore.onlineService {
		if item.ServiceType == serviceType {
			result = append(result, &item)
		}
	}
	return result
}

func getByServiceName(serviceName ServiceName) *RegisterInfo {
	for _, item := range registryStore.onlineService {
		if item.ServiceName == serviceName {
			return &item
		}
	}
	return nil
}

func (r *registry) Register(registerInfo RegisterInfo) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	timer := time.NewTimer(heartbeatTime * time.Second)
	registerInfo.HeartbeatTimer = timer
	r.onlineService = append(r.onlineService, registerInfo)
	go func() {
		select {
		case <-timer.C:
			_ = r.UnRegister(registerInfo.ServiceName)
		}
	}()
	return nil
}

func (r *registry) UnRegister(serviceName ServiceName) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for i, info := range r.onlineService {
		if info.ServiceName == serviceName {
			newValue := append(r.onlineService[:i], r.onlineService[i+1:]...)
			r.onlineService = newValue
		}
	}
	return nil
}

func (r *registry) Heartbeat(serviceName ServiceName) error {
	r.mutex.RLock()
	defer r.mutex.Unlock()
	registerInfo := getByServiceName(serviceName)
	registerInfo.HeartbeatTimer.Reset(heartbeatTime * time.Second)
	return nil
}

type RpcServer struct {
	pb.UnimplementedRegistryServer
}

func (s *RpcServer) Register(ctx context.Context, in *pb.RegisterInfo) (*pb.RegisterRes, error) {
	err := registryStore.Register(transRpcInfo2RegisterInfo(in))
	return nil, err
}
func (s *RpcServer) UnRegister(ctx context.Context, in *pb.UnRegisterInfo) (*pb.UnRegisterRes, error) {
	err := registryStore.UnRegister(ServiceName(in.ServiceName))
	return nil, err
}
func (s *RpcServer) Heartbeat(ctx context.Context, in *pb.HeartbeatInfo) (*pb.HeartbeatResp, error) {
	err := registryStore.Heartbeat(ServiceName(in.ServiceName))
	return nil, err
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

func transRpcInterfaces2RegisterInterfaces(rpcInterfaces []*pb.ServiceInterface) []ServiceInterface {
	result := make([]ServiceInterface, len(rpcInterfaces))
	for _, item := range rpcInterfaces {
		result = append(result, ServiceInterface{
			Path:     item.Path,
			Protocol: Protocol(item.Protocol),
		})
	}
	return result
}
