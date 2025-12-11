package drivers

import (
	"fmt"
	"sync"

	"github.com/wentf9/MyGoFileHub/internal/domain/vfs"
)

// DriverFactory 定义创建驱动的函数签名
type DriverFactory func() vfs.StorageDriver

var (
	driverFactories = make(map[string]DriverFactory)
	mu              sync.RWMutex
)

// Register 注册一个新的驱动类型 (在各驱动的 init() 函数中调用)
// 比如：drivers.Register("smb", NewSMBDriver)
func Register(name string, factory DriverFactory) {
	mu.Lock()
	defer mu.Unlock()
	if factory == nil {
		panic("driver factory is nil")
	}
	driverFactories[name] = factory
}

// CreateInstance 根据类型名称创建一个新的驱动实例
func CreateInstance(driverType string) (vfs.StorageDriver, error) {
	mu.RLock()
	defer mu.RUnlock()
	factory, ok := driverFactories[driverType]
	if !ok {
		return nil, fmt.Errorf("unknown driver type: %s", driverType)
	}
	return factory(), nil
}
