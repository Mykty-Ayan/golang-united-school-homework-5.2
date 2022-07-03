package cache

import (
	"time"
)

type cacheValue struct {
	value    string
	deadline time.Time
}

type Cache struct {
	cacheMap map[string]cacheValue
}

func NewCache() Cache {
	return Cache{
		cacheMap: map[string]cacheValue{},
	}
}

func (c Cache) Get(key string) (string, bool) {
	c.deleteExpiredKeys()
	if val, ok := c.cacheMap[key]; ok {
		return val.value, ok
	}
	return "", false
}

func (c *Cache) Put(key, value string) {
	cv := cacheValue{value: value, deadline: time.Time{}}
	c.cacheMap[key] = cv
}

func (c *Cache) Keys() []string {
	c.deleteExpiredKeys()
	keys := make([]string, 0)
	for key, _ := range c.cacheMap {
		keys = append(keys, key)
	}
	return keys
}

func (c *Cache) PutTill(key, value string, deadline time.Time) {
	if deadline.IsZero() {
		deadline = time.Now().Add(-time.Second)
	}
	cv := cacheValue{value, deadline}
	c.cacheMap[key] = cv
}

func (c *Cache) deleteExpiredKeys() {
	expiredKeys := make([]string, 0)
	for key, value := range c.cacheMap {
		if value.deadline.IsZero() {
			continue
		}
		if value.deadline.After(time.Now()) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	for _, key := range expiredKeys {
		delete(c.cacheMap, key)
	}
}
