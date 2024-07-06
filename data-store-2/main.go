package main

import "sync"

type DataStore struct {
	store         map[string]string
	lock          sync.Mutex
	inTransaction bool
	actions       []func()
}

func NewDataStore() *DataStore {
	return &DataStore{
		store:         make(map[string]string),
		inTransaction: false,
		actions:       make([]func(), 0),
	}
}

func (d *DataStore) Get(key string) (string, bool) {
	d.lock.Lock()
	defer d.lock.Unlock()
	va, ok := d.store[key]
	return va, ok
}

func (d *DataStore) Put(key, val string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.store[key] = val
	if !d.inTransaction {
		return
	}
	d.actions = append(d.actions, func() {
		delete(d.store, key)
	})
}

func (d *DataStore) Delete(key string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	val := d.store[key]
	delete(d.store, key)
	if !d.inTransaction {
		return
	}
	d.actions = append(d.actions, func() {
		d.store[key] = val
	})
}

func (d *DataStore) Commit() {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.actions = nil
}

func (d *DataStore) InTransaction() {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.inTransaction = true
}

func (d *DataStore) Rollback() {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.inTransaction = false
	for _, fun := range d.actions {
		fun()
	}
}
