package cache

import "time"

// A Cacher is the interface caching struct have to implement
type Cacher interface {
	SetDefaultExpiration(exp time.Duration)
	Set(id string, item interface{})
	SetWithExpiration(id string, item interface{}, exp time.Duration)
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
}

// ExpirationNotifier is an interface that cacheable structs can implement to be
// notified in case of expiration
type ExpirationNotifier interface {
	Expired(c Cacher, id string)
}
