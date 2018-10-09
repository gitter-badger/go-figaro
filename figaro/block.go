// Package figaro is the main package for go-figaro
package figaro

import (
	"bytes"
	"container/heap"
	"errors"
	"math"
	"time"

	"github.com/figaro-tech/go-fig-buf"
	"github.com/figaro-tech/go-fig-crypto/hasher"
	"github.com/figaro-tech/go-fig-crypto/signature/fastsig"
	"github.com/figaro-tech/go-fig-crypto/trie"
	"github.com/figaro-tech/go-fig-db/bloom"
)

const (
	// MaxCommitSize is the max number of commits for a block
	MaxCommitSize = math.MaxUint16
	// MaxTxSize is the max number of commits for a block
	MaxTxSize = math.MaxUint16
	// MaxFees is the maximum amount of commit or transaction fees, ensuring that
	// total fees does not overflow
	MaxFees = math.MaxUint32 / 2
)

var (
	// ErrExceedsBlockLimit is returned when adding a commit or tx would overflow the
	// max commit or transaction size.
	ErrExceedsBlockLimit = errors.New("figaro block: commit or tx would exceed limit")

	// protocol block fp = 0.03
	bloomfp = 0.03
)

// BlockHeader is the header for a block.
type BlockHeader struct {
	ID               []byte
	Signature        []byte
	Producer         Address
	Beneficiary      Address
	Number           uint64
	Timestamp        time.Time
	ParentBlock      BlockHash
	StateRoot        Root
	CommitsRoot      Root
	TransactionsRoot Root
	ReceiptsRoot     Root

	ChainConfig
}

// ToHash hashes the Block fields other than Signature, creating a unique ID.
// It is only valid for a sealed block.
func (bl *BlockHeader) ToHash() (BlockHash, error) {
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	cfg, err := bl.ChainConfig.Encode()
	if err != nil {
		return nil, err
	}
	e, err := enc.Encode(
		bl.Producer,
		bl.Beneficiary,
		bl.Number,
		bl.Timestamp,
		bl.ParentBlock,
		bl.StateRoot,
		bl.CommitsRoot,
		bl.TransactionsRoot,
		bl.ReceiptsRoot,
		cfg,
	)
	if err != nil {
		return nil, err
	}
	return hasher.Hash256(e), nil
}

// Sign cryptographically signs the block with the given
// private key, updating the block header.
func (bl *BlockHeader) Sign(privkey []byte) error {
	var err error
	bl.Signature, err = fastsig.Sign(privkey, bl.ID)
	return err
}

// VerifySignature verifies that the given block signature
// matches the block provider.
func (bl *BlockHeader) VerifySignature() bool {
	return fastsig.Verify(bl.Producer, bl.Signature, bl.ID)
}

// Encode deterministically encodes a BlockHeader to binary format.
func (bl BlockHeader) Encode() ([]byte, error) {
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	return enc.EncodeList(func(buf []byte) []byte {
		buf = enc.EncodeNextBytes(buf, bl.Signature)
		buf = enc.EncodeNextBytes(buf, bl.Producer)
		buf = enc.EncodeNextBytes(buf, bl.Beneficiary)
		buf = enc.EncodeNextUint64(buf, bl.Number)
		buf = enc.EncodeNextTextMarshaler(buf, bl.Timestamp)
		buf = enc.EncodeNextBytes(buf, bl.ParentBlock)
		buf = enc.EncodeNextBytes(buf, bl.StateRoot)
		buf = enc.EncodeNextBytes(buf, bl.CommitsRoot)
		buf = enc.EncodeNextBytes(buf, bl.TransactionsRoot)
		buf = enc.EncodeNextBytes(buf, bl.ReceiptsRoot)
		cfg, err := bl.ChainConfig.Encode()
		if err != nil {
			panic(err)
		}
		buf = enc.EncodeNextBytes(buf, cfg)
		return buf
	})
}

// Decode decodes a deterministically encoded BlockHeader from binary format.
func (bl *BlockHeader) Decode(buf []byte) error {
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	return dec.DecodeList(buf, func(r []byte) []byte {
		bl.Signature, r = dec.DecodeNextBytes(r)
		bl.Producer, r = dec.DecodeNextBytes(r)
		bl.Beneficiary, r = dec.DecodeNextBytes(r)
		bl.Number, r = dec.DecodeNextUint64(r)
		r = dec.DecodeNextTextUnmarshaler(r, &bl.Timestamp)
		bl.ParentBlock, r = dec.DecodeNextBytes(r)
		bl.StateRoot, r = dec.DecodeNextBytes(r)
		bl.CommitsRoot, r = dec.DecodeNextBytes(r)
		bl.TransactionsRoot, r = dec.DecodeNextBytes(r)
		bl.ReceiptsRoot, r = dec.DecodeNextBytes(r)
		var cfg []byte
		cfg, r = dec.DecodeNextBytes(r)
		err := bl.ChainConfig.Decode(cfg)
		if err != nil {
			panic(err)
		}
		return r
	})
}

// Block is a collection of ordered transactions that determine world state.
type Block struct {
	*BlockHeader
	CommitsBloom []byte
	Commits      []Commit
	TxBloom      []byte
	Transactions []*Transaction

	// local data
	cbloom   *bloom.Bloom
	txbloom  *bloom.Bloom
	receipts []*Receipt
}

// AddCommit adds a commit to the block, returning the total commits
// after the commit was added. Total commits cannot exceed `math.MaxUint16`.
func (bl *Block) AddCommit(commit Commit) (int, error) {
	if len(bl.Commits) == math.MaxUint16 {
		return 0, ErrExceedsBlockLimit
	}
	bl.Commits = append(bl.Commits, commit)
	return len(bl.Commits), nil
}

// HasCommit returns whether the txhash has a commit in the block.
func (bl *Block) HasCommit(txhash TxHash) bool {
	if !bl.cbloom.Has(txhash) {
		return false
	}
	for _, c := range bl.Commits {
		if bytes.Equal(txhash, c) {
			return true
		}
	}
	return false
}

// AddTx adds a transaction to the block, returning the total transactions
// after the tx was added. Total transactions cannot exceed `math.MaxUint16`.
func (bl *Block) AddTx(tx *Transaction, receipt *Receipt) (int, error) {
	if len(bl.Transactions) == math.MaxUint16 {
		return 0, ErrExceedsBlockLimit
	}
	bl.Transactions = append(bl.Transactions, tx)
	bl.receipts = append(bl.receipts, receipt)
	return len(bl.Transactions), nil
}

// HasTx returns whether the txhash has a transaction in the block.
func (bl *Block) HasTx(txhash TxHash) bool {
	if !bl.txbloom.Has(txhash) {
		return false
	}
	for _, tx := range bl.Transactions {
		if bytes.Equal(txhash, tx.ID) {
			return true
		}
	}
	return false
}

// SetBlooms recalculates and sets the block bloom filters. Should be called
// after transactions are added, or when hydrating a block from storage.
func (bl *Block) SetBlooms() error {
	var err error
	bl.cbloom = bloom.NewWithEstimates(uint64(len(bl.Commits)), bloomfp)
	bl.txbloom = bloom.NewWithEstimates(uint64(len(bl.Transactions)), bloomfp)
	for _, c := range bl.Commits {
		bl.cbloom.Add(c)
	}
	for _, t := range bl.Transactions {
		bl.txbloom.Add(t.ID)
	}
	bl.CommitsBloom, err = bl.cbloom.Marshal()
	if err != nil {
		return err
	}
	bl.TxBloom, err = bl.txbloom.Marshal()
	if err != nil {
		return err
	}
	return nil
}

// Seal seal as a block, writing commits, transactions, and receipts
// to the database and updating the block. Once sealed, a block is ready to be signed.
func (bl *Block) Seal(db BlockContentsDataService) error {
	var err error
	bl.CommitsRoot, err = db.ArchiveCommits(bl.Commits)
	if err != nil {
		return err
	}
	bl.TransactionsRoot, err = db.ArchiveTransactions(bl.Transactions)
	if err != nil {
		return err
	}
	var br []byte
	binreceipts := make([][]byte, len(bl.receipts))
	for i, r := range bl.receipts {
		err = db.SaveReceipt(*r)
		if err != nil {
			return err
		}
		br, err = r.Encode()
		if err != nil {
			return err
		}
		binreceipts[i] = br
	}
	bl.ReceiptsRoot = trie.Trie(binreceipts)
	if err != nil {
		return err
	}
	err = bl.SetBlooms()
	if err != nil {
		return err
	}
	bl.Timestamp = time.Now()
	return nil
}

// Encode deterministically encodes a Block to binary format.
// This is used for communication between nodes.
func (bl Block) Encode() ([]byte, error) {
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	return enc.EncodeList(func(buf []byte) []byte {
		head, err := bl.BlockHeader.Encode()
		if err != nil {
			panic(err)
		}
		buf = enc.EncodeNextBytes(buf, head)
		buf = enc.EncodeNextBytes(buf, bl.CommitsBloom)
		buf = enc.EncodeNextBytes(buf, bl.TxBloom)
		buf = enc.EncodeNextList(buf, func(buf []byte) []byte {
			for _, c := range bl.Commits {
				buf = enc.EncodeNextBytes(buf, c)
			}
			return buf
		})
		buf = enc.EncodeNextList(buf, func(buf []byte) []byte {
			for _, t := range bl.Transactions {
				e, err := t.Encode()
				if err != nil {
					panic(err)
				}
				buf = enc.EncodeNextBytes(buf, e)
			}
			return buf
		})
		return buf
	})
}

// Decode decodes a deterministically encoded Block from binary format.
// This is used for communication between nodes.
func (bl *Block) Decode(buf []byte) error {
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	return dec.DecodeList(buf, func(r []byte) []byte {
		var head []byte
		head, r = dec.DecodeNextBytes(r)
		err := bl.BlockHeader.Decode(head)
		if err != nil {
			panic(err)
		}
		bl.CommitsBloom, r = dec.DecodeNextBytes(r)
		bl.TxBloom, r = dec.DecodeNextBytes(r)
		r = dec.DecodeNextList(r, func(r []byte) []byte {
			var c []byte
			for len(r) > 0 {
				c, r = dec.DecodeNextBytes(r)
				bl.Commits = append(bl.Commits, c)
			}
			return r
		})
		r = dec.DecodeNextList(r, func(r []byte) []byte {
			var e []byte

			for len(r) > 0 {
				var t *Transaction
				e, r = dec.DecodeNextBytes(r)
				err := t.Decode(e)
				if err != nil {
					panic(err)
				}
				bl.Transactions = append(bl.Transactions, t)
			}
			return r
		})
		return r
	})
}

// BlockContentsDataService is a data service that can support commits, transactions, and receipts.
type BlockContentsDataService interface {
	CommitDataService
	TransactionDataService
	ReceiptDataService
}

// BlockDataService should save blocks directly into a key/value store.
type BlockDataService interface {
	SaveBlock(bl *Block) error
	FetchBlockHeader(id BlockHash) (*BlockHeader, error)
	FetchCompBlock(id BlockHash) (*CompBlock, error)
	FetchRefBlock(id BlockHash) (*RefBlock, error)
	FetchBlock(id BlockHash) (*Block, error)
}

// NewBlockHeap returns a BlockHeap, ready to use.
func NewBlockHeap() *BlockHeap {
	h := &BlockHeap{}
	heap.Init(h)
	return h
}

// BlockHeap is a priority heap of blocks sorted by block number. It implements `heap.Interface`.  It
// implements a number of functions for sorting and managing.
type BlockHeap []*BlockHeader

// PeekNextNumber returns the next block index on the heap without modifying the heap.
func (h BlockHeap) PeekNextNumber() uint64 {
	return h[len(h)-1].Number
}

func (h BlockHeap) Len() int           { return len(h) }
func (h BlockHeap) Less(i, j int) bool { return h[i].Number < h[j].Number }
func (h BlockHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push implements a heap.Interface. Use `heap.Push, etc`.
func (h *BlockHeap) Push(x interface{}) {
	*h = append(*h, x.(*BlockHeader))
}

// Pop implements a heap.Interface. Use `heap.Pop, etc`.
func (h *BlockHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
