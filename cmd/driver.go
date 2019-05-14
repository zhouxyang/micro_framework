package cmd

import (
	"fmt"
	"google.golang.org/grpc"
	"sync"
)

type RegisterFunc func(grpcServer *grpc.Server) error

var (
	serverDriversMu sync.RWMutex
	// ServerDrivers the login drivers
	ServerDrivers = make(map[string]RegisterFunc)
)

// RegisterRegisterFunc register a new driver for ip query
func PutRegisterFunc(name string, driver RegisterFunc) {
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

// GetRegisterFunc get the ip query func by pay_method
func GetRegisterFunc(name string) (RegisterFunc, bool) {
	serverDriversMu.RLock()
	defer serverDriversMu.RUnlock()

	driveri, exist := ServerDrivers[name]
	return driveri, exist
}
