// Package sha256 provides cryptographic hashing
package sha256

import "crypto/sha256"

// Hasher impelments figaro.Hasher
type Hasher struct{}

// Hash returns a hash of 0 or more []byte
func (h *Hasher) Hash(b ...[]byte) []byte {
	hash := sha256.New()
	for _, item := range b {
		hash.Write(item)
	}
	return hash.Sum(nil)
}
