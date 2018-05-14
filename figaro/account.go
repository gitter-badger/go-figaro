// Package figaro is the main package for go-figaro
package figaro

import "math/big"

// AddressSize is the size of an address, in bytes
const AddressSize = 32

// Address is an AddressSize length unique identifier
type Address []byte

// Account represents an account in Figaro
type Account struct {
	Nonce       *big.Int
	Address     []byte
	Balance     *big.Int
	Stake       *big.Int
	Code        []byte
	StorageRoot []byte
}

// AccountEncodingService should implement deterministic encoding/encoding of an account
type AccountEncodingService interface {
	EncodeAccount(account *Account) ([]byte, error)
	DecodeAccount(buf []byte) (*Account, error)
}

// AccountDataService should implement a Merkle database mapped to an account
type AccountDataService interface {
	// Account data services
	SaveAccount(ed AccountEncodingService, root []byte, account *Account) []byte
	FetchAccount(ed AccountEncodingService, root []byte, address [4]byte) (*Account, error)
	ProveAccount(ed AccountEncodingService, root []byte, address [4]byte) (*Account, [][][]byte, error)
	ValidateAccounted(ed AccountEncodingService, root []byte, account *Account, proof [][][]byte) bool

	// Account storage data services
	SaveAccountStorage(ed AccountEncodingService, oot []byte, account *Account, key, data []byte) []byte
	FetchAccountStorage(ed AccountEncodingService, account *Account, key []byte) ([]byte, error)
	ProveAccountStorage(ed AccountEncodingService, account *Account, key []byte) ([]byte, [][][]byte, error)
	ValidateAccountStorage(ed AccountEncodingService, account *Account, key, data []byte, proof [][][]byte) bool
}
