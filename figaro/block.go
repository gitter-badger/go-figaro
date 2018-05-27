// Package figaro is the main package for go-figaro
package figaro

import (
	"container/heap"
	"errors"
	"math"
	"reflect"
	"time"

	"github.com/figaro-tech/go-figaro/figbuf"
	"github.com/figaro-tech/go-figaro/figcrypto/hasher"
	"github.com/figaro-tech/go-figaro/figcrypto/signature/fastsig"
	"github.com/figaro-tech/go-figaro/figcrypto/trie"
	"github.com/figaro-tech/go-figaro/figdb/bloom"
)

// MaxFees ensures that total fees doesn't overflow.
const MaxFees = math.MaxUint32 / 2

var (
	// ErrExceedsBlockLimit is returned when adding a commit or tx would overflow the
	// max commit or transaction size.
	ErrExceedsBlockLimit = errors.New("figaro block: commit or tx would exceed limit")
)

const bloomfp = 0.03

// BlockHeader is the header for a block.
type BlockHeader struct {
	Signature        []byte
	Provider         Address
	Beneficiary      Address
	Number           uint64
	Timestamp        time.Time
	ParentBlock      BlockHash
	StateRoot        Root
	CommitsRoot      Root
	TransactionsRoot Root
	ReceiptsRoot     Root

	ChainConfig

	// local data
	id       []byte
	receipts []*Receipt
}

// Block is a collection of ordered transactions that determine world state.
type Block struct {
	*BlockHeader
	CommitsBloom []byte
	TxBloom      []byte
	Commits      []Commit
	Transactions []*Transaction

	// local data
	cbloom  *bloom.Bloom
	txbloom *bloom.Bloom
}

// ID hashes the Block fields other than Signature, creating a unique ID. It is only
// valid for a sealed block.
func (bl *BlockHeader) ID() (bh BlockHash, err error) {
	if len(bl.id) == 32 {
		bh = make(BlockHash, len(bl.id))
		copy(bh, bl.id)
		return
	}

	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	var e []byte
	e, err = enc.Encode(
		bl.Provider,
		bl.Beneficiary,
		bl.Number,
		bl.Timestamp,
		bl.ParentBlock,
		bl.StateRoot,
		bl.CommitsRoot,
		bl.TransactionsRoot,
		bl.ReceiptsRoot,

		bl.Authority,
		bl.Stake,
		bl.CommitFee,
		bl.TxFee,
		bl.WaitBlocks,
	)
	if err != nil {
		return
	}
	bh = hasher.Hash256(e)
	bl.id = make([]byte, len(bh))
	copy(bl.id, bh)
	return
}

// Sign cryptographically signs a block. It should already be sealed before signing.
func (bl *BlockHeader) Sign(privkey []byte) error {
	h, err := bl.ID()
	if err != nil {
		return err
	}
	sig, err := fastsig.Sign(privkey, h)
	if err != nil {
		return err
	}
	bl.Signature = sig
	return nil
}

// VerifySignature verifies the address that signed the block. It returns
// false if the block is not signed.
func (bl BlockHeader) VerifySignature() bool {
	h, err := bl.ID()
	if err != nil {
		return false
	}
	return fastsig.Verify(bl.Provider, bl.Signature, h)
}

// AddCommit adds a transaction to the block, returning the total commits
// after the commit was added. Total commits cannot exceed `math.MaxUint16`.
func (bl *Block) AddCommit(c Commit) (int, error) {
	if len(bl.Commits) == math.MaxUint16 {
		return 0, ErrExceedsBlockLimit
	}
	bl.Commits = append(bl.Commits, c)
	return len(bl.Commits), nil
}

// AddTx exectutes and adds a transaction to the block, returning the total transactions
// after the tx was added. Total transactions cannot exceed `math.MaxUint16`.
func (bl *Block) AddTx(db FullChainDataService, tx *Transaction) (int, error) {
	if len(bl.Transactions) == math.MaxUint16 {
		return 0, ErrExceedsBlockLimit
	}
	// Fetch the commit block
	cbid, err := db.FetchChainBlock(tx.CommitBlock)
	if err != nil {
		return 0, err
	}
	cblock, err := db.FetchBlock(cbid)
	if err != nil {
		return 0, err
	}
	var newroot Root
	var receipt *Receipt
	if tx.Validate(db, bl.BlockHeader, cblock) {
		newroot, receipt, err = tx.Execute(db, uint16(len(bl.Transactions)+1), bl.BlockHeader, cblock.BlockHeader)
		if err != nil {
			return 0, err
		}
	} else {
		newroot, receipt, err = tx.ExecuteInvalid(db, uint16(len(bl.Transactions)+1), bl.BlockHeader, cblock.BlockHeader)
		if err != nil {
			return 0, err
		}
	}
	bl.StateRoot = newroot
	bl.Transactions = append(bl.Transactions, tx)
	bl.receipts = append(bl.receipts, receipt)
	return len(bl.Transactions), nil
}

// Seal processes the block. Headers should be initialized and all commits and transactions
// should have been added before calling Seal. It saves all commits, transactions, and receipts
// to the database, and updates the block header. After a block is sealed, it is ready to be signed.
func (bl *Block) Seal(db FullChainDataService) error {
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
	cbloom := bloom.NewWithEstimates(uint64(len(bl.Commits)), bloomfp)
	for _, c := range bl.Commits {
		cbloom.Add(c)
	}
	cbloombits, err := cbloom.Marshal()
	if err != nil {
		return err
	}
	txbloom := bloom.NewWithEstimates(uint64(len(bl.Transactions)), bloomfp)
	for _, t := range bl.Transactions {
		id, err := t.ID()
		if err != nil {
			return err
		}
		txbloom.Add(id)
	}
	txloombits, err := txbloom.Marshal()
	if err != nil {
		return err
	}
	bl.CommitsBloom = cbloombits
	bl.cbloom = cbloom
	bl.TxBloom = txloombits
	bl.txbloom = txbloom
	bl.Timestamp = time.Now()
	// ensure that ID() will set the cache the first time it is called after sealing,
	// otherwise an early call could cause sublte bugs
	bl.id = nil
	return nil
}

// HasCommit returns whether the CompBlock contains a Commit. Only valid on a sealed
// or received block.
func (bl *Block) HasCommit(c Commit) bool {
	if bl.cbloom == nil {
		var err error
		bl.cbloom, err = bloom.Unmarshal(bl.CommitsBloom)
		if err != nil {
			panic("invalid Block.CommitsBloom")
		}
	}
	return bl.cbloom.Has(c)
}

// HasTx returns whether the Block contains a Transaction. Only valid on a sealed
// or received block.
func (bl *Block) HasTx(txhash TxHash) bool {
	if bl.txbloom == nil {
		var err error
		bl.txbloom, err = bloom.Unmarshal(bl.TxBloom)
		if err != nil {
			panic("invalid Block.TxBloom")
		}
	}
	return bl.txbloom.Has(txhash)
}

// WellFormed returns whether the Block is well-formed. It does no further
// validation beyond checking that the Block could possibly be a valid Block.
func (bl *Block) WellFormed() bool {
	if len(bl.Signature) != fastsig.SignatureSize {
		return false
	}
	if !bl.Provider.Valid() ||
		!bl.Beneficiary.Valid() ||
		!bl.ParentBlock.Valid() ||
		!bl.StateRoot.Valid() ||
		!bl.CommitsRoot.Valid() ||
		!bl.TransactionsRoot.Valid() ||
		!bl.ReceiptsRoot.Valid() {
		return false
	}
	if len(bl.Commits) > math.MaxInt16 {
		return false
	}
	if len(bl.Transactions) > math.MaxInt16 {
		return false
	}
	return true
}

// Sync syncs all commits, transactions, and receipts and then validates the Blockheader,
// returning whether or not the Block was valid.
func (bl *Block) Sync(db FullChainDataService) bool {
	btest := &Block{
		BlockHeader: &BlockHeader{
			Signature:   bl.Signature,
			Provider:    bl.Provider,
			Beneficiary: bl.Beneficiary,
			Number:      bl.Number,
			ParentBlock: bl.ParentBlock,
			ChainConfig: ChainConfig{
				Authority:  bl.Authority,
				Stake:      bl.Stake,
				CommitFee:  bl.CommitFee,
				TxFee:      bl.TxFee,
				WaitBlocks: bl.WaitBlocks,
			},
		},
	}
	// preallocate for some performance gain
	btest.Commits = make([]Commit, 0, len(bl.Commits))
	btest.Transactions = make([]*Transaction, 0, len(bl.Transactions))
	// This will sync the block by reprocessing it locally
	for _, c := range bl.Commits {
		btest.AddCommit(c)
	}
	for _, t := range bl.Transactions {
		// We assume we have already validated supplied tx signatures elsewhere.
		btest.AddTx(db, t)
	}
	btest.Seal(db)
	btest.Timestamp = bl.Timestamp
	return reflect.DeepEqual(btest, bl)
}

// Encode deterministically encodes a Block to binary format.
func (bl BlockHeader) Encode() ([]byte, error) {
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	return enc.EncodeList(func(buf []byte) []byte {
		buf = enc.EncodeNextBytes(buf, bl.Signature)
		buf = enc.EncodeNextBytes(buf, bl.Provider)
		buf = enc.EncodeNextBytes(buf, bl.Beneficiary)
		buf = enc.EncodeNextUint64(buf, bl.Number)
		buf = enc.EncodeNextTextMarshaler(buf, bl.Timestamp)
		buf = enc.EncodeNextBytes(buf, bl.ParentBlock)
		buf = enc.EncodeNextBytes(buf, bl.StateRoot)
		buf = enc.EncodeNextBytes(buf, bl.CommitsRoot)
		buf = enc.EncodeNextBytes(buf, bl.TransactionsRoot)
		buf = enc.EncodeNextBytes(buf, bl.ReceiptsRoot)

		buf = enc.EncodeNextBytes(buf, bl.Authority)
		buf = enc.EncodeNextUint64(buf, bl.Stake)
		buf = enc.EncodeNextUint32(buf, bl.CommitFee)
		buf = enc.EncodeNextUint32(buf, bl.TxFee)
		buf = enc.EncodeNextUint8(buf, bl.WaitBlocks)
		return buf
	})
}

// Decode decodes a deterministically encoded Block from binary format.
func (bl *BlockHeader) Decode(buf []byte) error {
	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	return dec.DecodeList(buf, func(r []byte) []byte {
		bl.Signature, r = dec.DecodeNextBytes(r)
		bl.Provider, r = dec.DecodeNextBytes(r)
		bl.Beneficiary, r = dec.DecodeNextBytes(r)
		bl.Number, r = dec.DecodeNextUint64(r)
		r = dec.DecodeNextTextUnmarshaler(r, &bl.Timestamp)
		bl.ParentBlock, r = dec.DecodeNextBytes(r)
		bl.StateRoot, r = dec.DecodeNextBytes(r)
		bl.CommitsRoot, r = dec.DecodeNextBytes(r)
		bl.TransactionsRoot, r = dec.DecodeNextBytes(r)
		bl.ReceiptsRoot, r = dec.DecodeNextBytes(r)

		bl.Authority, r = dec.DecodeNextBytes(r)
		bl.Stake, r = dec.DecodeNextUint64(r)
		bl.CommitFee, r = dec.DecodeNextUint32(r)
		bl.TxFee, r = dec.DecodeNextUint32(r)
		bl.WaitBlocks, r = dec.DecodeNextUint8(r)
		return r
	})
}

// BlockLDataService implements limited local data blocks.
type BlockLDataService interface {
	SaveBlockHeader(bl *BlockHeader) error
	FetchBlockHeader(id BlockHash) (*BlockHeader, error)
}

// BlockDataService should save blocks directly into a key/value store.
type BlockDataService interface {
	SaveBlock(bl *Block) error
	FetchCompBlock(id BlockHash) (*CompBlock, error)
	FetchRefBlock(id BlockHash) (*RefBlock, error)
	FetchBlock(id BlockHash) (*Block, error)
	BlockLDataService
}

// NewBlockHeap returns a BlockHeap, ready to use.
func NewBlockHeap() *BlockHeap {
	h := &BlockHeap{}
	heap.Init(h)
	return h
}

// BlockHeap is a priority heap of blocks sorted by block number. It implements `heap.Interface`.  It
// implements a number of functions for sorting and managing.
type BlockHeap []*Block

// PeekNextNumber returns the next block index on the heap without modifying the heap.
func (h BlockHeap) PeekNextNumber() uint64 {
	return h[len(h)-1].Number
}

func (h BlockHeap) Len() int           { return len(h) }
func (h BlockHeap) Less(i, j int) bool { return h[i].Number < h[j].Number }
func (h BlockHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push implements a heap.Interface. Use `heap.Push, etc`.
func (h *BlockHeap) Push(x interface{}) {
	*h = append(*h, x.(*Block))
}

// Pop implements a heap.Interface. Use `heap.Pop, etc`.
func (h *BlockHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
