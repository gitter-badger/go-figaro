// Package figaro is the main package for go-figaro
package figaro

// Store provides a key/value store
type Store interface {
	Get(key []byte) []byte
	Set(key []byte, value []byte)
	Delete(key []byte)
}
