package gateway

import "time"

type Registry interface {
	Register(registerInfo RegisterInfo) error
	UnRegister(serviceName ServiceName) error
	Heartbeat(serviceName ServiceName) error
}

type RegisterInfo struct {
	ServiceName       ServiceName
	ServiceType       ServiceType
	ServiceAddress    ServiceAddress
	ServiceInterfaces []ServiceInterface
	HeartbeatAddress  string
	HeartbeatTimer    *time.Timer
}

type ServiceName string

type Protocol string

type ServiceType string

const (
	// 	service type
	MateService = "MateService"
	DateService = "DateService"
)

type ServiceAddress struct {
	Host string
	Port string
}

type ServiceInterface struct {
	Path     string // Path or RPC method name
	Protocol Protocol
}
