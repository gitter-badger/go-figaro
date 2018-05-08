// Package figcrypto provides cryptographic hashing
package figcrypto

import (
	"crypto/sha256"
	"hash"
	"sync"
)

// HasherPool is a thread-safe pool of Hashers
var HasherPool = sync.Pool{
	New: func() interface{} {
		return NewHasher()
	},
}

// A Hasher hashes efficiently
type Hasher struct {
	hash hash.Hash
}

// NewHasher returns a Hasher ready to use
func NewHasher() *Hasher {
	return &Hasher{
		hash: sha256.New(),
	}
}

// Hash returns a hash of 0 or more []byte
func (h *Hasher) Hash(b ...[]byte) []byte {
	h.hash.Reset()
	for _, item := range b {
		h.hash.Write(item)
	}
	return h.hash.Sum(nil)
}
