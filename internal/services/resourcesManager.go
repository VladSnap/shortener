package services

import "sync"

// ResourceManager управляет жизненным циклом объектов, чтобы при завершении программы вызвать очистку
type ResourceManager struct {
	cleanupFuncs []func() error
	mu           sync.Mutex
}

func NewResourceManager() *ResourceManager {
	return &ResourceManager{}
}

// Register регистрирует функцию очистки
func (rm *ResourceManager) Register(cleanupFunc func() error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.cleanupFuncs = append(rm.cleanupFuncs, cleanupFunc)
}

// Cleanup вызывает все зарегистрированные функции очистки
func (rm *ResourceManager) Cleanup() error {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	for _, cleanup := range rm.cleanupFuncs {
		err := cleanup()
		if err != nil {
			return err
		}
	}

	return nil
}
