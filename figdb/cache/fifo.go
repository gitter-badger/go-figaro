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
func (c *FIFO) Get(key types.Key) (value interface{}, ok bool) {
	if c.MaxEntries == 0 {
		return
	}

	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key.String()]; hit {
		return ele.Value.(*entry).value, true
	}
	return
}
