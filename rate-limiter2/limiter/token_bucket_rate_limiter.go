package limiter

import (
	"github.com/rishu/design/rate-limiter2/inmemory"
	"sync"
	"time"
)

type TokenBucket struct {
	Rate       int
	BucketSize int
	MemoryDB   *inmemory.MemoryDB
	mu         sync.Mutex
}

func NewTokenBucket(rate, bucketSize int, db *inmemory.MemoryDB) *TokenBucket {
	return &TokenBucket{
		Rate:       rate,
		BucketSize: bucketSize,
		MemoryDB:   db,
	}
}

func (t *TokenBucket) IsRequestAllowed(apiKey string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	bucket := t.MemoryDB.GetBucket(apiKey)
	if bucket == nil {
		bucket = &inmemory.Bucket{
			Tokens:     t.BucketSize,
			LastRefill: time.Now(),
		}
		t.MemoryDB.SetBucket(apiKey, bucket)
	}
	now := time.Now()
	elapsed := now.Sub(bucket.LastRefill).Seconds()
	refillTokens := int(elapsed) * t.Rate
	if refillTokens > 0 {
		bucket.Tokens = min(refillTokens, t.BucketSize)
		bucket.LastRefill = now
		t.MemoryDB.SetBucket(apiKey, bucket)
	}
	if bucket.Tokens > 0 {
		bucket.Tokens--
		t.MemoryDB.SetBucket(apiKey, bucket)
		return true
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
