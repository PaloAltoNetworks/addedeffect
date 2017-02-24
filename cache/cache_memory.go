package cache

import (
	"sync"
	"time"
)

type memoryCache struct {
	data       map[string]*cacheItem
	lock       *sync.Mutex
	expiration time.Duration
	expirer    ExpirationNotifier
}

// NewMemoryCache returns a new generic cache.
func NewMemoryCache() Cacher {

	return &memoryCache{
		data:       map[string]*cacheItem{},
		lock:       &sync.Mutex{},
		expiration: -1,
	}
}

func (c *memoryCache) SetDefaultExpiration(exp time.Duration) {

	c.expiration = exp
}

func (c *memoryCache) SetDefaultExpirationNotifier(expirer ExpirationNotifier) {

	c.expirer = expirer
}

func (c *memoryCache) Get(id string) interface{} {

	c.lock.Lock()
	item, ok := c.data[id]
	c.lock.Unlock()

	if !ok {
		return nil
	}

	return item.data
}

func (c *memoryCache) GetReset(id string) interface{} {

	c.lock.Lock()
	defer c.lock.Unlock()

	item, ok := c.data[id]
	if !ok {
		return nil
	}

	if c.expiration != -1 {
		if item.timer != nil {
			item.timer.Stop()
		}
		item.timer = time.AfterFunc(c.expiration, func() { c.delNotify(id, true) })
	}

	return item.data
}

func (c *memoryCache) Set(id string, item interface{}) {

	c.SetWithExpiration(id, item, c.expiration)
}

func (c *memoryCache) SetWithExpiration(id string, item interface{}, exp time.Duration) {

	c.SetWithExpirationAndNotifier(id, item, c.expiration, c.expirer)
}

func (c *memoryCache) SetWithExpirationAndNotifier(id string, item interface{}, exp time.Duration, expirer ExpirationNotifier) {

	var timer *time.Timer
	if exp != -1 {
		timer = time.AfterFunc(exp, func() { c.delNotify(id, true) })
	}

	ci := &cacheItem{
		identifier: id,
		data:       item,
		timestamp:  time.Now(),
		timer:      timer,
		expirer:    expirer,
	}

	c.lock.Lock()
	if item, ok := c.data[id]; ok && item.timer != nil {
		item.timer.Stop()
	}
	c.data[id] = ci
	c.lock.Unlock()
}

func (c *memoryCache) delNotify(id string, notify bool) {

	c.lock.Lock()
	item, ok := c.data[id]
	if ok && item.timer != nil {
		item.timer.Stop()
	}
	delete(c.data, id)
	c.lock.Unlock()

	if !ok {
		return
	}

	if !notify {
		return
	}

	if item.expirer != nil {
		item.expirer(c, id, item.data)
	}
}

func (c *memoryCache) Del(id string) {

	c.delNotify(id, false)
}

func (c *memoryCache) Exists(id string) bool {

	c.lock.Lock()
	_, ok := c.data[id]
	c.lock.Unlock()

	return ok
}

func (c *memoryCache) All() map[string]interface{} {

	out := map[string]interface{}{}

	c.lock.Lock()
	for k, i := range c.data {
		out[k] = i.data
	}
	c.lock.Unlock()

	return out
}
