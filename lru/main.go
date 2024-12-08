package main

import (
	"container/list"
	"fmt"
)

type LRUCache struct {
	cap   int
	cache map[int]*list.Element
	list  *list.List
}

type Entry struct {
	Key   int
	Value int
}

func NewLRUCache(cap int) *LRUCache {
	return &LRUCache{
		cap:   cap,
		cache: make(map[int]*list.Element),
		list:  list.New(),
	}
}

func (l *LRUCache) Get(key int) int {
	if ele, found := l.cache[key]; found {
		l.list.MoveToFront(ele)
		return ele.Value.(*Entry).Value
	}
	return -1
}

func (l *LRUCache) Put(key, val int) {
	if ele, found := l.cache[key]; found {
		ele.Value.(*Entry).Value = val
		l.list.MoveToFront(ele)
	} else {
		if l.list.Len() == l.cap {
			tail := l.list.Back()
			if tail != nil {
				delete(l.cache, tail.Value.(*Entry).Key)
				l.list.Remove(tail)
			}
		}
		newE := &Entry{
			Key:   key,
			Value: val,
		}
		ele = l.list.PushFront(newE)
		l.cache[key] = ele
	}
}

func main() {
	cache := NewLRUCache(2)
	cache.Put(1, 1)
	cache.Put(2, 2)
	fmt.Println(cache.Get(1)) // Output: 1
	cache.Put(3, 3)           // Evicts key 2
	fmt.Println(cache.Get(2)) // Output: -1
	fmt.Println(cache.Get(3)) // Output: 3
	cache.Put(4, 4)           // Evicts key 1
	fmt.Println(cache.Get(1)) // Output: -1
	fmt.Println(cache.Get(3)) // Output: 3
	fmt.Println(cache.Get(4)) // Output: 4
}
