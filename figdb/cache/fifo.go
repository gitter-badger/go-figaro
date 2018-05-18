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

	"github.com/figaro-tech/go-figaro/figdb/types"
)

// FIFO is an FIFO cache. It's just like a LRU, only it doesn't
// update on Get, only Add.
type FIFO struct {
	LRU
}

// NewFIFO creates a new LRU. If maxEntries == 0, the cache
// is disabled and does nothing.
func NewFIFO(maxEntries int) *FIFO {
	return &FIFO{
		LRU{
			MaxEntries: maxEntries,
			ll:         list.New(),
			cache:      make(map[string]*list.Element),
		},
	}
}

// Get looks up a key's value from the cache.
func (c *FIFO) Get(key types.Key) (value []byte, ok bool) {
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
