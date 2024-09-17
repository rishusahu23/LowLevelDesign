package main

import (
	"fmt"
	"sync"
	"time"
)

type IRateLimiter interface {
	Allow(clientId string) bool
}

type IStrategy interface {
	Allow(clientId string) bool
}

type Config struct {
	Request int
	Window  time.Duration
}

type RateLimiter struct {
	strategy IStrategy
}

func NewRateLimiter(strategy IStrategy) *RateLimiter {
	return &RateLimiter{strategy: strategy}
}

func (r *RateLimiter) Allow(clientId string) bool {
	return r.strategy.Allow(clientId)
}

type FixedWindow struct {
	requests  int
	window    time.Duration
	timestamp map[string]time.Time
	count     map[string]int
	mu        sync.Mutex
}

func (f *FixedWindow) Allow(clientId string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	ts, exists := f.timestamp[clientId]
	if !exists || time.Now().Sub(ts) > f.window {
		f.timestamp[clientId] = time.Now()
		f.count[clientId] = 1
		return true
	}
	if f.count[clientId] < f.requests {
		f.count[clientId]++
		return true
	}
	return false
}

func NewFixedWindow(config *Config) *FixedWindow {
	return &FixedWindow{
		requests:  config.Request,
		window:    config.Window,
		timestamp: make(map[string]time.Time),
		count:     make(map[string]int),
	}
}

type SlidingWindow struct {
	requests int
	window   time.Duration
	mu       sync.Mutex
	records  map[string][]time.Time
}

func (s *SlidingWindow) Allow(clientId string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	records, exists := s.records[clientId]
	if !exists {
		s.records[clientId] = []time.Time{now}
		return true
	}
	var updatedRecord []time.Time
	for _, ts := range records {
		if now.Sub(ts) <= s.window {
			updatedRecord = append(updatedRecord, ts)
		}
	}
	if len(updatedRecord) < s.requests {
		updatedRecord = append(updatedRecord, now)
		s.records[clientId] = updatedRecord
		return true
	}
	s.records[clientId] = updatedRecord
	return false
}

func NewSlidingWindow(config *Config) *SlidingWindow {
	return &SlidingWindow{
		requests: config.Request,
		window:   config.Window,
		records:  make(map[string][]time.Time),
	}
}

type TokenBucket struct {
	cap        int
	tokens     map[string]int
	fillRate   float64
	lastAccess map[string]time.Time
	mu         sync.Mutex
}

func (t *TokenBucket) Allow(clientId string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	ls, ok := t.lastAccess[clientId]
	if !ok {
		t.tokens[clientId] = t.cap - 1
		t.lastAccess[clientId] = now
		return true
	}
	elapsed := now.Sub(ls).Seconds()
	newTokens := int(elapsed * t.fillRate)
	if newTokens > 0 {
		t.lastAccess[clientId] = now
		t.tokens[clientId] = min(t.cap, t.tokens[clientId]+newTokens)
	}
	if t.tokens[clientId] > 0 {
		t.tokens[clientId]--
		return true
	}
	return false
}

func NewTokenBucket(config *Config, lastAccess map[string]time.Time, mu sync.Mutex) *TokenBucket {
	return &TokenBucket{
		cap:        config.Request,
		tokens:     make(map[string]int),
		fillRate:   float64(config.Request) / config.Window.Seconds(),
		lastAccess: make(map[string]time.Time),
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	config := &Config{
		Request: 6,
		Window:  5 * time.Second,
	}
	fixed := NewFixedWindow(config)
	rate := NewRateLimiter(fixed)
	fmt.Print(rate.Allow("hell"))
}
