// Package figaro is the main package for go-figaro
package figaro

import "math/big"

// Account represents an account in Figaro
type Account struct {
	Nonce       *big.Int
	Address     [4]byte
	Balance     *big.Int
	Stake       *big.Int
	Code        []byte
	StorageRoot [4]byte
}

// AccountService should implement a Merkle database mapped to an account
type AccountService interface {
	SaveAccount(root []byte, account *Account) []byte
	FetchAccount(root []byte, address [4]byte) (*Account, error)
	ProveAccount(root []byte, address [4]byte) (*Account, [][][]byte, error)
	ValidateAccount(root []byte, account *Account, proof [][][]byte) bool
}

// AccountStorageService should implement a Merkle database mapped to an account's storage
type AccountStorageService interface {
	SaveAccountStorage(worldroot []byte, account *Account, key, data []byte) []byte
	FetchAccountStorage(account *Account, key []byte) ([]byte, error)
	ProveAccountStorage(account *Account, key []byte) ([]byte, [][][]byte, error)
	ValidateAccountStorage(account *Account, key, data []byte, proof [][][]byte) bool
}
