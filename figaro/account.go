// Package figaro is the main package for go-figaro
package figaro

import "math/big"

// Account represents an account in Figaro
type Account struct {
	Nonce       *big.Int
	Address     []byte
	Balance     *big.Int
	Stake       *big.Int
	Code        []byte
	StorageRoot []byte
}

// AccountDataService should implement a Merkle database mapped to an account
type AccountDataService interface {
	SaveAccount(root []byte, account *Account) []byte
	FetchAccount(root []byte, address [4]byte) (*Account, error)
	ProveAccount(root []byte, address [4]byte) (*Account, [][][]byte, error)
	ValidateAccount(root []byte, account *Account, proof [][][]byte) bool
}

// AccountStorageDataService should implement a Merkle database mapped to an account's storage
type AccountStorageDataService interface {
	SaveAccountStorage(worldroot []byte, account *Account, key, data []byte) []byte
	FetchAccountStorage(account *Account, key []byte) ([]byte, error)
	ProveAccountStorage(account *Account, key []byte) ([]byte, [][][]byte, error)
	ValidateAccountStorage(account *Account, key, data []byte, proof [][][]byte) bool
}
