// Package figaro is the main package for go-figaro
package figaro

// A Hasher provides cryptographic hashing
type Hasher interface {
	// Hash accepts 0 or more []byte and hashes them, returning a []byte
	Hash(b ...[]byte) []byte
}
