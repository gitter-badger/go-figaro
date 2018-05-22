// Package figaro is the main package for go-figaro
package figaro

import (
	"container/heap"
	"math/big"
	"time"

	"github.com/figaro-tech/go-figaro/figbuf"
	"github.com/figaro-tech/go-figaro/figcrypto/hash"
	"github.com/figaro-tech/go-figaro/figcrypto/signature/fastsig"
)

// A TxCommit must be mined into a block.
type TxCommit []byte

// A Tx must be mined into a block.
type Tx struct {
	Signature []byte
	Nonce     uint64
	From      Address
	To        Address
	Stake     *big.Int
	Value     *big.Int
	Data      []byte
}

// TxID returns a cryptographic hash of Tx fields.
func (tx Tx) TxID() (h []byte, err error) {
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	var e []byte
	e, err = enc.Encode(tx.Nonce, tx.To, tx.Stake, tx.Value, tx.Data)
	if err != nil {
		return
	}
	h = hash.Hash256(e)
	return
}

// Sign cryptographically signs a signature.
func (tx *Tx) Sign(privkey []byte) error {
	h, err := tx.TxID()
	if err != nil {
		return err
	}
	sig, err := fastsig.Sign(privkey, h)
	if err != nil {
		return err
	}
	tx.Signature = sig
	return nil
}

// Verify verifies the address that signed the transaction. It returns
// false if the transaction is not signed.
func (tx Tx) Verify() bool {
	h, err := tx.TxID()
	if err != nil {
		return false
	}
	return fastsig.Verify(tx.From.Bytes(), tx.Signature, h)
}

// Encode deterministically encodes a transaction to binary format.
func (tx Tx) Encode() ([]byte, error) {
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	return enc.EncodeList(func(buf []byte) []byte {
		buf = enc.EncodeNextBytes(buf, tx.Signature)
		buf = enc.EncodeNextUint64(buf, tx.Nonce)
		buf = enc.EncodeNextBytes(buf, tx.To.Bytes())
		buf = enc.EncodeNextTextMarshaler(buf, tx.Stake)
		buf = enc.EncodeNextTextMarshaler(buf, tx.Value)
		buf = enc.EncodeNextBytes(buf, tx.Data)
		return buf
	})
}

// Decode decodes a deterministically encoded transaction from binary format.
func (tx *Tx) Decode(buf []byte) error {
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	tx = &Tx{}
	return dec.DecodeList(buf, func(r []byte) []byte {
		var toBytes []byte
		tx.Signature, r = dec.DecodeNextBytes(r)
		tx.Nonce, r = dec.DecodeNextUint64(r)
		toBytes, r = dec.DecodeNextBytes(r)
		err := tx.To.SetBytes(toBytes)
		if err != nil {
			panic(err)
		}
		r = dec.DecodeNextTextUnmarshaler(r, tx.Stake)
		r = dec.DecodeNextTextUnmarshaler(r, tx.Value)
		tx.Data, r = dec.DecodeNextBytes(r)
		return r
	})
}

// TxDataService provides merkelized data services for TxCommit and Tx histories
type TxDataService interface {
	// TxCommit data services
	CreateTxCommits(commits ...TxCommit) ([]byte, error)
	SaveTxCommits(key, set []byte) ([]byte, error)
	HasTxCommits(key []byte, commits ...TxCommit) ([]bool, error)
	HasTxCommit(key []byte, commit TxCommit) (bool, error)

	// Tx data services
	ArchiveTxs(commits ...*Tx) ([]byte, error)
	RetrieveTxs(root []byte) ([]*Tx, error)
	GetTx(root []byte, index int) (*Tx, error)
	GetAndProveTx(root []byte, index int) (*Tx, [][]byte, error)
	ValidateTx(root []byte, index int, commit *Tx, proof [][]byte) bool
}

// A PendingTxCommit implements the commit phase of MPTx.
type PendingTxCommit struct {
	TxCommit

	Received time.Time
}

// A PendingTx is a Tx that has been mined into a block.
type PendingTx struct {
	Tx

	TxID     []byte
	Received time.Time
}

// NewTxCommitHeap returns a TxCommitHeap, ready to use.
func NewTxCommitHeap() *TxCommitHeap {
	h := &TxCommitHeap{}
	heap.Init(h)
	return h
}

// NewTxHeap returns a TxHeap, ready to use.
func NewTxHeap() *TxHeap {
	h := &TxHeap{}
	heap.Init(h)
	return h
}

// TxCommitHeap is a min heap of pending tx commits. It implements `heap.Interface`.
type TxCommitHeap []*PendingTxCommit

func (h TxCommitHeap) Len() int           { return len(h) }
func (h TxCommitHeap) Less(i, j int) bool { return h[i].Received.Before(h[j].Received) }
func (h TxCommitHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push implements a heap.Interface. Use `heap.Push, etc`.
func (h *TxCommitHeap) Push(x interface{}) {
	*h = append(*h, x.(*PendingTxCommit))
}

// Pop implements a heap.Interface. Use `heap.Pop, etc`.
func (h *TxCommitHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// TxHeap is a priority Heap of pending transactions. It implements `heap.Interface`.
type TxHeap []*PendingTx

func (h TxHeap) Len() int           { return len(h) }
func (h TxHeap) Less(i, j int) bool { return h[i].Received.Before(h[j].Received) }
func (h TxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push implements a heap.Interface. Use `heap.Push, etc`.
func (h *TxHeap) Push(x interface{}) {
	*h = append(*h, x.(*PendingTx))
}

// Pop implements a heap.Interface. Use `heap.Pop, etc`.
func (h *TxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
