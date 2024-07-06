package cache3

import (
	"sync"
	"time"
)

type ICache interface {
	Insert(key string, value interface{})
	Fetch(key string) (interface{}, bool)
	Remove(key string)
}

type EvictionPolicy interface {
	RecordAccess(key string)
	Evict(mp map[string]interface{}) string
}

type LRUPolicy struct {
	order      []string
	keyToIndex map[string]int
}

func NewLRUPolicy() *LRUPolicy {
	return &LRUPolicy{
		order:      make([]string, 0),
		keyToIndex: make(map[string]int),
	}
}

func (l *LRUPolicy) RecordAccess(key string) {
	if ind, ok := l.keyToIndex[key]; ok {
		l.order = append(l.order[:ind], l.order[ind+1:]...)
	}
	l.order = append(l.order, key)
	l.keyToIndex[key] = len(l.order) - 1
}

func (l *LRUPolicy) Evict(mp map[string]interface{}) string {
	evictKey := l.order[0]
	l.order = l.order[1:]
	delete(l.keyToIndex, evictKey)
	return evictKey
}

type TimeEvictionPolicy struct {
	timestamps map[string]time.Time
	timeout    time.Duration
}

func NewTimeEvictionPolicy(timeout time.Duration) *TimeEvictionPolicy {
	return &TimeEvictionPolicy{
		timestamps: make(map[string]time.Time),
		timeout:    timeout,
	}
}

func (t *TimeEvictionPolicy) RecordAccess(key string) {
	t.timestamps[key] = time.Now()
}

func (t *TimeEvictionPolicy) Evict(mp map[string]interface{}) string {
	oldestKey := ""
	oldestTime := time.Now()
	for key, timee := range t.timestamps {
		if time.Since(timee) > t.timeout {
			oldestKey = key
			break
		}
		if timee.Before(oldestTime) {
			oldestKey = key
			oldestTime = timee
		}
	}
	if oldestKey != "" {
		delete(t.timestamps, oldestKey)
	}
	return oldestKey
}

type Cache struct {
	data           map[string]interface{}
	evictionPolicy EvictionPolicy
	mu             sync.Mutex
}

func NewCache(evictionPolicy EvictionPolicy) *Cache {
	return &Cache{
		evictionPolicy: evictionPolicy,
		data:           make(map[string]interface{}),
	}
}

func (c *Cache) Insert(key string, val interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.data[key]; ok {
		if len(c.data) >= 10 {
			evictedKey := c.evictionPolicy.Evict(c.data)
			delete(c.data, evictedKey)
		}
	}
	c.data[key] = val
	c.evictionPolicy.RecordAccess(key)
}

func (c *Cache) Fetch(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.data[key]
	if ok {
		c.evictionPolicy.RecordAccess(key)
	}
	return val, ok
}

func (c *Cache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

func main() {
	lruCache := NewCache(NewLRUPolicy())
	timeCache := NewCache(NewTimeEvictionPolicy(5 * time.Second))

	lruCache.Insert("q", 1)
	timeCache.Insert("a", 2)
}