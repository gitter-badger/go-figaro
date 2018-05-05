// Package figcrypto provides a subset of convenient crypto functions
package figcrypto

import "github.com/figaro-tech/figaro/pkg/figcrypto/internal/sha256"

// Hasher provides a convenient and fast cryptographic Hash function
type Hasher struct {
	sha256.Hasher
}
