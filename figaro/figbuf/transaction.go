// Package figbuf implements figaro domain specific wrappers for figbuf
package figbuf

import (
	"log"

	"github.com/figaro-tech/go-figaro/figaro"
	"github.com/figaro-tech/go-figaro/figbuf"
)

// EncodeTxCommit encodes
func (ed EncoderDecoder) EncodeTxCommit(tx *figaro.TxCommit) ([]byte, error) {
	return tx.TxHash, nil
}

// DecodeTxCommit decodes
func (ed EncoderDecoder) DecodeTxCommit(buf []byte) (*figaro.TxCommit, error) {
	return &figaro.TxCommit{TxHash: buf}, nil
}

// EncodeTransaction encodes
func (ed EncoderDecoder) EncodeTransaction(tx *figaro.Transaction) (buf []byte, err error) {
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
		buf = enc.EncodeNextBytes(buf, tx.Signature)
		buf = enc.EncodeNextBytes(buf, tx.Sender)
		buf = enc.EncodeNextBytes(buf, tx.To)
		buf = enc.EncodeNextTextMarshaler(buf, tx.Nonce)
		buf = enc.EncodeNextTextMarshaler(buf, tx.Stake)
		buf = enc.EncodeNextTextMarshaler(buf, tx.Value)
		buf = enc.EncodeNextBytes(buf, tx.Data)
		return buf
	})
	return buf, nil
}

// DecodeTransaction decodes
func (ed EncoderDecoder) DecodeTransaction(buf []byte) (tx *figaro.Transaction, err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				tx = nil
				err = re
			} else {
				log.Panic(r)
			}
		}
	}()
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	tx = &figaro.Transaction{}
	var r []byte
	r = dec.DecodeNextList(buf, func(b []byte) {
		tx.Signature, r = dec.DecodeNextBytes(b)
		tx.Sender, r = dec.DecodeNextBytes(r)
		tx.To, r = dec.DecodeNextBytes(r)
		r = dec.DecodeNextTextUnmarshaler(r, tx.Nonce)
		r = dec.DecodeNextTextUnmarshaler(r, tx.Stake)
		r = dec.DecodeNextTextUnmarshaler(r, tx.Value)
		tx.Data, r = dec.DecodeNextBytes(r)
		if len(r) > 0 {
			panic(ErrInvalidData)
		}
	})
	if len(r) > 0 {
		return nil, ErrInvalidData
	}
	return tx, nil
}
