package main

import (
	"container/list"
	"fmt"
	"sync"
)

type EvictionPolicy interface {
	Evict(cache *Cache)
}

type Item struct {
	Key string
	Val interface{}
}

type LRU struct {
	maxSize int
	items   map[string]*list.Element
	order   *list.List
}

func NewLRU() *LRU {
	return &LRU{
		maxSize: 5,
		items:   make(map[string]*list.Element),
		order:   list.New(),
	}
}

func (l *LRU) Evict(cache *Cache) {
	if l.order.Len() == 0 {
		return
	}
	oldest := l.order.Back()
	if oldest == nil {
		return
	}
	item := oldest.Value.(*Item)
	l.order.Remove(oldest)
	delete(l.items, item.Key)
}

func (l *LRU) Update(item *Item) {
	if e, exist := l.items[item.Key]; exist {
		l.order.Remove(e)
	}
	l.items[item.Key] = l.order.PushFront(item)
}

type FIFO struct {
	maxSize int
	order   []string
}

func NewFifo() *FIFO {
	return &FIFO{
		maxSize: 3,
		order:   make([]string, 0),
	}
}

func (f *FIFO) Evict(c *Cache) {
	if len(f.order) == 0 {
		return
	}
	oldestKey := f.order[0]
	f.order = f.order[1:]
	fmt.Println(oldestKey)
}

func (f *FIFO) AddKey(key string) {
	f.order = append(f.order, key)
}

type Cache struct {
	store          map[string]interface{}
	evictionPolicy EvictionPolicy
	mu             sync.Mutex
	maxSize        int
}

func NewCache(policy EvictionPolicy) *Cache {
	return &Cache{
		store:          make(map[string]interface{}),
		evictionPolicy: policy,
		maxSize:        3,
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.store) >= c.maxSize {
		c.evictionPolicy.Evict(c)
	}

	c.store[key] = value

	switch policy := c.evictionPolicy.(type) {
	case *LRU:
		policy.Update(&Item{
			Key: key,
			Val: value,
		})
	case *FIFO:
		policy.AddKey(key)
	}
}
