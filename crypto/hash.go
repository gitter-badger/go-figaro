package crypto

import (
	"crypto/sha256"
	"hash"

	"golang.org/x/crypto/sha3"
)

// Cypher is a convenience for hashing
type Cypher struct{}

// Hash hashes a []byte array according the default algorithm
func (c *Cypher) Hash(data []byte) []byte {
	return Hash(data)
}

// NewHash returns a new hash.Hash according to the default algorithm
func (c *Cypher) NewHash() hash.Hash {
	return NewHash()
}

// Hash hashes a []byte array according the default algorithm
func Hash(data []byte) []byte {
	return Sha256(data)
}

// NewHash returns a new hash.Hash according to the default algorithm
func NewHash() hash.Hash {
	return NewSha256()
}

// Sha256 hashes a []byte
func Sha256(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

// NewSha256 returns a new hash.Hash according for Sha256
func NewSha256() hash.Hash {
	return sha256.New()
}

// Sha3 hashes a []byte
func Sha3(data []byte) []byte {
	h := sha3.Sum256(data)
	return h[:]
}

// NewSha3 returns a new hash.Hash according for Sha3
func NewSha3() hash.Hash {
	return sha3.New256()
}
