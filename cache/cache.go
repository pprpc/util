// Package cache manager
package cache

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Cache cache
type Cache struct {
	count int32
	max   int32
	cache *sync.Map
}

// RangeCall  range call back
type RangeCall func(k, v interface{}) bool

// NewCache set amx
func NewCache(max int32) *Cache {
	cc := new(Cache)
	cc.max = max
	cc.count = 0
	cc.cache = new(sync.Map)
	return cc
}

// Del delete
func (cc *Cache) Del(key string) {
	_, err := cc.Get(key)
	if err != nil {
		return
	}
	atomic.AddInt32(&cc.count, -1)
	cc.cache.Delete(key)
}

// Delete delete
func (cc *Cache) Delete(key string) {
	cc.Del(key)
}

// Set set value
func (cc *Cache) Set(key string, v interface{}) (c int32, err error) {
	c = atomic.LoadInt32(&cc.count)
	_, err = cc.Get(key)
	if err != nil {
		err = nil
		if c >= cc.max && cc.max != 0 {
			err = fmt.Errorf("cc, max connect over flow(%d)", cc.max)
			return
		}
		atomic.AddInt32(&cc.count, 1)
		cc.cache.Store(key, v)
		c = c + 1
	} else {
		cc.cache.Store(key, v)
	}
	return
}

// AddORUpdate set value
func (cc *Cache) AddORUpdate(key string, v interface{}) (c int32, err error) {
	c, err = cc.Set(key, v)
	return
}

// Get get value
func (cc *Cache) Get(key string) (v interface{}, err error) {
	var ok bool
	v, ok = cc.cache.Load(key)
	if !ok {
		err = fmt.Errorf("cc, not find: %s", key)
	}
	return
}

// Range range
func (cc *Cache) Range(fn RangeCall) {
	cc.cache.Range(fn)
}

// Len get record
func (cc *Cache) Len() int32 {
	return atomic.LoadInt32(&cc.count)
}
