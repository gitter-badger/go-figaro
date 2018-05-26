package common

import (
	"encoding/hex"
	"errors"
	"os"

	"github.com/btcsuite/btcutil/base58"
	"github.com/figaro-tech/go-figaro/figcrypto/hasher"
)

// AddressSize is the size of an address, in bytes.
const AddressSize = 25

var (
	// ErrInvalidKey is a self-explantory error
	ErrInvalidKey = errors.New("figcrypto signature: invalid public or private key")
	// ErrInvalidSignature is a self-explantory error
	ErrInvalidSignature = errors.New("figcrypto signature: invalid signature")
	// Using different versions for different nets means addresses won't
	// cross-validate. Override this with 1-byte hex encoded env var
	// `ADDRESS_VERSION_CODE`
	version = []byte{0x00}
)

func init() {
	if v, ok := os.LookupEnv("ADDRESS_VERSION_CODE"); ok {
		b, _ := hex.DecodeString(v)
		if len(b) == 1 {
			version = b
		}
	}
}

// ToHumanAddress converts a binary address to a Base58 encoded "human readable" address
func ToHumanAddress(binaddr []byte) string {
	return base58.Encode(binaddr)
}

// ToBinaryAddress converts a Base58 encoded "human readable" address to a binary address
func ToBinaryAddress(humaddr string) []byte {
	return base58.Decode(humaddr)
}

// AddressFromPublicKey creates an address from a public key, following the Bitcoin address protocol.
func AddressFromPublicKey(pubkey []byte) (address []byte) {
	h := hasher.Hash256(pubkey)
	h = hasher.Hash160(h)
	address = append(version, h...)
	h2 := hasher.Hash256(address)
	h2 = hasher.Hash256(h2)
	checksum := h2[:4]
	address = append(address, checksum...)
	return
}
