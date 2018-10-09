// Package fastsig provides cryptographic functions
package fastsig

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"io"
	"math/big"

	"github.com/figaro-tech/go-fig-crypto/hasher"
	"github.com/figaro-tech/go-fig-crypto/signature/fastsig/secp256k1"

	"github.com/figaro-tech/go-fig-crypto/signature/common"
)

// Constants
const (
	AddressSize               = common.AddressSize
	SignatureSize             = 65
	PublicKeySize             = 33
	DecompressedPublicKeySize = 65
	PrivateKeySize            = 32
)

// ToHumanAddress converts a binary address to a Base58 encoded "human readable" address
func ToHumanAddress(binaddr []byte) string {
	return common.ToHumanAddress(binaddr)
}

// ToBinaryAddress converts a Base58 encoded "human readable" address to a binary address
func ToBinaryAddress(humaddr string) []byte {
	return common.ToBinaryAddress(humaddr)
}

// DecompressPubkey parses a public key in the 33-byte compressed format.
func DecompressPubkey(pubkey []byte) ([]byte, error) {
	if len(pubkey) != PublicKeySize {
		return nil, common.ErrInvalidKey
	}
	x, y := secp256k1.DecompressPubkey(pubkey)
	if x == nil {
		return nil, common.ErrInvalidKey
	}
	return elliptic.Marshal(secp256k1.S256(), x, y), nil
}

// CompressPubkey encodes a public key to the 33-byte compressed format.
func CompressPubkey(pubkey []byte) ([]byte, error) {
	if len(pubkey) != DecompressedPublicKeySize {
		return nil, common.ErrInvalidKey
	}
	x, y := elliptic.Unmarshal(secp256k1.S256(), pubkey)
	return secp256k1.CompressPubkey(x, y), nil
}

// GenerateKey generates an public/private key pair, along with an address,
// that can be used to sign and verify messages.
// If rand is nil, crypto/rand.Reader will be used.
func GenerateKey(rando io.Reader) (privkey, pubkey, address []byte, err error) {
	if rando == nil {
		rando = rand.Reader
	}
	var key *ecdsa.PrivateKey
	key, err = ecdsa.GenerateKey(secp256k1.S256(), rando)
	if err != nil {
		return
	}
	privkey = key.D.Bytes()
	pubkey = secp256k1.CompressPubkey(key.X, key.Y)
	address = common.AddressFromPublicKey(pubkey)
	return
}

// GenerateKeyFromSeed creates an public/private key pair, along with an address,
// that can be used to verify messages, from a seed string or strings.
func GenerateKeyFromSeed(seed ...string) (privkey, pubkey, address []byte, err error) {
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
	c := secp256k1.S256()
	x, y := c.ScalarBaseMult(privkey)
	pubkey = secp256k1.CompressPubkey(x, y)
	address = common.AddressFromPublicKey(pubkey)
	return
}

// RecoverFromSignature gets an address which can verify transactions from a signature and message.
func RecoverFromSignature(signature, message []byte) (pubkey, address []byte, err error) {
	if len(signature) != SignatureSize {
		return nil, nil, common.ErrInvalidSignature
	}
	h := hasher.Hash256(message)
	rec, err := secp256k1.RecoverPubkey(h, signature)
	if err != nil {
		return nil, nil, err
	}
	x, y := elliptic.Unmarshal(secp256k1.S256(), rec)
	pubkey = secp256k1.CompressPubkey(x, y)
	address = common.AddressFromPublicKey(pubkey)
	return
}

// SignerFromPrivateKey returns a `crypto.Signer` from private key bytes.
func SignerFromPrivateKey(privkey []byte) (crypto.Signer, error) {
	if len(privkey) != PrivateKeySize {
		return nil, common.ErrInvalidKey
	}
	c := secp256k1.S256()
	d := new(big.Int).SetBytes(privkey)
	x, y := c.ScalarBaseMult(d.Bytes())
	return &ecdsa.PrivateKey{
		D:         d,
		PublicKey: ecdsa.PublicKey{Curve: c, X: x, Y: y},
	}, nil
}

// Sign signs the message using the private key
func Sign(privkey, message []byte) (signature []byte, err error) {
	if len(privkey) != PrivateKeySize {
		return nil, common.ErrInvalidKey
	}
	h := hasher.Hash256(message)
	return secp256k1.Sign(h, privkey)
}

// Verify verifies that a message was signed by the owner of the compact encoded public key.
func Verify(pubkey, signature, message []byte) bool {
	if len(pubkey) != PublicKeySize || len(signature) != SignatureSize {
		return false
	}
	h := hasher.Hash256(message)
	return secp256k1.VerifySignature(pubkey, h, signature[:SignatureSize-1])
}

// VerifyAddress verifies that the public key is valid for the address.
func VerifyAddress(publickey, address []byte) bool {
	if len(publickey) != PublicKeySize || len(address) != AddressSize {
		return false
	}
	derived := common.AddressFromPublicKey(publickey)
	return bytes.Equal(address, derived)
}

// VerifyWithAddress verifies that a message was signed by the owner of the address.
func VerifyWithAddress(address, signature, message []byte) bool {
	if len(address) != AddressSize {
		return false
	}
	_, rec, err := RecoverFromSignature(signature, message)
	if err != nil {
		return false
	}
	return bytes.Equal(address, rec)
}
