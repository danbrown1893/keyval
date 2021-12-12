package cache

import (
	"errors"
	"sync"
	"time"
)

// ErrNotFound is returned when the key is not found
var ErrNotFound = errors.New("not found")

// Cache -
type Cache struct {
	cache    map[string]CacheItem
	deadline time.Time
	m        sync.Mutex
}

// CacheItem represents an item to be stored in the cache
type CacheItem struct {
	Value    interface{}
	Deadline time.Time
}

// NewCache initializes a new Cache
func NewCache() *Cache {
	return &Cache{
		cache: make(map[string]CacheItem),
	}
}

// EnsureKey will not allow you to overwrite an existing value, use Set directly if you want to overwrite
func (c *Cache) EnsureKey(key string, value interface{}, deadline time.Time) (interface{}, bool, error) {
	v, err := c.Get(key)
	if err == ErrNotFound {
		return value, true, c.Set(key, value, deadline)
	}
	if err != nil {
		return nil, false, err
	}
	return v, false, nil
}

// KeyExists will check if the Key exists in the cache
func (c *Cache) KeyExists(key string) bool {
	_, err := c.Get(key)
	return err == nil
}

// Get will return the value for a given key if it exists
func (c *Cache) Get(key string) (interface{}, error) {
	c.m.Lock()
	defer c.m.Unlock()

	v, found := c.cache[key]
	if !found {
		return nil, ErrNotFound
	}
	if !v.Deadline.IsZero() && time.Now().UTC().After(v.Deadline) {
		delete(c.cache, key)
		return nil, ErrNotFound
	}
	return v.Value, nil
}

// Set will set a key directly - It will overwrite an existing key value pair
func (c *Cache) Set(key string, value interface{}, deadline time.Time) error {
	c.m.Lock()
	defer c.m.Unlock()

	c.cache[key] = CacheItem{
		Value:    value,
		Deadline: deadline.UTC(),
	}

	// Will delete the key/value pair when it is expired
	if !deadline.IsZero() {
		go func() {
			time.Sleep(deadline.Sub(time.Now()) + time.Second)
			if c != nil {
				c.Get(key)
			}
		}()
	}

	return nil
}

// Delete will delete a key if it exists
func (c *Cache) Delete(key string) {
	delete(c.cache, key)
}
