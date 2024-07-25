package singleton

import "sync"

type Singleton struct {
}

var once sync.Once
var instance *Singleton

func getInstance() *Singleton {
	once.Do(func() {
		instance = &Singleton{}
	})
	return instance
}
