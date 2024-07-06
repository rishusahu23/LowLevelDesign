package inmemory

import (
	"sync"
	"time"
)

type Bucket struct {
	Tokens     int
	LastRefill time.Time
}

type MemoryDB struct {
	Buckets map[string]*Bucket
	mu      sync.RWMutex
}

func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		Buckets: make(map[string]*Bucket),
	}
}

func (m *MemoryDB) GetBucket(apiKey string) *Bucket {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Buckets[apiKey]
}

func (m *MemoryDB) SetBucket(apiKey string, bucket *Bucket) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Buckets[apiKey] = bucket
}
