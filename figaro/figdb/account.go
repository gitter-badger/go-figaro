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
func (db DB) SaveAccount(root figaro.Root, account *figaro.Account) (newroot figaro.Root, err error) {
	var buf []byte
	buf, err = account.Encode()
	if err != nil {
		return
	}
	newroot, err = db.State.Set(root, account.Address, buf)
	if err != nil {
		return
	}
	return
}

// FetchAccount returns an account from the database
func (db DB) FetchAccount(root figaro.Root, address figaro.Address) (account *figaro.Account, err error) {
	var buf []byte
	buf, err = db.State.Get(root, address)
	if err != nil {
		return
	}
	if len(buf) > 0 {
		err = account.Decode(buf)
		if err != nil {
			return
		}
	}
	account.Address = address
	return
}

// ProveAccount returns an account from the database, along with a proof
func (db DB) ProveAccount(root figaro.Root, address figaro.Address) (account *figaro.Account, proof [][][]byte, err error) {
	var buf []byte
	buf, proof, err = db.State.GetAndProve(root, address)
	if err != nil {
		return
	}
	if len(buf) > 0 {
		err = account.Decode(buf)
		if err != nil {
			return
		}
	}
	account.Address = address
	return
}

// ValidateAccount validates an account against a proof.
func ValidateAccount(root figaro.Root, account *figaro.Account, proof [][][]byte) bool {
	buf, err := account.Encode()
	if err != nil {
		return false
	}
	return figdb.ValidateState(root, account.Address, buf, proof)
}

// SaveAccountStorage saves binary key/value pair to the account's storage.
// Requires passing the world state root as the first param, and returns the new
// world state root created as a result of the account storage root change.
func (db DB) SaveAccountStorage(root figaro.Root, account *figaro.Account, key, data []byte) (newroot figaro.Root, err error) {
	var storageroot []byte
	storageroot, err = db.State.Set(account.StorageRoot, key, data)
	if err != nil {
		return
	}
	account.StorageRoot = storageroot
	newroot, err = db.SaveAccount(root, account)
	return
}

// FetchAccountStorage fetches a value at key in the account storage root.
func (db DB) FetchAccountStorage(account *figaro.Account, key []byte) ([]byte, error) {
	return db.State.Get(account.StorageRoot, key)
}

// ProveAccountStorage fetches a value at key in the account storage root, and also returning a Merkle proof.
func (db DB) ProveAccountStorage(account *figaro.Account, key []byte) ([]byte, [][][]byte, error) {
	return db.State.GetAndProve(account.StorageRoot, key)
}

// ValidateAccountStorage validates a value at key in the account storage root against the Merkle proof.
func ValidateAccountStorage(account *figaro.Account, key, data []byte, proof [][][]byte) bool {
	return figdb.ValidateState(account.StorageRoot, key, data, proof)
}
