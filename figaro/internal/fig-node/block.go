package internal

import (
	"reflect"

	"github.com/figaro-tech/go-figaro/figaro"
	"github.com/figaro-tech/go-figaro/figaro/internal/fig-node/figdb"
)

// VerifyTxSignatures will verify all tx signatures in the block
// at once, returning whether even a single signature is fraudulent.
func VerifyTxSignatures(bl *figaro.Block) bool {
	// TODO: because we can do this without reference to the db,
	// we can run this concurrently and on machines with multiple
	// processors, in parallel
	for _, tx := range bl.Transactions {
		if !tx.VerifySignature() {
			return false
		}
	}
	return true
}

// SyncBlock will add all commits and transactions to the database, returning
// whether the block header is valid for the block data. If the block is invalid,
// it will unwind any changes.
func SyncBlock(db *figdb.DB, prev, bl *figaro.Block) error {
	db.FigDB.Store.Batch()
	defer db.FigDB.Store.Discard() // noop if Write is called upon success, otherwise will discard

	btest := &figaro.Block{
		BlockHeader: &figaro.BlockHeader{
			Signature:   bl.Signature,
			Producer:    bl.Producer,
			Beneficiary: bl.Beneficiary,
			Number:      bl.Number,
			ParentBlock: bl.ParentBlock,
			ChainConfig: bl.ChainConfig,
			StateRoot:   prev.StateRoot,
		},
	}
	btest.Commits = bl.Commits
	btest.Transactions = make([]*figaro.Transaction, 0, len(bl.Transactions))

	for i, tx := range bl.Transactions {
		cblockhash, err := db.FetchChainBlock(tx.CommitBlock)
		if err != nil {
			return err
		}
		cblock, err := db.FetchBlock(cblockhash)
		if err != nil {
			return err
		}
		var receipt *figaro.Receipt
		valid, err := ValidateTx(db, tx, btest.BlockHeader, cblock)
		if err != nil {
			return err
		}
		if valid {
			btest.StateRoot, receipt, err = ExecuteTx(db, tx, uint16(i), btest.BlockHeader, cblock.BlockHeader)
			if err != nil {
				return err
			}
		} else {
			btest.StateRoot, receipt, err = ExecuteInvalidTx(db, tx, uint16(i), btest.BlockHeader, cblock.BlockHeader)
			if err != nil {
				return err
			}
		}
		_, err = btest.AddTx(tx, receipt)
		if err != nil {
			return err
		}
	}
	btest.Seal(db)
	btest.Timestamp = bl.Timestamp
	if !reflect.DeepEqual(btest, bl) {
		return figaro.ErrInvalidBlock
	}
	return db.FigDB.Store.Write()
}

// ProduceBlock takes a freshly primed block and adds as many commits and transactions as possible, based
// on the pending pools, before sealing and signing the block.
func ProduceBlock(db *figdb.DB, bl *figaro.Block) {

}
