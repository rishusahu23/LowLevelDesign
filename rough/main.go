package main

import (
	"fmt"
	"sync"
)

// Key is a key in the key-value store
type Key string

// Value is a value associated with a key in the key-value store
type Value string

// Request represents a read or write request
type Request interface {
	Process(store *KeyValueStore)
}

// ReadRequest represents a read request
type ReadRequest struct {
	Key    Key
	Result chan<- Value
}

// Process executes the read request
func (r *ReadRequest) Process(store *KeyValueStore) {
	store.Read(r)
}

// WriteRequest represents a write request
type WriteRequest struct {
	Key   Key
	Value Value
}

// Process executes the write request
func (w *WriteRequest) Process(store *KeyValueStore) {
	store.Write(w)
}

// KeyValueStore represents the key-value store
type KeyValueStore struct {
	data  map[Key]Value
	queue chan Request
	lock  sync.Mutex
}

// NewKeyValueStore creates a new instance of KeyValueStore
func NewKeyValueStore() *KeyValueStore {
	kv := &KeyValueStore{
		data:  make(map[Key]Value),
		queue: make(chan Request),
	}
	go kv.processQueue()
	return kv
}

// Set sets a value for the given key
func (kv *KeyValueStore) Set(key Key, value Value) {
	kv.queue <- &WriteRequest{Key: key, Value: value}
}

// Get retrieves the value for the given key
func (kv *KeyValueStore) Get(key Key) Value {
	result := make(chan Value)
	kv.queue <- &ReadRequest{Key: key, Result: result}
	return <-result
}

// processQueue processes incoming read and write requests
func (kv *KeyValueStore) processQueue() {
	for req := range kv.queue {
		req.Process(kv)
	}
}

// Read processes a read request
func (kv *KeyValueStore) Read(req *ReadRequest) {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	req.Result <- kv.data[req.Key]
}

// Write processes a write request
func (kv *KeyValueStore) Write(req *WriteRequest) {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	kv.data[req.Key] = req.Value
}

func main() {
	kv := NewKeyValueStore()

	// Simulate concurrent read and write requests
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 5; i++ {
		go func(i int) {
			defer wg.Done()
			key := Key(fmt.Sprintf("key%d", i))
			kv.Set(key, Value(fmt.Sprintf("value%d", i)))
			fmt.Printf("Set: %s\n", key)
		}(i)
	}

	for i := 0; i < 5; i++ {
		go func(i int) {
			defer wg.Done()
			key := Key(fmt.Sprintf("key%d", i))
			fmt.Printf("Get: %s = %s\n", key, kv.Get(key))
		}(i)
	}

	wg.Wait()
}
