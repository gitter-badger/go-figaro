// Package figbuf implements figaro domain specific wrappers for figbuf
package figbuf

import (
	"log"

	"github.com/figaro-tech/go-figaro/figaro"
	"github.com/figaro-tech/go-figaro/figbuf"
)

// EncodeAccount figbuf encodes an Account
func (ed EncoderDecoder) EncodeAccount(account *figaro.Account) (buf []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				buf = nil
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	buf = enc.EncodeNextList(nil, func(buf []byte) []byte {
		buf = enc.EncodeNextTextMarshaler(buf, account.Nonce)
		buf = enc.EncodeNextTextMarshaler(buf, account.Stake)
		buf = enc.EncodeNextTextMarshaler(buf, account.Balance)
		buf = enc.EncodeNextBytes(buf, account.Code)
		buf = enc.EncodeNextBytes(buf, account.StorageRoot)
		return buf
	})

	return buf, nil
}

// DecodeAccount figbuf decodes an Account buffer
func (ed EncoderDecoder) DecodeAccount(buf []byte) (acc *figaro.Account, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				acc = nil
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	acc = &figaro.Account{}
	var r []byte
	r = dec.DecodeNextList(buf, func(b []byte) {
		r = dec.DecodeNextTextUnmarshaler(b, acc.Nonce)
		r = dec.DecodeNextTextUnmarshaler(r, acc.Stake)
		r = dec.DecodeNextTextUnmarshaler(r, acc.Balance)
		acc.Code, r = dec.DecodeNextBytes(r)
		acc.StorageRoot, r = dec.DecodeNextBytes(r)
		if len(r) > 0 {
			panic(ErrInvalidData)
		}
	})
	if len(r) > 0 {
		return nil, ErrInvalidData
	}
	return acc, nil
}
