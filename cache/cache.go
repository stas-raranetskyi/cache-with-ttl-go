package cache

import (
	"runtime"
	"time"
)

type Item struct {
	value    any
	expireAt int64
}

type Cache struct {
	storage   map[string]Item
	stopClean chan bool
}

func (c Cache) Set(key string, value any, ttl time.Duration) {
	item := Item{value, time.Now().Add(ttl).UnixMilli()}
	c.storage[key] = item
}

func (c Cache) Get(key string) any {
	return c.storage[key].value
}

func (c Cache) Delete(key string) {
	delete(c.storage, key)
}

func (cache Cache) purge() {
	now := time.Now().UnixMilli()
	for key, data := range cache.storage {
		if data.expireAt < now {
			cache.Delete(key)
		}
	}
}

func clean(cache Cache) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			cache.purge()
		case <-cache.stopClean:
			ticker.Stop()
		}
	}
}

func stopCleaning(cache *Cache) {
	cache.stopClean <- true
}

func New() Cache {
	storage := make(map[string]Item)
	stop := make(chan bool)
	cache := Cache{storage, stop}
	go clean(cache)
	runtime.SetFinalizer(&cache, stopCleaning)
	return cache
}
