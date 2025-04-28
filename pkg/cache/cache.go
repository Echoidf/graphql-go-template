package cache

import (
	"sync"

	"go.uber.org/zap"
)

type Cache[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex
}

func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		data: make(map[K]V),
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

func (c *Cache[K, V]) Set(key K, val V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = val
}

func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

func (c *Cache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

func (c *Cache[K, V]) Update(key K, updateFn func(*V) error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val := c.data[key]
	if err := updateFn(&val); err != nil {
		zap.L().Error("failed to update cache", zap.Error(err))
		return
	}
	c.data[key] = val
}

func (c *Cache[K, V]) Iterator() <-chan struct {
	Key K
	Val V
} {
	ch := make(chan struct {
		Key K
		Val V
	}, len(c.data))
	go func() {
		c.mu.RLock()
		defer c.mu.RUnlock()
		defer close(ch)

		for k, v := range c.data {
			ch <- struct {
				Key K
				Val V
			}{k, v}
		}
	}()
	return ch
}