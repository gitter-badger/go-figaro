// Package figaro is the main package for go-figaro
package figaro

import (
	"github.com/figaro-tech/go-figaro/figbuf"
)

// MaxCodeSize is the max length, in bytes, of account code storage. This is
// a network configuration value, and does not impact consensus or validation
// of existing data.
const MaxCodeSize = 24576

// Account represents an account in Figaro
type Account struct {
	Address     Address
	Nonce       uint64
	Bonded      bool
	Stake       uint64
	Balance     uint64
	StorageRoot Root
	Code        []byte
}

// Encode deterministically encodes an account to binary format.
func (acc Account) Encode() ([]byte, error) {
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	return enc.EncodeList(func(buf []byte) []byte {
		buf = enc.EncodeNextUint64(buf, acc.Nonce)
		buf = enc.EncodeNextBool(buf, acc.Bonded)
		buf = enc.EncodeNextUint64(buf, acc.Stake)
		buf = enc.EncodeNextUint64(buf, acc.Balance)
		buf = enc.EncodeNextBytes(buf, acc.StorageRoot)
		buf = enc.EncodeNextBytes(buf, acc.Code)
		return buf
	})
}

// Decode decodes a deterministically encoded account from binary format.
func (acc *Account) Decode(buf []byte) error {
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	return dec.DecodeList(buf, func(r []byte) []byte {
		acc.Nonce, r = dec.DecodeNextUint64(r)
		acc.Bonded, r = dec.DecodeNextBool(r)
		acc.Stake, r = dec.DecodeNextUint64(r)
		acc.Balance, r = dec.DecodeNextUint64(r)
		acc.StorageRoot, r = dec.DecodeNextBytes(r)
		acc.Code, r = dec.DecodeNextBytes(r)
		return r
	})
}

// AccountLDataService can retreive data from either the local database
// or the p2p network.
type AccountLDataService interface {
	FetchAccount(root Root, address Address) (*Account, error)
	ProveAccount(root Root, address Address) (*Account, [][][]byte, error)
	ValidateAccount(root Root, account *Account, proof [][][]byte) bool

	FetchAccountStorage(account *Account, key []byte) ([]byte, error)
	ProveAccountStorage(account *Account, key []byte) ([]byte, [][][]byte, error)
	ValidateAccountStorage(account *Account, key, data []byte, proof [][][]byte) bool
}

// AccountDataService should implement a Merkle database mapped to an account.
type AccountDataService interface {
	AccountLDataService
	SaveAccount(root Root, account *Account) (Root, error)
	SaveAccountStorage(root Root, account *Account, key, data []byte) (Root, error)
}
