// Package figdb implements figaro domain specific wrappers for figdb
package figdb

import (
	"log"

	"github.com/figaro-tech/figaro"
)

// SaveAccount saves an account to the db, returning a new root
func (db *DB) SaveAccount(root []byte, account *figaro.Account) []byte {
	b, err := db.EncDec.Encode(account)
	if err != nil {
		log.Panic(err)
	}
	key := account.Address[:]
	return db.State.Set(root, key, b)
}

// FetchAccount returns an account from the database
func (db *DB) FetchAccount(root []byte, address [4]byte) *figaro.Account {
	key := address[:]
	b := db.State.Get(root, key)
	acc := &figaro.Account{}
	_, err := db.EncDec.Decode(b, acc)
	if err != nil {
		log.Panic(err)
	}
	return acc
}

// ProveAccount returns an account from the database, along with a proof
func (db *DB) ProveAccount(root []byte, address [4]byte) (*figaro.Account, [][][]byte) {
	key := address[:]
	b, proof := db.State.Prove(root, key)
	acc := &figaro.Account{}
	_, err := db.EncDec.Decode(b, acc)
	if err != nil {
		log.Panic(err)
	}
	return acc, proof
}

// ValidateAccount validates an account against a proof
func (v *Validator) ValidateAccount(root []byte, account *figaro.Account, proof [][][]byte) bool {
	b, err := v.EncDec.Encode(account)
	if err != nil {
		log.Panic(err)
	}
	key := account.Address[:]
	return v.State.Validate(root, key, b, proof)
}

// SaveAccountStorage saves binary key/value pair to the account's storage
//
// Requires passing the world state root as the first param, and returns the new world state
// root created as a result of the account storage root change.
func (db *DB) SaveAccountStorage(worldroot []byte, account *figaro.Account, key, data []byte) []byte {
	root := account.Address[:]
	copy(account.StorageRoot[:], db.State.Set(root, key, data))
	return db.SaveAccount(worldroot, account)
}

// FetchAccountStorage fetches a value at key in the account storage root
func (db *DB) FetchAccountStorage(account *figaro.Account, key []byte) ([]byte, error) {
	root := account.Address[:]
	b := db.State.Get(root, key)
	return b, nil
}

// ProveAccountStorage fetches a value at key in the account storage root, and
// also returning a Merkle proof
func (db *DB) ProveAccountStorage(account *figaro.Account, key []byte) ([]byte, [][][]byte, error) {
	root := account.Address[:]
	b, proof := db.State.Prove(root, key)
	return b, proof, nil
}

// ValidateAccountStorage validates a value at key in the account storage root against
// the Merkle proof
func (v *Validator) ValidateAccountStorage(account *figaro.Account, key, data []byte, proof [][][]byte) bool {
	root := account.Address[:]
	return v.State.Validate(root, key, data, proof)
}
