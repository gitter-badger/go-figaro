// Package figdb implements figaro domain specific wrappers for figdb
package figdb

import (
	"errors"

	"github.com/figaro-tech/go-figaro/figaro"
	"github.com/figaro-tech/go-figaro/figbuf"
	"github.com/figaro-tech/go-figaro/figdb"
)

// ErrInvalidAccount is returned when an address or db value for an account is not valid
var ErrInvalidAccount = errors.New("figdb account: invalid address or data")

// SaveAccount saves an account to the db, returning a new root
func (db *DB) SaveAccount(root []byte, account *figaro.Account) ([]byte, error) {
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	buf, err := enc.EncodeNextTextMarshaler(nil, account.Balance)
	if err != nil {
		return nil, err
	}
	buf, err = enc.EncodeNextTextMarshaler(buf, account.Stake)
	if err != nil {
		return nil, err
	}
	buf = enc.EncodeNextBytes(buf, account.Code[:])
	buf = enc.EncodeNextBytes(buf, account.StorageRoot[:])
	buf, err = enc.EncodeNextTextMarshaler(buf, account.Nonce)
	if err != nil {
		return nil, err
	}
	buf = enc.EncodeNextList(buf, 0)

	return db.State.Set(root, account.Address, buf)
}

// FetchAccount returns an account from the database
func (db *DB) FetchAccount(root []byte, address []byte) (*figaro.Account, error) {
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	b, err := db.State.Get(root, address)
	if err != nil {
		return nil, err
	}
	acc := &figaro.Account{Address: address}
	if len(b) == 0 {
		return acc, nil
	}
	l, r, err := dec.DecodeNextList(b)
	if err != nil {
		return nil, err
	}
	if len(r) > 0 {
		return nil, ErrInvalidAccount
	}
	r, err = dec.DecodeNextTextUnmarshaler(l, acc.Balance)
	if err != nil {
		return nil, err
	}
	r, err = dec.DecodeNextTextUnmarshaler(r, acc.Stake)
	if err != nil {
		return nil, err
	}
	acc.Code, r, err = dec.DecodeNextBytes(r)
	if err != nil {
		return nil, err
	}
	acc.StorageRoot, r, err = dec.DecodeNextBytes(r)
	if err != nil {
		return nil, err
	}
	r, err = dec.DecodeNextTextUnmarshaler(r, acc.Nonce)
	if err != nil {
		return nil, err
	}
	if len(r) > 0 {
		return nil, ErrInvalidAccount
	}
	return acc, nil
}

// ProveAccount returns an account from the database, along with a proof
func (db *DB) ProveAccount(root []byte, address []byte) (*figaro.Account, [][][]byte, error) {
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	b, proof, err := db.State.GetAndProve(root, address)
	if err != nil {
		return nil, nil, err
	}
	acc := &figaro.Account{Address: address}
	if len(b) == 0 {
		return acc, proof, nil
	}

	l, r, err := dec.DecodeNextList(b)
	if err != nil {
		return nil, nil, err
	}
	if len(r) > 0 {
		return nil, nil, ErrInvalidAccount
	}
	r, err = dec.DecodeNextTextUnmarshaler(l, acc.Balance)
	if err != nil {
		return nil, nil, err
	}
	r, err = dec.DecodeNextTextUnmarshaler(r, acc.Stake)
	if err != nil {
		return nil, nil, err
	}
	acc.Code, r, err = dec.DecodeNextBytes(r)
	if err != nil {
		return nil, nil, err
	}
	acc.StorageRoot, r, err = dec.DecodeNextBytes(r)
	if err != nil {
		return nil, nil, err
	}
	r, err = dec.DecodeNextTextUnmarshaler(r, acc.Nonce)
	if err != nil {
		return nil, nil, err
	}
	if len(r) > 0 {
		return nil, nil, ErrInvalidAccount
	}
	return acc, proof, nil
}

// ValidateAccount validates an account against a proof
func ValidateAccount(root []byte, account *figaro.Account, proof [][][]byte) bool {
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	buf, err := enc.EncodeNextTextMarshaler(nil, account.Balance)
	if err != nil {
		return false
	}
	buf, err = enc.EncodeNextTextMarshaler(buf, account.Stake)
	if err != nil {
		return false
	}
	buf = enc.EncodeNextBytes(buf, account.Code)
	buf = enc.EncodeNextBytes(buf, account.StorageRoot)
	buf, err = enc.EncodeNextTextMarshaler(buf, account.Nonce)
	if err != nil {
		return false
	}
	buf = enc.EncodeNextList(buf, 0)

	return figdb.StateValidate(root, account.Address, buf, proof)
}

// SaveAccountStorage saves binary key/value pair to the account's storage
//
// Requires passing the world state root as the first param, and returns the new world state
// root created as a result of the account storage root change.
func (db *DB) SaveAccountStorage(root []byte, account *figaro.Account, key, data []byte) ([]byte, error) {
	h, err := db.State.Set(account.StorageRoot, key, data)
	if err != nil {
		return nil, err
	}
	account.StorageRoot = h
	return db.SaveAccount(root, account)
}

// FetchAccountStorage fetches a value at key in the account storage root
func (db *DB) FetchAccountStorage(account *figaro.Account, key []byte) ([]byte, error) {
	return db.State.Get(account.StorageRoot, key)
}

// ProveAccountStorage fetches a value at key in the account storage root, and
// also returning a Merkle proof
func (db *DB) ProveAccountStorage(account *figaro.Account, key []byte) ([]byte, [][][]byte, error) {
	return db.State.GetAndProve(account.StorageRoot, key)
}

// ValidateAccountStorage validates a value at key in the account storage root against
// the Merkle proof
func ValidateAccountStorage(account *figaro.Account, key, data []byte, proof [][][]byte) bool {
	return figdb.StateValidate(account.StorageRoot, key, data, proof)
}
