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
	onlineService                map[ServiceName]*RegisterInfo
	mutex                        *sync.RWMutex
	continuouslyIncreasingNumber int64
}

func (r *registry) Register(registerInfo RegisterInfo) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	timer := time.NewTimer(heartbeatTime * time.Second)
	registerInfo.HeartbeatTimer = timer
	r.onlineService[registerInfo.ServiceName] = &registerInfo
	log.Println("register", registerInfo.ServiceName, "successful")
	r.continuouslyIncreasingNumber++
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
	registerInfo, ok := r.onlineService[serviceName]
	if !ok {
		log.Println("unregister a service already removed")
		return nil
	}
	delete(r.onlineService, serviceName)
	registerInfo.HeartbeatTimer.Stop()
	return nil
}

func (r *registry) Heartbeat(serviceName ServiceName) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	registerInfo, ok := r.onlineService[serviceName]
	if !ok {
		log.Println("heartbeat get info error")
		return errors.New("heartbeat get info error")
	}
	registerInfo.HeartbeatTimer.Reset(heartbeatTime * time.Second)
	//fmt.Println(registerInfo, "heartbeat...")
	return nil
}
