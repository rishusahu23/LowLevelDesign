package caches

type Cache interface {
	Set(key string, val interface{})
	Get(key string) (val interface{}, bool2 bool)
	Remove(key string)
}

type EvictionStrategy interface {
	Evict(c Cache)
	OnAdd(c Cache, key string)
	OnAccess(c Cache, key string)
	OnRemove(c Cache, key string)
}
