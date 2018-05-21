// Package signature provides cryptographic signing functions
package signature

import (
	"crypto"
	"io"

	"github.com/figaro-tech/go-figaro/figcrypto/signature/common"
	"github.com/figaro-tech/go-figaro/figcrypto/signature/ecdsa"
)

// Constants
const (
	PrivateKeySize = ecdsa.PrivateKeySize
	SignatureSize  = ecdsa.SignatureSize
	PublicKeySize  = ecdsa.PublicKeySize
	AddressSize    = ecdsa.AddressSize
)

// ToHumanAddress converts a binary address to a Base58 encoded "human readable" address
func ToHumanAddress(binaddr []byte) string {
	return common.ToHumanAddress(binaddr)
}

// ToBinaryAddress converts a Base58 encoded "human readable" address to a binary address
func ToBinaryAddress(humaddr string) []byte {
	return common.ToBinaryAddress(humaddr)
}

// GenerateKey generates an address/private key pair that can be used to sign and verify messages,
// If rand is nil, crypto/rand.Reader will be used.
func GenerateKey(rando io.Reader) (privkey, publickey, address []byte, err error) {
	return ecdsa.GenerateKey(rando)
}

// GenerateKeyFromSeed creates an public/private key pair, along with an address,
// that can be used to verify messages, from a seed string or strings.
func GenerateKeyFromSeed(seed ...string) (privkey, publickey, address []byte, err error) {
	return ecdsa.GenerateKeyFromSeed(seed...)
}

// RecoverFromPrivateKey gets a public key and address which can verify transactions from a private key.
func RecoverFromPrivateKey(privkey []byte) (pubkey, address []byte, err error) {
	return ecdsa.RecoverFromPrivateKey(privkey)
}

// SignerFromPrivateKey returns a `crypto.Signer` from private key bytes.
func SignerFromPrivateKey(privkey []byte) (crypto.Signer, error) {
	return ecdsa.SignerFromPrivateKey(privkey)
}

// Sign signs the message using the private key.
func Sign(privkey, message []byte) (signature []byte, err error) {
	return ecdsa.Sign(privkey, message)
}

// Verify verifies that a message was signed by the owner of the address.
func Verify(pubkey, signature, message []byte) bool {
	return ecdsa.Verify(pubkey, signature, message)
}

// VerifyAddress verifies that the public key is valid for the address.
func VerifyAddress(publickey, address []byte) bool {
	return ecdsa.VerifyAddress(publickey, address)
}

// VerifyWithAddress verifies that the public key is valid for the address.
func VerifyWithAddress(address, signature, message []byte) bool {
	return ecdsa.VerifyWithAddress(address, signature, message)
}
