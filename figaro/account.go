// Package figaro is the main package for go-figaro
package figaro

import (
	"math/big"

	"github.com/figaro-tech/go-figaro/figbuf"
)

// Account represents an account in Figaro
type Account struct {
	Address     Address
	Nonce       uint64
	Stake       *big.Int
	Balance     *big.Int
	StorageRoot []byte
	Code        []byte
}

// Encode deterministically encodes an account to binary format.
func (acc Account) Encode() ([]byte, error) {
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	return enc.EncodeList(func(buf []byte) []byte {
		buf = enc.EncodeNextUint64(buf, acc.Nonce)
		buf = enc.EncodeNextTextMarshaler(buf, acc.Stake)
		buf = enc.EncodeNextTextMarshaler(buf, acc.Balance)
		buf = enc.EncodeNextBytes(buf, acc.Code)
		buf = enc.EncodeNextBytes(buf, acc.StorageRoot)
		return buf
	})
}

// Decode decodes a deterministically encoded account from binary format.
func (acc *Account) Decode(buf []byte) error {
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	acc = &Account{}
	return dec.DecodeList(buf, func(r []byte) []byte {
		acc.Nonce, r = dec.DecodeNextUint64(r)
		r = dec.DecodeNextTextUnmarshaler(r, acc.Stake)
		r = dec.DecodeNextTextUnmarshaler(r, acc.Balance)
		acc.Code, r = dec.DecodeNextBytes(r)
		acc.StorageRoot, r = dec.DecodeNextBytes(r)
		return r
	})
}

// AccountDataService should implement a Merkle database mapped to an account
type AccountDataService interface {
	// Account data services
	SaveAccount(root []byte, account *Account) []byte
	FetchAccount(root []byte, address Address) (*Account, error)
	ProveAccount(root []byte, address Address) (*Account, [][][]byte, error)
	ValidateAccount(root []byte, account *Account, proof [][][]byte) bool

	// Account storage data services
	SaveAccountStorage(root []byte, account *Account, key, data []byte) []byte
	FetchAccountStorage(account *Account, key []byte) ([]byte, error)
	ProveAccountStorage(account *Account, key []byte) ([]byte, [][][]byte, error)
	ValidateAccountStorage(account *Account, key, data []byte, proof [][][]byte) bool
}
