package main

type EvictionPolicy interface {
	KeyAccessed(key string)
	Evict() string
}


