// Package figaro is the main package for go-figaro
package figaro

import "github.com/figaro-tech/go-figaro/figdb/bloom"

// A RefBlock is a block where Transactions
// is replaced with TxHashes. Useful where a requester
// already has a list of transactions, and will follow-up with
// a request for transaction data that they are missing.
type RefBlock struct {
	*BlockHeader
	Commits []Commit
	TxIDs   []TxHash
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

// A CompBlock is a block where Commits and Transactions
// are replaced with bloom filters.
type CompBlock struct {
	*BlockHeader
	CommitsBloom []byte
	TxBloom      []byte

	// local data
	cbloom  *bloom.Bloom
	txbloom *bloom.Bloom
}

// HasCommit returns whether the CompBlock contains a Commit.
func (cb *CompBlock) HasCommit(c Commit) bool {
	if cb.cbloom == nil {
		var err error
		cb.cbloom, err = bloom.Unmarshal(cb.CommitsBloom)
		if err != nil {
			panic("invalid CompBlock.CommitsBloom")
		}
	}
	return cb.cbloom.Has(c)
}

// HasTx returns whether the CompBlock contains a Transaction.
func (cb *CompBlock) HasTx(txhash TxHash) bool {
	if cb.txbloom == nil {
		var err error
		cb.txbloom, err = bloom.Unmarshal(cb.TxBloom)
		if err != nil {
			panic("invalid CompBlock.TxBloom")
		}
	}
	return cb.txbloom.Has(txhash)
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
	cb.cbloom = cbloom
	cb.TxBloom = txloombits
	cb.txbloom = txbloom
	return
}

// BigBlock is a Block where Commits and Transactions co-exist alongside
// their bloombits, allowing faster membership testing at the cost of
// memory/storage.
type BigBlock struct {
	*BlockHeader
	*CompBlock
	*Block
	TxIDs []TxHash
}

// Expand converts a Block into a BigBlock. The Block should already be sealed and signed before calling Expand.
func (bl *Block) Expand() (*BigBlock, error) {
	cb, err := bl.Compress()
	if err != nil {
		return nil, err
	}
	ref, err := bl.Ref()
	if err != nil {
		return nil, err
	}
	return &BigBlock{bl.BlockHeader, cb, bl, ref.TxIDs}, nil
}
