package main

import (
	"fmt"
	"github.com/rishu/design/rate-limiter2/inmemory"
	"github.com/rishu/design/rate-limiter2/limiter"
	"time"
)

func main() {
	memoryDB := inmemory.NewMemoryDB()
	rateLimiter := limiter.NewTokenBucket(1, 10, memoryDB)
	apiKey := "user-123"
	for i := 0; i < 15; i++ {
		if rateLimiter.IsRequestAllowed(apiKey) {
			fmt.Printf("Request %d allowed\n", i+1)
		} else {
			fmt.Printf("Request %d denied\n", i+1)
		}
		time.Sleep(500 * time.Millisecond)
	}
}
