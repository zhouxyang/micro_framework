package cmd

import (
	"fmt"
	"micro_framework/configfile"
	"sync"

	"google.golang.org/grpc"
)

// InitServiceFunc 服务初始化函数
type InitServiceFunc func(grpcServer *grpc.Server, config *configfile.Config) error

var (
	serverDriversMu sync.RWMutex
	// ServerDrivers the login drivers
	ServerDrivers = make(map[string]InitServiceFunc)
)

// RegisterService register a new driver for service
func RegisterService(name string, driver InitServiceFunc) {
	serverDriversMu.Lock()
	defer serverDriversMu.Unlock()

	if driver == nil {
		panic(fmt.Sprintf("Register Driver %s is nil", name))
	}

	if _, dup := ServerDrivers[name]; dup {
		panic(fmt.Sprintf("Register Called Twice for Driver %s", name))
	}

	ServerDrivers[name] = driver
}

// GetService get the initService func by service name
func GetService(name string) (InitServiceFunc, bool) {
	serverDriversMu.RLock()
	defer serverDriversMu.RUnlock()

	driveri, exist := ServerDrivers[name]
	return driveri, exist
}
