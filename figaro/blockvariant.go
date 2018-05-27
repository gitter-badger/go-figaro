// Package figaro is the main package for go-figaro
package figaro

// A RefBlock is a Block where Transactions is replacedwith TxIDs.
// Useful for requesting only missing Transactions.
type RefBlock struct {
	*BlockHeader
	Commits []Commit
	TxIDs   []TxHash
}

// Ref converts a Block into a RefBlock.
// The Block should already be sealed and signed before calling Ref.
func (bl Block) Ref() (rf *RefBlock) {
	rf.BlockHeader = bl.BlockHeader
	rf.Commits = bl.Commits
	rf.TxIDs = make([]TxHash, len(bl.Transactions))
	for i, t := range bl.Transactions {
		id, err := t.ID()
		if err != nil {
			panic("error getting transaction id")
		}
		rf.TxIDs[i] = id
	}
	return
}

// A CompBlock is a Block with only CommitsBloom and TxBloom.
// Useful for checking inclusion of transactions in light-clients.
type CompBlock struct {
	*BlockHeader
	CommitsBloom []byte
	TxBloom      []byte
}

// Compress converts a Block into a CompBlock.
// The Block should already be sealed and signed before calling Compress.
func (bl Block) Compress() (cb *CompBlock) {
	cb.BlockHeader = bl.BlockHeader
	cb.CommitsBloom = bl.CommitsBloom
	cb.TxBloom = bl.TxBloom
	return
}
