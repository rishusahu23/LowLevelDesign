package main

import (
	"fmt"
	"sync"
)

type DataStore struct {
	store          map[string]string
	inTransactions bool
	mutex          sync.Mutex
	actions        []func()
}

func newDataStore() *DataStore {
	return &DataStore{
		store:          make(map[string]string),
		inTransactions: false,
		actions:        make([]func(), 0),
	}
}

func (s *DataStore) Get(key string) (string, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	val, ok := s.store[key]
	return val, ok
}

func (s *DataStore) Put(key, val string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.store[key] = val
	if !s.inTransactions {
		return
	}
	s.actions = append(s.actions, func() {
		delete(s.store, key)
	})
}

func (s *DataStore) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	val, ok := s.store[key]
	if !ok {
		return
	}
	delete(s.store, key)
	if !s.inTransactions {
		return
	}
	s.actions = append(s.actions, func() {
		s.store[key] = val
	})
	return
}

func (s *DataStore) Print() {
	for key, val := range s.store {
		fmt.Printf("key: %v, val: %v \n", key, val)
	}
	fmt.Println()
}

func (s *DataStore) beginTransaction() {
	s.mutex.Lock()
	s.mutex.Unlock()
	s.inTransactions = true
}

func (s *DataStore) commit() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.actions = nil
	s.inTransactions = false
}

func (s *DataStore) rollback() {
	s.mutex.Lock()
	s.mutex.Unlock()
	for _, fun := range s.actions {
		fun()
	}
}

func main() {
	ds := newDataStore()
	ds.Put("name", "rishu")
	ds.beginTransaction()
	ds.Delete("name")
	ds.rollback()
	ds.Print()
}
