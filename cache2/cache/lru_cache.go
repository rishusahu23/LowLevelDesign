package caches

import (
	"container/list"
	"sync"
)

type item struct {
	key string
	val interface{}
}

type LRUCache struct {
	cap      int
	items    map[string]*list.Element
	order    *list.List
	strategy EvictionStrategy
	mu       sync.Mutex
}

func NewLRUCache(cap int, strategy EvictionStrategy) *LRUCache {
	return &LRUCache{
		cap:      cap,
		items:    make(map[string]*list.Element),
		order:    list.New(),
		strategy: strategy,
	}
}

func (c *LRUCache) Set(key string, val interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.items) >= c.cap {
		c.strategy.Evict(c)
	}

	if elem, found := c.items[key]; found {
		c.order.MoveToFront(elem)
		elem.Value.(*item).val = val
	} else {
		elem = c.order.PushFront(&item{key, val})
		c.items[key] = elem
	}

	c.strategy.OnAdd(c, key)
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.items[key]; found {
		c.order.MoveToFront(elem)
		c.strategy.OnAccess(c, key)
		return elem.Value.(*item).val, true
	}
	return nil, false
}

func (c *LRUCache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.items[key]; found {
		delete(c.items, key)
		c.order.Remove(elem)
		c.strategy.OnRemove(c, key)
	}
}
