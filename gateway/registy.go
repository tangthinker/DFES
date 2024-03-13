package gateway

import (
	"errors"
	"log"
	"sync"
	"time"
)

const (
	heartbeatTime             = 10 // 10s
	DefaultRegistryServerAddr = ":6001"
)

type registry struct {
	onlineService map[ServiceName]*RegisterInfo
	mutex         *sync.RWMutex
}

func (r *registry) Register(registerInfo RegisterInfo) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	timer := time.NewTimer(heartbeatTime * time.Second)
	registerInfo.HeartbeatTimer = timer
	r.onlineService[registerInfo.ServiceName] = &registerInfo
	go func() {
		for {
			select {
			case <-timer.C:
				_ = r.UnRegister(registerInfo.ServiceName)
				log.Println("UnRegister service:", registerInfo)
				return
			}
		}
	}()
	return nil
}

func (r *registry) UnRegister(serviceName ServiceName) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	registerInfo, ok := GetByServiceName(serviceName)
	if ok {
		return nil
	}
	delete(r.onlineService, serviceName)
	registerInfo.HeartbeatTimer.Stop()
	return nil
}

func (r *registry) Heartbeat(serviceName ServiceName) error {
	r.mutex.Lock()
	registerInfo, ok := GetByServiceName(serviceName)
	r.mutex.Unlock()
	if !ok {
		log.Println("heartbeat get info error")
		return errors.New("heartbeat get info error")
	}
	registerInfo.HeartbeatTimer.Reset(heartbeatTime * time.Second)
	//fmt.Println(registerInfo, "heartbeat...")
	return nil
}
