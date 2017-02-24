package cache

import "time"

// ExpirationNotifier is a function which will be called every time a cache
// expires an item
type ExpirationNotifier func(c Cacher, id string, item interface{})

// A Cacher is the interface caching struct have to implement
type Cacher interface {
	SetDefaultExpiration(exp time.Duration)
	SetDefaultExpirationNotifier(expNotifier ExpirationNotifier)
	Set(id string, item interface{})
	SetWithExpiration(id string, item interface{}, exp time.Duration)
	SetWithExpirationAndNotifier(id string, item interface{}, exp time.Duration, expNotifier ExpirationNotifier)
	Get(id string) interface{}
	GetReset(id string) interface{}
	Del(id string)
	Exists(id string) bool
	All() map[string]interface{}
}

type cacheItem struct {
	timestamp  time.Time
	identifier string
	data       interface{}
	timer      *time.Timer
	expirer    ExpirationNotifier
}
