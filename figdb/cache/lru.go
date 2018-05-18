// Copyright 2018 The Figaro Authors.
// <License goes here>
// Based on Google's Groupcache LRU https://github.com/golang/groupcache/blob/master/lru/lru.go

/*
Copyright 2013 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package cache implements caches.
package cache

import (
	"container/list"
	"sync"

	"github.com/figaro-tech/go-figaro/figdb/types"
)

// LRU is an LRU cache.
type LRU struct {
	lock       sync.RWMutex
	MaxEntries int
	ll         *list.List
	cache      map[string]*list.Element
}

type entry struct {
	key   types.Key
	value string
}

// NewLRU creates a new LRU. If maxEntries == 0, the cache
// is disabled and does nothing.
func NewLRU(maxEntries int) *LRU {
	return &LRU{
		MaxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
	}
}

// Add adds a value to the cache.
func (c *LRU) Add(key types.Key, value []byte) {
	if c.MaxEntries == 0 {
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.cache == nil {
		c.cache = make(map[string]*list.Element)
		c.ll = list.New()
	}
	if ee, ok := c.cache[key.String()]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = string(value)
		return
	}
	ele := c.ll.PushFront(&entry{key, string(value)})
	c.cache[key.String()] = ele
	if c.ll.Len() > c.MaxEntries {
		c.RemoveOldest()
	}
}

// Get looks up a key's value from the cache.
func (c *LRU) Get(key types.Key) (value []byte, ok bool) {
	if c.MaxEntries == 0 {
		return
	}
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key.String()]; hit {
		return []byte(ele.Value.(*entry).value), true
	}
	return
}

// Remove removes the provided key from the cache.
func (c *LRU) Remove(key types.Key) {
	if c.MaxEntries == 0 {
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key.String()]; hit {
		c.removeElement(ele)
	}
}

// RemoveOldest removes the oldest item from the cache.
func (c *LRU) RemoveOldest() {
	if c.MaxEntries == 0 {
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *LRU) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key.String())
}

// Len returns the number of items in the cache.
func (c *LRU) Len() int {
	if c.MaxEntries == 0 {
		return 0
	}
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

// Clear purges all stored items from the cache.
func (c *LRU) Clear() {
	if c.MaxEntries == 0 {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	c.ll = nil
	c.cache = nil
}
