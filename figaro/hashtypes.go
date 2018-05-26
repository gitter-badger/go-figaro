// Package figaro is the main package for go-figaro
package figaro

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/figaro-tech/go-figaro/figcrypto/signature/fastsig"
)

// HashTypes are []byte hash types of length 32 with special functionality.

var (
	// ZeroAddress is a special system address of all zeros.
	ZeroAddress Address = bytes.Repeat([]byte{0x00}, fastsig.AddressSize)
	// ErrInvalidAddressData is a self-explantory error.
	ErrInvalidAddressData = errors.New("figaro address: invalid data")
	// ErrInvalidRootData is a self-explantory error.
	ErrInvalidRootData = errors.New("figaro root: invalid data")
	// ErrInvalidBlockHashData is a self-explantory error.
	ErrInvalidBlockHashData = errors.New("figaro block: invalid BlockHash data")
	// ErrInvalidTxHashData is a self-explantory error.
	ErrInvalidTxHashData = errors.New("figaro tx: invalid TxHash data")
	// ZeroAddress is a special zeroeth address used to denote
	// transactions to the network.
)

// NewAddressFromHuman is a convenience helper to create an address from a Base58 encoded string.
func NewAddressFromHuman(humaddr string) (address *Address, err error) {
	address = &Address{}
	err = address.SetHuman(humaddr)
	return
}

// AddressSize is the size, in bytes, of an Address.
const AddressSize = fastsig.AddressSize

// An Address a unique address used for accounts.
type Address []byte

// Valid returns whether a []byte is a valid Address.
func (addr Address) Valid() bool {
	return len(addr) == AddressSize
}

// IsZeroAddress returns whether an Address is the ZeroAddress.
func (addr Address) IsZeroAddress() bool {
	return len(addr) == 0 || bytes.Equal(addr, ZeroAddress)
}

// String converts to a string.
func (addr Address) String() string { return fmt.Sprintf("%#x", addr) }

// Human converts to a Base58 encoded string.
func (addr Address) Human() string { return fastsig.ToHumanAddress(addr) }

// Hex converts to a hex encoded string.
func (addr Address) Hex() string { return hex.EncodeToString(addr) }

// SetHuman sets an address from a Base58 encoded string.
func (addr Address) SetHuman(humaddr string) error {
	// TODO: validate the address
	b := fastsig.ToBinaryAddress(humaddr)
	if len(b) != fastsig.AddressSize {
		return ErrInvalidAddressData
	}
	copy(addr, b)
	return nil
}

// SetHex sets an address from a hex encoded string.
func (addr Address) SetHex(h string) error {
	b, err := hex.DecodeString(h)
	if err != nil {
		return err
	}
	copy(addr, b)
	return nil
}

// RootSize is the size, in bytes, of a Root.
const RootSize = 32

// An Root a unique root used for accounts.
type Root []byte

// Valid returns whether a []byte is a valid Root.
func (root Root) Valid() bool {
	return len(root) == RootSize
}

// String converts to a string.
func (root Root) String() string { return fmt.Sprintf("%#x", root) }

// Hex converts to a hex encoded string.
func (root Root) Hex() string { return hex.EncodeToString(root) }

// SetHex sets an root from a hex encoded string.
func (root Root) SetHex(h string) error {
	b, err := hex.DecodeString(h)
	if err != nil {
		return err
	}
	copy(root, b)
	return nil
}

// BlockHashSize is the size, in bytes, of a BlockHash.
const BlockHashSize = 32

// An BlockHash a unique root used for accounts.
type BlockHash []byte

// Valid returns whether a []byte is a valid BlockHash.
func (bh BlockHash) Valid() bool {
	return len(bh) == BlockHashSize
}

// String converts to a string.
func (bh BlockHash) String() string { return fmt.Sprintf("%#x", bh) }

// Hex converts to a hex encoded string.
func (bh BlockHash) Hex() string { return hex.EncodeToString(bh) }

// SetHex sets an bh from a hex encoded string.
func (bh BlockHash) SetHex(h string) error {
	b, err := hex.DecodeString(h)
	if err != nil {
		return err
	}
	copy(bh, b)
	return nil
}

// TxHashSize is the size, in bytes, of a TxHash.
const TxHashSize = 32

// TxHash is a "fingerprint" of a Tx suitable for use as an ID.
type TxHash []byte

// Valid returns whether a []byte is a valid BlockHash.
func (txhash TxHash) Valid() bool {
	return len(txhash) == TxHashSize
}

// String converts to a string.
func (txhash TxHash) String() string { return fmt.Sprintf("%#x", txhash) }

// Hex converts to a hex encoded string.
func (txhash TxHash) Hex() string { return hex.EncodeToString(txhash) }

// SetHex sets an txhash from a hex encoded string.
func (txhash TxHash) SetHex(h string) error {
	b, err := hex.DecodeString(h)
	if err != nil {
		return err
	}
	copy(txhash, b)
	return nil
}
