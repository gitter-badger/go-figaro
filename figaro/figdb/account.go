// Package figdb implements figaro domain specific wrappers for figdb
package figdb

import (
	"errors"

	"github.com/figaro-tech/go-figaro/figaro"
	"github.com/figaro-tech/go-figaro/figdb"
)

// ErrInvalidAccount is returned when an address or db value for an account is not valid
var ErrInvalidAccount = errors.New("figdb account: invalid address")

// SaveAccount saves an account to the db, returning a new root
func (db DB) SaveAccount(ed figaro.AccountEncodingService, root []byte, account *figaro.Account) ([]byte, error) {
	buf, err := ed.EncodeAccount(account)
	if err != nil {
		return nil, err
	}
	return db.State.Set(root, account.Address, buf)
}

// FetchAccount returns an account from the database
func (db DB) FetchAccount(ed figaro.AccountEncodingService, root []byte, address []byte) (*figaro.Account, error) {
	buf, err := db.State.Get(root, address)
	if err != nil {
		return nil, err
	}
	acc, err := ed.DecodeAccount(buf)
	if err != nil {
		return nil, err
	}
	acc.Address = address
	return acc, nil
}

// ProveAccount returns an account from the database, along with a proof
func (db DB) ProveAccount(ed figaro.AccountEncodingService, root []byte, address []byte) (*figaro.Account, [][][]byte, error) {
	buf, proof, err := db.State.GetAndProve(root, address)
	if err != nil {
		return nil, nil, err
	}
	acc, err := ed.DecodeAccount(buf)
	if err != nil {
		return nil, nil, err
	}
	acc.Address = address
	return acc, proof, nil
}

// ValidateAccount validates an account against a proof.
func ValidateAccount(ed figaro.AccountEncodingService, root []byte, account *figaro.Account, proof [][][]byte) bool {
	buf, err := ed.EncodeAccount(account)
	if err != nil {
		return false
	}
	return figdb.StateValidate(root, account.Address, buf, proof)
}

// SaveAccountStorage saves binary key/value pair to the account's storage. Requires passing the world state root as the first param, and returns the new world state root created as a result of the account storage root change.
func (db DB) SaveAccountStorage(ed figaro.AccountEncodingService, root []byte, account *figaro.Account, key, data []byte) ([]byte, error) {
	h, err := db.State.Set(account.StorageRoot, key, data)
	if err != nil {
		return nil, err
	}
	account.StorageRoot = h
	return db.SaveAccount(ed, root, account)
}

// FetchAccountStorage fetches a value at key in the account storage root.
func (db DB) FetchAccountStorage(ed figaro.AccountEncodingService, account *figaro.Account, key []byte) ([]byte, error) {
	return db.State.Get(account.StorageRoot, key)
}

// ProveAccountStorage fetches a value at key in the account storage root, and also returning a Merkle proof.
func (db DB) ProveAccountStorage(ed figaro.AccountEncodingService, account *figaro.Account, key []byte) ([]byte, [][][]byte, error) {
	return db.State.GetAndProve(account.StorageRoot, key)
}

// ValidateAccountStorage validates a value at key in the account storage root against the Merkle proof.
func ValidateAccountStorage(ed figaro.AccountEncodingService, account *figaro.Account, key, data []byte, proof [][][]byte) bool {
	return figdb.StateValidate(account.StorageRoot, key, data, proof)
}
