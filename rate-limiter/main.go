package main

import (
	"fmt"
	builder2 "github.com/rishu/design/rate-limiter/builder"
	director2 "github.com/rishu/design/rate-limiter/director"
	"time"
)

func main() {
	builder := builder2.NewTokenBucketBuilder()
	director := director2.NewRateLimiterDirector(builder)
	rateLimiter := director.Construct()
	clientID := "client-1"
	for i := 0; i < 15; i++ {
		if rateLimiter.IsRequestAllowed(clientID) {
			fmt.Printf("Request %d from %s allowed\n", i+1, clientID)
		} else {
			fmt.Printf("Request %d from %s denied\n", i+1, clientID)
		}
		time.Sleep(1 * time.Millisecond)
	}
}
