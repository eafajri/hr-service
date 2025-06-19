package cache

import "sync"

type MemoryCache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string)
}

type MemoryCacheImpl struct {
	cache map[string]interface{}
	mu    sync.Mutex
}

func NewMemoryCache() *MemoryCacheImpl {
	return &MemoryCacheImpl{
		cache: make(map[string]interface{}),
	}
}

func (m *MemoryCacheImpl) Get(key string) (interface{}, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	value, exists := m.cache[key]
	return value, exists
}

func (m *MemoryCacheImpl) Set(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache[key] = value
}

func (m *MemoryCacheImpl) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.cache, key)
}
