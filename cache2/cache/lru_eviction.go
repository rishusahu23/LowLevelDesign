// lru_eviction.go
package caches

type LRUEviction struct{}

func (l *LRUEviction) Evict(c Cache) {
	if lruCache, ok := c.(*LRUCache); ok {
		elem := lruCache.order.Back()
		if elem != nil {
			lruCache.Remove(elem.Value.(*item).key)
		}
	}
}

func (l *LRUEviction) OnAdd(c Cache, key string) {
	// No-op for LRU
}

func (l *LRUEviction) OnAccess(c Cache, key string) {
	// No-op for LRU
}

func (l *LRUEviction) OnRemove(c Cache, key string) {
	// No-op for LRU
}
