package main

import (
	"sync"
	"time"
)

type IRateLimiter interface {
	Allow(clientId string) bool
}

type ILimiterStrategy interface {
	Allow(clientId string) bool
}

type Config struct {
	MaxRequest int
	Window     time.Duration
}

type FixedWindow struct {
	RequestCount map[string]int
	TimeWindow   map[string]time.Time
	ConfigMap    map[string]*Config
	mu           sync.Mutex
}

func NewFixedWindow() *FixedWindow {
	return &FixedWindow{
		RequestCount: make(map[string]int),
		TimeWindow:   make(map[string]time.Time),
		ConfigMap:    make(map[string]*Config),
	}
}

func (f *FixedWindow) IsAllow(clientId string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	ts, ok := f.TimeWindow[clientId]
	if !ok || time.Now().Sub(ts) > f.ConfigMap[clientId].Window {
		f.TimeWindow[clientId] = time.Now()
		f.RequestCount[clientId] = 1
		return true
	}

	if f.RequestCount[clientId] < f.ConfigMap[clientId].MaxRequest {
		f.RequestCount[clientId]++
		return true
	}
	return false
}

type SlidingWindow struct {
	TimeWindow map[string][]time.Time
	ConfigMap  map[string]*Config
	mu         sync.Mutex
}

func NewSlidingWindow() *SlidingWindow {
	return &SlidingWindow{
		TimeWindow: make(map[string][]time.Time),
		ConfigMap:  make(map[string]*Config),
	}
}

func (s *SlidingWindow) IsAllow(clientId string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	records, exist := s.TimeWindow[clientId]
	if !exist {
		s.TimeWindow[clientId] = []time.Time{time.Now()}
		return true
	}

	var updatedRecord []time.Time
	for _, ts := range records {
		if now.Sub(ts) <= s.ConfigMap[clientId].Window {
			updatedRecord = append(updatedRecord, ts)
		}
	}
	if len(updatedRecord) < s.ConfigMap[clientId].MaxRequest {
		updatedRecord = append(updatedRecord, time.Now())
		s.TimeWindow[clientId] = updatedRecord
		return true
	}
	s.TimeWindow[clientId] = updatedRecord
	return false
}
