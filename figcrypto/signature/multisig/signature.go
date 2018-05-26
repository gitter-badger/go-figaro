// Package multisig provides cryptographic functions
package multisig

// NOTE: multiig is a very much a work in progress

import (
	"bytes"
	"crypto"
	"io"

	"github.com/figaro-tech/go-figaro/figcrypto/hasher"
	"github.com/figaro-tech/go-figaro/figcrypto/signature/common"

	"golang.org/x/crypto/ed25519"
)

// Constants
const (
	AddressSize    = common.AddressSize
	PublicKeySize  = ed25519.PublicKeySize
	SignatureSize  = ed25519.SignatureSize
	PrivateKeySize = ed25519.PrivateKeySize
)

// ToHumanAddress converts a binary address to a Base58 encoded "human readable" address
func ToHumanAddress(binaddr []byte) string {
	return common.ToHumanAddress(binaddr)
}

// ToBinaryAddress converts a Base58 encoded "human readable" address to a binary address
func ToBinaryAddress(humaddr string) []byte {
	return common.ToBinaryAddress(humaddr)
}

// GenerateKey generates an public/private key pair, along with an address,
// that can be used to sign and verify messages.
// If rand is nil, crypto/rand.Reader will be used.
func GenerateKey(rando io.Reader) (privkey, pubkey, address []byte, err error) {
	pubkey, privkey, err = ed25519.GenerateKey(rando)
	address = common.AddressFromPublicKey(pubkey)
	return
}

// GenerateKeyFromSeed creates an public/private key pair, along with an address,
// that can be used to verify messages, from a seed string or strings.
func GenerateKeyFromSeed(seed ...string) (privkey, publickey, address []byte, err error) {
	var h []byte
	for _, s := range seed {
		h = hasher.Hash512(h, []byte(s))
	}
	rando := bytes.NewReader(h)
	return GenerateKey(rando)
}

// RecoverFromPrivateKey gets an address which can verify transactions from a private key.
func RecoverFromPrivateKey(privkey []byte) (pubkey, address []byte, err error) {
	if len(privkey) != PrivateKeySize {
		err = common.ErrInvalidKey
		return
	}
	pubkey = ed25519.PrivateKey(privkey).Public().([]byte)
	address = common.AddressFromPublicKey(pubkey)
	return
}

// SignerFromPrivateKey returns a `crypto.Signer` from private key bytes.
func SignerFromPrivateKey(privkey []byte) (crypto.Signer, error) {
	if len(privkey) != PrivateKeySize {
		return nil, common.ErrInvalidKey
	}
	return ed25519.PrivateKey(privkey), nil
}

// Sign signs the message using the private key
func Sign(privkey, message []byte) (signature []byte, err error) {
	if len(privkey) != PrivateKeySize {
		return nil, common.ErrInvalidKey
	}
	return ed25519.Sign(privkey, message), nil
}

// Verify verifies that a message was signed by the owner of the address
func Verify(publickey, signature, message []byte) bool {
	if len(publickey) != PublicKeySize || len(signature) != SignatureSize {
		return false
	}
	return ed25519.Verify(publickey, message, signature)
}

// VerifyAddress verifies that the public key is valid for the address.
func VerifyAddress(publickey, address []byte) bool {
	if len(publickey) != PublicKeySize || len(address) != AddressSize {
		return false
	}
	derived := common.AddressFromPublicKey(publickey)
	return bytes.Equal(address, derived)
}
