// Package hash provides cryptographic functions
package hash

import (
	"crypto"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ripemd160"
)

// Hash is the crypto hash used by figcrypto
const Hash = crypto.BLAKE2b_256

// Hash256 returns a hash of 0 or more []byte
func Hash256(b ...[]byte) []byte {
	if len(b) == 1 {
		h := blake2b.Sum256(b[0])
		return h[:]
	}
	h, err := blake2b.New256(nil)
	if err != nil {
		panic(err)
	}
	for _, item := range b {
		h.Write(item)
	}
	return h.Sum(nil)
}

// Hash512 returns a hash of 0 or more []byte
func Hash512(b ...[]byte) []byte {
	if len(b) == 1 {
		h := blake2b.Sum512(b[0])
		return h[:]
	}
	h, err := blake2b.New512(nil)
	if err != nil {
		panic(err)
	}
	for _, item := range b {
		h.Write(item)
	}
	return h.Sum(nil)
}

// Hash160 returns a hash of 0 or more []byte
func Hash160(b ...[]byte) []byte {
	h := ripemd160.New()
	for _, item := range b {
		h.Write(item)
	}
	return h.Sum(nil)
}
