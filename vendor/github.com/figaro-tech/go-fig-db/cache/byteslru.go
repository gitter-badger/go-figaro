// Copyright 2018 The Figaro Authors.
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

	"github.com/figaro-tech/go-fig-db/types"
)

// BytesLRU is an LRU cache with bytes values.
type BytesLRU struct {
	LRU
}

// NewBytesLRU creates a new LRU. If maxEntries == 0, the cache
// is disabled and does nothing.
func NewBytesLRU(maxEntries int) *BytesLRU {
	return &BytesLRU{
		LRU{
			MaxEntries: maxEntries,
			ll:         list.New(),
			cache:      make(map[string]*list.Element),
		},
	}
}

// Add adds a value to the cache.
func (c *BytesLRU) Add(key types.Key, value []byte) {
	c.LRU.Add(key, value)
}

// Get looks up a key's value from the cache.
func (c *BytesLRU) Get(key types.Key) (value []byte, ok bool) {
	v, ok := c.LRU.Get(key)
	return v.([]byte), ok
}
