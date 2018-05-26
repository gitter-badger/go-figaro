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

// MaxFees ensures that combined fees are not larger than the max value
// that can be represented by an individual fee.
const MaxFees = math.MaxUint32 / 2

var (
	// ErrExceedsBlockLimit is raised when a commit or tx is added when the addition would
	// exceed the max number of commits or transactions, which is MaxUint16.
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

// A CompBlock is a block where Commits and Transactions
// are replaced with bloom filters.
type CompBlock struct {
	*BlockHeader
	CommitsBloom []byte
	TxBloom      []byte
}

// A RefBlock is a block where Transactions
// is replaced with TxHashes. Useful where a requester
// already has a list of hashes, and will follow-up with
// a request for transaction data that they are missing.
type RefBlock struct {
	*BlockHeader
	Commits []Commit
	TxIDs   []TxHash
}

// Block is a collection of ordered transactions that updated state.
type Block struct {
	*BlockHeader
	Commits      []Commit
	Transactions []*Transaction
}

// BigBlock is a Block where Commits and Transactions co-exist alongside
// their bloombits, allowing faster membership testing at the cost of
// memory/storage.
type BigBlock struct {
	*BlockHeader
	CommitsBloom *bloom.Bloom
	Commits      []Commit
	TxBloom      *bloom.Bloom
	Transactions []*Transaction
}

// ID hashes the Block fields other than Signature, creating a unique ID.
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
	cbig, err := db.FetchBigBlock(cbid)
	if err != nil {
		return 0, err
	}
	var newroot Root
	var receipt *Receipt
	if tx.Validate(db, cbig, bl.BlockHeader) {
		newroot, receipt, err = tx.Execute(db, uint16(len(bl.Transactions)+1), bl.BlockHeader, cbig.BlockHeader)
		if err != nil {
			return 0, err
		}
	} else {
		newroot, receipt, err = tx.ExecuteInvalid(db, uint16(len(bl.Transactions)+1), bl.BlockHeader, cbig.BlockHeader)
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
	bl.Timestamp = time.Now()
	return nil
}

// CheckChainConfig returns whether the block chain config matches the provided chain config.
func (bl *Block) CheckChainConfig(config ChainConfig) bool {
	return reflect.DeepEqual(bl.ChainConfig, config)
}

// ValidateAndSync returns whether a block is valid or not. It relies on `engine` to determine whether
// the Producer is valid for the block. It then syncs all transactions and validates the Block against
// the results. The caller is responsible for any data cleanup if the Block does not validate.
func (bl *Block) ValidateAndSync(db FullChainDataService, prev *BlockHeader, engine ConsensusEngine) bool {
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
	if !engine.ValidateBlock(db, bl.BlockHeader) {
		return false
	}
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
	btest.Transactions = make([]*Transaction, len(bl.Transactions))

	// This will sync the block by reprocessing it locally
	for _, c := range bl.Commits {
		btest.AddCommit(c)
	}
	// for _, t := range bl.Transactions {
	// btest.AddTx(db, t) // houston, we have a problem
	// }
	btest.Seal(db)
	btest.Timestamp = bl.Timestamp
	return reflect.DeepEqual(btest, bl)
}

// Compress converts a Block into a CompBlock. The Block should already be sealed and signed before calling Compress.
func (bl Block) Compress() (cb *CompBlock, err error) {
	cb.BlockHeader = bl.BlockHeader
	cbloom := bloom.NewWithEstimates(uint64(len(bl.Commits)), bloomfp)
	for _, c := range bl.Commits {
		cbloom.Add(c)
	}
	var cbloombits []byte
	cbloombits, err = cbloom.Marshal()
	if err != nil {
		return
	}
	txbloom := bloom.NewWithEstimates(uint64(len(bl.Transactions)), bloomfp)
	var id TxHash
	for _, t := range bl.Transactions {
		id, err = t.ID()
		if err != nil {
			return
		}
		txbloom.Add(id)
	}
	var txloombits []byte
	txloombits, err = txbloom.Marshal()
	if err != nil {
		return
	}
	cb.CommitsBloom = cbloombits
	cb.TxBloom = txloombits
	return
}

// Ref converts a Block into a RefBlock. The Block should already be sealed and signed before calling Ref.
func (bl Block) Ref() (rf *RefBlock, err error) {
	rf.BlockHeader = bl.BlockHeader
	rf.Commits = make([]Commit, len(bl.Commits))
	copy(rf.Commits, bl.Commits)
	rf.TxIDs = make([]TxHash, len(bl.Transactions))
	var id TxHash
	for i, t := range bl.Transactions {
		id, err = t.ID()
		if err != nil {
			return
		}
		rf.TxIDs[i] = id
	}
	return
}

// Expand converts a Block into a BigBlock. The Block should already be sealed and signed before calling Expand.
func (bl Block) Expand() (bb *BigBlock, err error) {
	var cb *CompBlock
	cb, err = bl.Compress()
	if err != nil {
		return
	}
	bb.BlockHeader = bl.BlockHeader
	var cbloom, tbloom *bloom.Bloom
	cbloom, err = bloom.Unmarshal(cb.CommitsBloom)
	if err != nil {
		return
	}
	tbloom, err = bloom.Unmarshal(cb.TxBloom)
	if err != nil {
		return
	}
	bb.CommitsBloom = cbloom
	bb.Commits = bl.Commits
	bb.TxBloom = tbloom
	bb.Transactions = bl.Transactions
	return
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
	SaveBlock(bl *Block) error
	FetchBlockHeader(id BlockHash) (*BlockHeader, error)
}

// BlockDataService should save blocks directly into a key/value store.
type BlockDataService interface {
	FetchCompBlock(id BlockHash) (*CompBlock, error)
	FetchRefBlock(id BlockHash) (*RefBlock, error)
	FetchBlock(id BlockHash) (*Block, error)
	FetchBigBlock(id BlockHash) (*BigBlock, error)
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
