// Package figaro is the main package for go-figaro
package figaro

import (
	"github.com/figaro-tech/go-figaro/figbuf"
)

// Receipt is a record of a processed transaction.
type Receipt struct {
	TxID          TxHash
	BlockNum      uint64
	Index         uint16
	PrevStateRoot Root
	StateRoot     Root
	TotalFees     uint32
	Success       bool
}

// Encode encodes to binary.
func (rc Receipt) Encode() ([]byte, error) {
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	return enc.EncodeList(func(buf []byte) []byte {
		buf = enc.EncodeNextUint64(buf, rc.BlockNum)
		buf = enc.EncodeNextUint16(buf, rc.Index)
		buf = enc.EncodeNextBytes(buf, rc.PrevStateRoot)
		buf = enc.EncodeNextBytes(buf, rc.StateRoot)
		buf = enc.EncodeNextUint32(buf, rc.TotalFees)
		buf = enc.EncodeNextBool(buf, rc.Success)
		return buf
	})
}

// Decode decodes from binary.
func (rc *Receipt) Decode(buf []byte) error {
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	return dec.DecodeList(buf, func(r []byte) []byte {
		rc.BlockNum, r = dec.DecodeNextUint64(r)
		rc.Index, r = dec.DecodeNextUint16(r)
		rc.PrevStateRoot, r = dec.DecodeNextBytes(r)
		rc.StateRoot, r = dec.DecodeNextBytes(r)
		rc.TotalFees, r = dec.DecodeNextUint32(r)
		rc.Success, r = dec.DecodeNextBool(r)
		return r
	})
}

// ReceiptLDataService handles limited local storage of receipts.
type ReceiptLDataService interface {
	FetchReceipt(txid TxHash) (*Receipt, error)
}

// ReceiptDataService handles db storage of receipts.
type ReceiptDataService interface {
	SaveReceipt(r Receipt) error
	ReceiptLDataService
}
