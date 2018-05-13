// Package hash provides cryptographic functions
package hash

import (
	"crypto"
	"hash"
	"sync"

	"golang.org/x/crypto/blake2b"
)

// Hash is the crypto hash used by figcrypto
const Hash = crypto.BLAKE2b_256

// HasherPool is a thread-safe pool of Hashers
var HasherPool = sync.Pool{
	New: func() interface{} {
		return NewHasher()
	},
}

// A Hasher hashes efficiently
type Hasher struct {
	Hasher hash.Hash
}

// NewHasher returns a Hasher ready to use
func NewHasher() *Hasher {
	h, err := blake2b.New256(nil)
	if err != nil {
		panic(err)
	}
	return &Hasher{
		Hasher: h,
	}
}

// Hash256 returns a hash of 0 or more []byte
func (h *Hasher) Hash256(b ...[]byte) []byte {
	h.Hasher.Reset()
	for _, item := range b {
		h.Hasher.Write(item)
	}
	return h.Hasher.Sum(nil)
}

// Hash256 returns a hash of 0 or more []byte
func Hash256(b ...[]byte) []byte {
	if len(b) == 1 {
		h := blake2b.Sum256(b[0])
		return h[:]
	}
	h := HasherPool.Get().(*Hasher)
	defer HasherPool.Put(h)
	return h.Hash256(b...)
}
