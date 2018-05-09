// Package figcrypto provides cryptographic functions
package figcrypto

import (
	"hash"
	"sync"

	"golang.org/x/crypto/blake2b"
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
	h, err := blake2b.New256(nil)
	if err != nil {
		panic(err)
	}
	return &Hasher{
		hash: h,
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

// Hash returns a hash of 0 or more []byte
func Hash(b ...[]byte) []byte {
	if len(b) == 1 {
		h := blake2b.Sum256(b[0])
		return h[:]
	}
	h := HasherPool.Get().(*Hasher)
	defer HasherPool.Put(h)
	return h.Hash(b...)
}
