package gateway

import (
	"errors"
	"fmt"
	"sync"
)

var registryStore = &registry{
	onlineService: make(map[ServiceName]*RegisterInfo),
	mutex:         &sync.RWMutex{},
}

func GetProvideService(serviceType ServiceType) (RegisterInfo, error) {
	provideServices := GetProvideServices(serviceType)
	if len(provideServices) == 0 {
		return RegisterInfo{}, errors.New("register info empty")
	}
	return provideServices[0], nil
}

func PrintOnlineServices() {
	registryStore.mutex.RLock()
	defer registryStore.mutex.RUnlock()
	for _, item := range registryStore.onlineService {
		fmt.Printf("%s %s %s:%s\n", item.ServiceName, item.ServiceType, item.ServiceAddress.Host, &item.ServiceAddress.Port)
	}
}

func GetProvideServices(serviceType ServiceType) []RegisterInfo {
	result := make([]RegisterInfo, 0)
	for _, item := range registryStore.onlineService {
		if item.ServiceType == serviceType {
			result = append(result, *item)
		}
	}
	return result
}

func GetByServiceName(serviceName ServiceName) (RegisterInfo, bool) {
	value, ok := registryStore.onlineService[serviceName]
	if ok {
		return *value, ok
	}
	return RegisterInfo{}, ok
}
