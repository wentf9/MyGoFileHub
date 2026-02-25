package drivers

import (
	"fmt"
	"sync"

	"github.com/wentf9/MyGoFileHub/internal/domain/model"
	"github.com/wentf9/MyGoFileHub/internal/domain/vfs"
)

// DriverFactory 定义创建驱动的函数签名
type DriverFactory func() vfs.StorageDriver

type driverEntry struct {
	factory DriverFactory
	schema  model.StorageDriverSchema
}

var (
	driverFactories = make(map[string]driverEntry)
	mu              sync.RWMutex
)

// Register 注册一个新的驱动类型 (在各驱动的 init() 函数中调用)
// 比如：drivers.Register("smb", NewSMBDriver, schema)
func Register(name string, factory DriverFactory, schema model.StorageDriverSchema) {
	mu.Lock()
	defer mu.Unlock()
	if factory == nil {
		panic("driver factory is nil")
	}
	driverFactories[name] = driverEntry{
		factory: factory,
		schema:  schema,
	}
}

// CreateInstance 根据类型名称创建一个新的驱动实例
func CreateInstance(driverType string) (vfs.StorageDriver, error) {
	mu.RLock()
	defer mu.RUnlock()
	entry, ok := driverFactories[driverType]
	if !ok {
		return nil, fmt.Errorf("unknown driver type: %s", driverType)
	}
	return entry.factory(), nil
}

// GetRegisteredSchemas 返回所有已注册驱动的描述结构
func GetRegisteredSchemas() []model.StorageDriverSchema {
	mu.RLock()
	defer mu.RUnlock()
	schemas := make([]model.StorageDriverSchema, 0, len(driverFactories))
	for _, entry := range driverFactories {
		schemas = append(schemas, entry.schema)
	}
	return schemas
}
