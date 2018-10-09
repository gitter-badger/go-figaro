// Package cache implements caches.
package cache

import (
	"container/list"

	"github.com/figaro-tech/go-fig-db/types"
)

// BytesFIFO is an FIFO cache. It's just like a LRU, only it doesn't
// update on Get, only Add.
type BytesFIFO struct {
	FIFO
}

// NewBytesFIFO creates a new LRU. If maxEntries == 0, the cache
// is disabled and does nothing.
func NewBytesFIFO(maxEntries int) *BytesFIFO {
	return &BytesFIFO{
		FIFO{
			LRU{
				MaxEntries: maxEntries,
				ll:         list.New(),
				cache:      make(map[string]*list.Element),
			},
		},
	}
}

// Add adds a value to the cache.
func (c *BytesFIFO) Add(key types.Key, value []byte) {
	c.FIFO.Add(key, value)
}

// Get looks up a key's value from the cache.
func (c *BytesFIFO) Get(key types.Key) (value []byte, ok bool) {
	v, ok := c.FIFO.Get(key)
	return v.([]byte), ok
}
