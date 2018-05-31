// Package figaro is the main package for go-figaro
package figaro

import (
	"errors"
	"time"

	"github.com/figaro-tech/go-figaro/figbuf"
	"github.com/figaro-tech/go-figaro/figcrypto/hasher"
	"github.com/figaro-tech/go-figaro/figcrypto/signature/fastsig"
)

var (
	// ErrInvalidTransaction is a self-explantory error
	ErrInvalidTransaction = errors.New("figaro tx: invalid transaction")
)

// MaxTxDataSize is the max length, in bytes, of tx data. This is
// a network configuration value, and does not impact consensus or validation
// of existing data.
const MaxTxDataSize = 4096

// A Transaction must be mined into a block. It contains transaction information along with a
// cryptographic signature over the Transaction hash by the sender, and a nonce value which
// must match the account nonce of the sender at the time it is mined into a block.
type Transaction struct {
	ID          []byte
	Signature   []byte
	From        Address
	To          Address
	Nonce       uint64
	Type        TxType
	CommitBlock uint64
	Value       uint64
	Data        []byte
}

// ToHash hashes the Tx fields other than Signature, creating a unique ID.
func (tx Transaction) ToHash() (TxHash, error) {
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	e, err := enc.Encode(tx.Nonce, tx.CommitBlock, tx.From, tx.To, tx.Type, tx.Value, tx.Data)
	if err != nil {
		return nil, err
	}
	return hasher.Hash256(e), nil
}

// Sign cryptographically signs a transaction.
func (tx *Transaction) Sign(privkey []byte) error {
	var err error
	tx.Signature, err = fastsig.Sign(privkey, tx.ID)
	return err
}

// VerifySignature verifies the address that signed the transaction.
func (tx Transaction) VerifySignature() bool {
	return fastsig.Verify(tx.From, tx.Signature, tx.ID)
}

// Encode deterministically encodes a transaction to binary format.
func (tx Transaction) Encode() ([]byte, error) {
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	return enc.EncodeList(func(buf []byte) []byte {
		buf = enc.EncodeNextBytes(buf, tx.Signature)
		buf = enc.EncodeNextUint64(buf, tx.Nonce)
		buf = enc.EncodeNextUint64(buf, tx.CommitBlock)
		buf = enc.EncodeNextBytes(buf, tx.From)
		buf = enc.EncodeNextBytes(buf, tx.To)
		buf = enc.EncodeNextBinaryMarshaler(buf, tx.Type)
		buf = enc.EncodeNextUint64(buf, tx.Value)
		buf = enc.EncodeNextBytes(buf, tx.Data)
		return buf
	})
}

// Decode decodes a deterministically encoded transaction from binary format.
func (tx *Transaction) Decode(buf []byte) error {
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	return dec.DecodeList(buf, func(r []byte) []byte {
		tx.Signature, r = dec.DecodeNextBytes(r)
		tx.Nonce, r = dec.DecodeNextUint64(r)
		tx.CommitBlock, r = dec.DecodeNextUint64(r)
		tx.From, r = dec.DecodeNextBytes(r)
		tx.To, r = dec.DecodeNextBytes(r)
		r = dec.DecodeNextBinaryUnmarshaler(r, &tx.Type)
		tx.Value, r = dec.DecodeNextUint64(r)
		tx.Data, r = dec.DecodeNextBytes(r)
		return r
	})
}

// TransactionLDataService implements only limited local data.
type TransactionLDataService interface {
	RetrieveTransactions(root Root) ([]*Transaction, error)
	GetTransaction(root Root, index int) (*Transaction, error)
	GetAndProveTransaction(root Root, index int) (transactions *Transaction, proof [][]byte, err error)
	ValidateTransaction(root Root, index int, commit Transaction, proof [][]byte) bool
}

// TransactionDataService provides merkelized data services for Transaction archives.
type TransactionDataService interface {
	ArchiveTransactions(transactions []*Transaction) (root Root, err error)
	TransactionLDataService
}

// A ReceivedTx is a Transaction that is waiting to be processed.
type ReceivedTx struct {
	Transaction

	Received time.Time
}

// Ordinal returns 1 if this is a future tx, 0 if it is a current tx, and -1 if this is a stale tx.
// Ordinality is only valid in reference to the given account nonce.
func (tx ReceivedTx) Ordinal(accnonce uint64) int {
	if tx.Nonce > accnonce {
		return 1
	}
	if tx.Nonce < accnonce {
		return -1
	}
	return 0
}

// TxReceivedHeap is a priority Heap of pending transactions, sorted by Received timestamp.
type TxReceivedHeap []*ReceivedTx

func (h TxReceivedHeap) Len() int           { return len(h) }
func (h TxReceivedHeap) Less(i, j int) bool { return h[i].Received.Before(h[j].Received) }
func (h TxReceivedHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push implements a heap.Interface. Use `heap.Push, etc`.
func (h *TxReceivedHeap) Push(x interface{}) {
	*h = append(*h, x.(*ReceivedTx))
}

// Pop implements a heap.Interface. Use `heap.Pop, etc`.
func (h *TxReceivedHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// TxNonceHeap is a is a priority Heap of transactions, sorted by Nonce.
type TxNonceHeap []*ReceivedTx

func (h TxNonceHeap) Len() int           { return len(h) }
func (h TxNonceHeap) Less(i, j int) bool { return h[i].Nonce < h[j].Nonce }
func (h TxNonceHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push implements a heap.Interface. Use `heap.Push, etc`.
func (h *TxNonceHeap) Push(x interface{}) {
	*h = append(*h, x.(*ReceivedTx))
}

// Pop implements a heap.Interface. Use `heap.Pop, etc`.
func (h *TxNonceHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
