package main

import (
	"errors"
	"sync"
)

type Memory struct {
	db map[string]string
	mu sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		db: make(map[string]string),
	}
}

func (m *Memory) Store(key, val string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.db[key] = val
	return nil
}

func (m *Memory) Recall(key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	val, ok := m.db[key]
	if !ok {
		return "", errors.New("not found")
	}
	return val, nil
}
