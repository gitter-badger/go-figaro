package internal

import (
	"github.com/figaro-tech/go-fig-crypto/signature/fastsig"
	"github.com/figaro-tech/go-figaro/figaro"
	"github.com/figaro-tech/go-figaro/figaro/internal/fig-node/figdb"
)

// ValidateTx returns whether the transaction will fail if it is processed as the next transaction.
// Assumes that signature is already verified as authentic.
func ValidateTx(db *figdb.DB, tx *figaro.Transaction, txblock *figaro.BlockHeader, commitblock *figaro.Block) (bool, error) {
	fromAcc, err := db.FetchAccount(txblock.StateRoot, tx.From)
	if err != nil {
		return false, err
	}
	if fromAcc == nil {
		return false, nil
	}
	// Nonce must match
	if tx.Nonce != fromAcc.Nonce {
		return false, nil
	}
	// Is a valid TxType
	if !figaro.ValidTxType(tx.Type) {
		return false, nil
	}
	// Follows data limits
	if len(tx.Data) > figaro.MaxTxDataSize {
		return false, nil
	}
	// Is signed
	if len(tx.Signature) != fastsig.SignatureSize {
		return false, nil
	}
	// Sanity check
	if tx.CommitBlock != commitblock.Number {
		return false, nil
	}
	// MPTx rules
	diffN := txblock.Number - commitblock.Number
	if diffN < uint64(txblock.WaitBlocks) || diffN > 2*uint64(txblock.WaitBlocks)+1 {
		return false, nil
	}
	if !commitblock.HasCommit(tx.ID) {
		return false, nil
	}
	// No free money
	var totalFees uint32
	if !commitblock.Beneficiary.IsZeroAddress() {
		totalFees += commitblock.CommitFee
	}
	if !txblock.Beneficiary.IsZeroAddress() {
		totalFees += txblock.TxFee
	}
	switch tx.Type {
	case figaro.StakeTx:
		if tx.Value > fromAcc.Stake || uint64(totalFees) > fromAcc.Balance {
			return false, nil
		}
	case figaro.BalanceTx:
		if tx.Value+uint64(totalFees) > fromAcc.Balance {
			return false, nil
		}
	default:
		return false, nil
	}
	return true, nil
}

// ExecuteTx executes a transaction, returning a transaction Receipt.
// It assumes that the transaction is valid for processing, and will perform no checks.
func ExecuteTx(db *figdb.DB, tx *figaro.Transaction, index uint16, txblock, commitblock *figaro.BlockHeader) (figaro.Root, *figaro.Receipt, error) {
	var fromAcc, toAcc, cbAcc, txbAcc *figaro.Account
	fromAcc, err := db.FetchAccount(txblock.StateRoot, tx.From)
	if err != nil {
		return nil, nil, err
	}
	toAcc, err = db.FetchAccount(txblock.StateRoot, tx.To)
	if err != nil {
		return nil, nil, err
	}
	fromAcc.Nonce++
	var totalFees uint32
	if !commitblock.Beneficiary.IsZeroAddress() {
		cbAcc, err = db.FetchAccount(txblock.StateRoot, commitblock.Beneficiary)
		if err != nil {
			return nil, nil, err
		}
		cbAcc.Balance += uint64(commitblock.CommitFee)
		fromAcc.Balance -= uint64(commitblock.CommitFee)
		totalFees += commitblock.CommitFee
	}
	if !txblock.Beneficiary.IsZeroAddress() {
		txbAcc, err = db.FetchAccount(txblock.StateRoot, txblock.Beneficiary)
		if err != nil {
			return nil, nil, err
		}
		txbAcc.Balance += uint64(txblock.TxFee)
		fromAcc.Balance -= uint64(txblock.TxFee)
		totalFees += txblock.TxFee
	}
	switch tx.Type {
	case figaro.StakeTx:
		fromAcc.Stake -= tx.Value
		toAcc.Stake += tx.Value
	case figaro.BalanceTx:
		fromAcc.Balance -= tx.Value
		toAcc.Balance += tx.Value
	default:
		return nil, nil, err
	}

	// TODO: create contract or execute data against contract

	newroot := txblock.StateRoot
	newroot, err = db.SaveAccount(newroot, fromAcc)
	if err != nil {
		return nil, nil, err
	}
	newroot, err = db.SaveAccount(newroot, toAcc)
	if err != nil {
		return nil, nil, err
	}
	newroot, err = db.SaveAccount(newroot, cbAcc)
	if err != nil {
		return nil, nil, err
	}
	newroot, err = db.SaveAccount(newroot, txbAcc)
	if err != nil {
		return nil, nil, err
	}

	receipt := &figaro.Receipt{
		TxID:          tx.ID,
		BlockNum:      txblock.Number,
		Index:         index,
		PrevStateRoot: txblock.StateRoot,
		StateRoot:     newroot,
		TotalFees:     totalFees,
		Success:       true,
	}
	return newroot, receipt, nil
}

// ExecuteInvalidTx executes an invalid transaction, returning a transaction Receipt. It assumes
// that the transaction is invalid for processing, and will perform no checks. Invalid executions
// still pay fees to discourage spam txs, and still generate a receipt.
func ExecuteInvalidTx(db *figdb.DB, tx *figaro.Transaction, index uint16, txblock, commitblock *figaro.BlockHeader) (figaro.Root, *figaro.Receipt, error) {
	fromAcc, err := db.FetchAccount(txblock.StateRoot, tx.From)
	if err != nil {
		return nil, nil, err
	}
	fromAcc.Nonce++
	var totalFees uint32
	var feeRatio float64
	if !commitblock.Beneficiary.IsZeroAddress() {
		totalFees += commitblock.CommitFee
	}
	if !txblock.Beneficiary.IsZeroAddress() {
		totalFees += txblock.TxFee
	}
	if uint64(totalFees) > fromAcc.Balance {
		feeRatio = float64(fromAcc.Balance) / float64(totalFees)
	} else {
		feeRatio = 1
	}
	var cbAcc, txbAcc *figaro.Account
	if !commitblock.Beneficiary.IsZeroAddress() {
		cbAcc, err = db.FetchAccount(txblock.StateRoot, commitblock.Beneficiary)
		if err != nil {
			return nil, nil, err
		}
		cfee := uint64(feeRatio * float64(commitblock.CommitFee))
		if cfee > fromAcc.Balance {
			// this should never happen
			panic("invalid commit fee")
		}
		cbAcc.Balance += uint64(cfee)
		fromAcc.Balance -= uint64(cfee)
	}
	if !txblock.Beneficiary.IsZeroAddress() {
		txbAcc, err = db.FetchAccount(txblock.StateRoot, txblock.Beneficiary)
		if err != nil {
			return nil, nil, err
		}
		txfee := uint64(feeRatio * float64(txblock.TxFee))
		if txfee > fromAcc.Balance {
			// this should never happen
			panic("invalid tx fee")
		}
		txbAcc.Balance += uint64(txfee)
		fromAcc.Balance -= uint64(txfee)
	}
	newroot := txblock.StateRoot
	newroot, err = db.SaveAccount(newroot, fromAcc)
	if err != nil {
		return nil, nil, err
	}
	newroot, err = db.SaveAccount(newroot, cbAcc)
	if err != nil {
		return nil, nil, err
	}
	newroot, err = db.SaveAccount(newroot, txbAcc)
	if err != nil {
		return nil, nil, err
	}
	receipt := &figaro.Receipt{
		TxID:          tx.ID,
		BlockNum:      txblock.Number,
		Index:         index,
		PrevStateRoot: txblock.StateRoot,
		StateRoot:     newroot,
		TotalFees:     totalFees,
		Success:       false,
	}
	return newroot, receipt, nil
}
