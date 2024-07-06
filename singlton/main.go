package main

import "sync"

type Singleton struct {
	Data string
}

var once sync.Once
var ins *Singleton

func GetInstance(data string) *Singleton {
	once.Do(func() {
		ins = &Singleton{
			Data: data,
		}
	})
	return ins
}

func main() {
	ins1 := GetInstance("rishu")
	println(ins1.Data)

	ins2 := GetInstance("sahu")
	println(ins2.Data)

	ins3 := GetInstance("rishu")
	println(ins3.Data)
}
