package limiter

import (
	"math"
	"sync"
	"time"
)

type TokenBucket struct {
	Rate       int
	BucketSize int
	Tokens     int
	LastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucket(rate, bucketSize int) *TokenBucket {
	return &TokenBucket{
		Rate:       rate,
		BucketSize: bucketSize,
		Tokens:     bucketSize,
		LastRefill: time.Now(),
	}
}

func (t *TokenBucket) IsRequestAllowed(requestId string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(t.LastRefill).Seconds()
	refillToken := int(elapsed) * t.Rate
	if refillToken > 0 {
		t.Tokens = int(math.Min(float64(refillToken), float64(t.BucketSize)))
		t.LastRefill = now
	}
	if t.Tokens > 0 {
		t.Tokens--
		return true
	}
	return false
}
