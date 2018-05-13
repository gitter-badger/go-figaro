// Package signature provides cryptographic functions
package signature

import (
	"errors"
	"io"

	"golang.org/x/crypto/ed25519"
)

const (
	// PublicKeySize is the size, in bytes, of public keys as used in this package.
	PublicKeySize = ed25519.PublicKeySize
	// PrivateKeySize is the size, in bytes, of private keys as used in this package.
	PrivateKeySize = ed25519.PrivateKeySize
	// SignatureSize is the size, in bytes, of signatures generated and verified by this package.
	SignatureSize = ed25519.SignatureSize
)

// Errors used in package
var (
	ErrInvalidKey = errors.New("figcrypto signature: invalid public or private key")
)

// Signer is a crypto.Signer interface suitable for HD wallets
type Signer ed25519.PrivateKey

// GenerateKey generates a public/private key pair that can be used to sign and verify messages.
// If rand is nil, crypto/rand.Reader will be used.
func GenerateKey(rand io.Reader) (publicKey, privateKey []byte, err error) {
	return ed25519.GenerateKey(rand)
}

// Sign signs the message using the 64-byte private key
func Sign(privateKey, message []byte) ([]byte, error) {
	if len(privateKey) != PrivateKeySize {
		return nil, ErrInvalidKey
	}
	return ed25519.Sign(ed25519.PrivateKey(privateKey), message), nil
}

// Verify verifies that a message was signed by the owner of the public key
func Verify(publicKey, message, sig []byte) bool {
	if len(publicKey) != PublicKeySize {
		return false
	}
	return ed25519.Verify(ed25519.PublicKey(publicKey), message, sig)
}
