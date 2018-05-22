// Package figaro is the main package for go-figaro
package figaro

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/figaro-tech/go-figaro/figcrypto/signature/fastsig"
)

// ErrInvalidAddressData is a self-explantory error.
var ErrInvalidAddressData = errors.New("figaro address: invalid data")

// An Address a unique address used for accounts.
type Address [fastsig.AddressSize]byte

// NewAddressFromHuman is a convenience helper to create an address from a Base58 encoded string.
func NewAddressFromHuman(humaddr string) (address *Address, err error) {
	address = &Address{}
	err = address.SetHuman(humaddr)
	return
}

// String converts to a string.
func (addr Address) String() string { return fmt.Sprintf("%#x", addr) }

// Bytes converts to bytes.
func (addr Address) Bytes() []byte { return addr[:] }

// Human converts to a Base58 encoded string.
func (addr Address) Human() string { return fastsig.ToHumanAddress(addr.Bytes()) }

// Hex converts to a hex encoded string.
func (addr Address) Hex() string { return hex.EncodeToString(addr.Bytes()) }

// SetBytes sets an address from bytes.
func (addr *Address) SetBytes(b []byte) error {
	if len(b) != fastsig.AddressSize {
		return ErrInvalidAddressData
	}
	copy(addr.Bytes(), b)
	return nil
}

// SetHuman sets an address from a Base58 encoded string.
func (addr *Address) SetHuman(humaddr string) error {
	// TODO: validate the address
	b := fastsig.ToBinaryAddress(humaddr)
	if len(b) != fastsig.AddressSize {
		return ErrInvalidAddressData
	}
	copy(addr.Bytes(), b)
	return nil
}

// SetHex sets an address from a hex encoded string.
func (addr *Address) SetHex(h string) error {
	b, err := hex.DecodeString(h)
	if err != nil {
		return err
	}
	copy(addr.Bytes(), b)
	return nil
}
