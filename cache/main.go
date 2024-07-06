package main

import (
	"fmt"
	"sync"
)

type CacheItem struct {
	Key  string
	Val  string
	Prev *CacheItem
	Next *CacheItem
}

type LRUCache struct {
	items map[string]*CacheItem
	cap   int
	head  *CacheItem
	tail  *CacheItem
	size  int
	lock  sync.Mutex
}

func NewLRUCache(cap int) *LRUCache {
	return &LRUCache{
		cap:   cap,
		items: make(map[string]*CacheItem, cap),
	}
}

func (c *LRUCache) Get(key string) (string, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if item, found := c.items[key]; found {
		c.moveToFront(item)
		return item.Val, true
	}
	return "", false
}

func (c *LRUCache) Put(key, val string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if item, found := c.items[key]; found {
		item.Val = val
		c.moveToFront(item)
		return
	}
	newItem := &CacheItem{
		Key: key,
		Val: val,
	}
	if c.size == c.cap {
		delete(c.items, c.tail.Key)
		c.removeLast()
	}
	c.items[key] = newItem
	c.addToFront(newItem)
	c.size++
}

func (c *LRUCache) moveToFront(item *CacheItem) {
	if item == c.head {
		return
	}
	if item.Next == nil {
		c.tail = item.Prev
	} else {
		item.Next.Prev = item.Prev
	}
	item.Prev.Next = item.Next
	item.Prev = nil
	item.Next = c.head
	c.head.Prev = item
	c.head = item
	if c.tail == nil {
		c.tail = item
	}
}

func (c *LRUCache) addToFront(item *CacheItem) {
	item.Prev = nil
	item.Next = c.head
	if c.head != nil {
		c.head.Prev = item
	}
	c.head = item
	if c.tail == nil {
		c.tail = item
	}
}

func (c *LRUCache) removeLast() {
	if c.tail != nil {
		c.tail = c.tail.Prev
		if c.tail != nil {
			c.tail.Next = nil
		}
	}
}

func main() {
	cache := NewLRUCache(2)

	cache.Put("one", "1")
	cache.Put("two", "2")
	fmt.Println(cache.Get("one")) // Should print "1 true"

	cache.Put("three", "3") // This should evict "two"
	_, found := cache.Get("two")
	fmt.Println(found) // Should print "false"

	cache.Put("four", "4") // This should evict "one"
	_, found = cache.Get("one")
	fmt.Println(found) // Should print "false"
}
