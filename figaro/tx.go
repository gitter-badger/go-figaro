// Package figaro is the main package for go-figaro
package figaro

import (
	"bytes"
	"errors"
	"log"
	"time"

	"github.com/figaro-tech/go-figaro/figbuf"
	"github.com/figaro-tech/go-figaro/figcrypto/hasher"
	"github.com/figaro-tech/go-figaro/figcrypto/signature/fastsig"
)

// MaxTxDataSize is the max length, in bytes, of tx data. This is
// a network configuration value, and does not impact consensus or validation
// of existing data.
const MaxTxDataSize = 4096

var (
	// ErrInvalidTxTypeData is a self-explantory error.
	ErrInvalidTxTypeData = errors.New("figaro tx: invalid TxType data")
)

// TxType is a supported transaction type. It is used to save
// space by only including a single `Value` field and a single
// byte to represent the type of value being transferred.
type TxType byte

const (
	// BalanceTx transactions transfer Fia Balance from one account to another.
	BalanceTx TxType = iota
	// StakeTx transactions transfer FIG Stake from one account to another.
	StakeTx
)

// ValidTxType is returns whether a TxType is a valid TxType
func ValidTxType(t TxType) bool {
	switch t {
	case BalanceTx:
		return true
	case StakeTx:
		return true
	default:
		return false
	}
}

// MarshalBinary implements encoding.BinaryMarshaler
func (tx TxType) MarshalBinary() ([]byte, error) {
	return []byte{byte(tx)}, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (tx *TxType) UnmarshalBinary(b []byte) error {
	if len(b) > 1 {
		return ErrInvalidTxTypeData
	}
	*tx = TxType(b[0])
	return nil
}

// A ReceivedTx is a Transaction that is waiting to be mined into a block.
type ReceivedTx struct {
	Transaction

	Received time.Time
}

// A Transaction must be mined into a block. It contains transaction information along with a
// cryptographic signature over the Transaction hash by the sender, and a nonce value which
// must match the account nonce of the sender at the time it is mined into a block.
type Transaction struct {
	Signature   []byte
	From        Address
	To          Address
	Nonce       uint64
	Type        TxType
	CommitBlock uint64
	Value       uint64
	Data        []byte

	// Sender is responsible for providing commit block,
	// this saves processing time by server, but also
	// helps ensure sender and receiver are on same
	// canonical chain

	// local data
	id []byte
}

// ID hashes the Tx fields other than Signature, creating a unique ID.
func (tx Transaction) ID() (txh TxHash, err error) {
	if len(tx.id) == 32 {
		txh = make(TxHash, len(tx.id))
		copy(txh, tx.id)
		return
	}

	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	var e []byte
	e, err = enc.Encode(tx.Nonce, tx.CommitBlock, tx.From, tx.To, tx.Type, tx.Value, tx.Data)
	if err != nil {
		return
	}
	txh = hasher.Hash256(e)
	tx.id = make([]byte, 32)
	copy(tx.id, txh)
	return
}

// Sign cryptographically signs a transaction.
func (tx *Transaction) Sign(privkey []byte) error {
	h, err := tx.ID()
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

// VerifySignature verifies the address that signed the transaction. It returns
// false if the transaction is not signed. If this check does not pass, the transaction
// should be rejected and no further processing is required.
func (tx Transaction) VerifySignature() bool {
	h, err := tx.ID()
	if err != nil {
		return false
	}
	return fastsig.Verify(tx.From, tx.Signature, h)
}

// QuickCheck returns whether the transaction is internally consistent and follows
// network protocol rules. As the name implies, it is meant to run quickly and should
// be called for "live" transactions prior to signature verification or tx validation,
// as tx rejected here save significant processing. CAUTION: because network configuration
// can change, it should NOT be used for historical validation.
func (tx Transaction) QuickCheck() bool {
	// Is a valid TxType
	if !ValidTxType(tx.Type) {
		return false
	}
	// Follows data limits
	if len(tx.Data) > MaxTxDataSize {
		return false
	}
	// Is signed
	if len(tx.Signature) < fastsig.SignatureSize {
		return false
	}
	return true
}

// Ordinal returns 1 if this is a future tx, 0 if it is a current tx, and -1 if this is a stale tx.
func (tx Transaction) Ordinal(db AccountFetchService, root Root) int {
	fromAcc, err := db.FetchAccount(root, tx.From)
	if err != nil {
		log.Panic(err)
	}
	if tx.Nonce > fromAcc.Nonce {
		return 1
	}
	if tx.Nonce < fromAcc.Nonce {
		return -1
	}
	return 0
}

// Validate returns whether the transaction will fail if it is processed as the next transaction.
func (tx Transaction) Validate(db FullChainDataService, commmitBlock *BigBlock, block *BlockHeader) bool {
	fromAcc, err := db.FetchAccount(block.StateRoot, tx.From)
	if err != nil {
		log.Panic(err)
	}
	// Sanity check
	if tx.CommitBlock != commmitBlock.Number {
		return false
	}
	// MPTx rules
	diffN := block.Number - commmitBlock.Number
	if diffN < uint64(block.WaitBlocks) || diffN > 2*uint64(block.WaitBlocks)+1 {
		return false
	}
	txid, err := tx.ID()
	if err != nil {
		return false
	}
	if !commmitBlock.CommitsBloom.Has(txid) {
		return false
	}
	for _, c := range commmitBlock.Commits {
		if bytes.Equal(c, txid) {
			break
		}
		return false
	}
	// commit fees use rules/beneficiary from commit block
	var totalFees uint32
	if !commmitBlock.Beneficiary.IsZeroAddress() {
		totalFees += commmitBlock.CommitFee
	}
	if !block.Beneficiary.IsZeroAddress() {
		totalFees += block.TxFee
	}
	// No free money (valid type check already occured)
	if tx.Type == StakeTx && (tx.Value > fromAcc.Stake || uint64(totalFees) > fromAcc.Balance) {
		return false
	} else if tx.Type == BalanceTx && (uint64(totalFees)+tx.Value) > fromAcc.Balance {
		return false
	}
	return true
}

// Execute executes a transaction, returning a transaction Receipt. It assumes
// that the transaction is valid for processing, and will perform no checks.
func (tx Transaction) Execute(db FullChainDataService, index uint16, block, commitBlock *BlockHeader) (newroot Root, receipt *Receipt, err error) {
	var fromAcc, toAcc, cbAcc, txbAcc *Account
	fromAcc, err = db.FetchAccount(block.StateRoot, tx.From)
	if err != nil {
		return
	}
	toAcc, err = db.FetchAccount(block.StateRoot, tx.To)
	if err != nil {
		return
	}
	fromAcc.Nonce++
	var totalFees uint32
	if !commitBlock.Beneficiary.IsZeroAddress() {
		cbAcc, err = db.FetchAccount(block.StateRoot, commitBlock.Beneficiary)
		if err != nil {
			return
		}
		cbAcc.Balance += uint64(commitBlock.CommitFee)
		fromAcc.Balance -= uint64(commitBlock.CommitFee)
		totalFees += commitBlock.CommitFee
	}
	if !block.Beneficiary.IsZeroAddress() {
		txbAcc, err = db.FetchAccount(block.StateRoot, block.Beneficiary)
		if err != nil {
			return
		}
		txbAcc.Balance += uint64(block.TxFee)
		fromAcc.Balance -= uint64(block.TxFee)
		totalFees += block.TxFee
	}
	if tx.Type == BalanceTx {
		fromAcc.Balance -= tx.Value
		toAcc.Balance += tx.Value
	} else if tx.Type == StakeTx {
		fromAcc.Stake -= tx.Value
		toAcc.Stake += tx.Value
	}
	// TODO: create contract or execute data against contract
	newroot = block.StateRoot
	newroot = db.SaveAccount(newroot, fromAcc)
	newroot = db.SaveAccount(newroot, toAcc)
	newroot = db.SaveAccount(newroot, cbAcc)
	newroot = db.SaveAccount(newroot, txbAcc)
	var h TxHash
	h, err = tx.ID()
	if err != nil {
		return
	}
	receipt = &Receipt{
		TxID:          h,
		BlockNum:      block.Number,
		Index:         index,
		PrevStateRoot: block.StateRoot,
		StateRoot:     newroot,
		TotalFees:     totalFees,
		Success:       true,
	}
	return
}

// ExecuteInvalid executes an invalid transaction, returning a transaction Receipt. It assumes
// that the transaction is invalid for processing, and will perform no checks. Invalid executions
// still pay fees to discourage spam txs, and still generate a receipt.
func (tx Transaction) ExecuteInvalid(db AccountDataService, index uint16, block, commitBlock *BlockHeader) (newroot Root, receipt *Receipt, err error) {
	var fromAcc, toAcc, cbAcc, txbAcc *Account
	fromAcc, err = db.FetchAccount(block.StateRoot, tx.From)
	if err != nil {
		return
	}
	toAcc, err = db.FetchAccount(block.StateRoot, tx.To)
	if err != nil {
		return
	}
	fromAcc.Nonce++
	var totalFees uint32
	var feeRatio float64
	if !commitBlock.Beneficiary.IsZeroAddress() {
		totalFees += commitBlock.CommitFee
	}
	if !block.Beneficiary.IsZeroAddress() {
		totalFees += block.TxFee
	}
	if uint64(totalFees) > fromAcc.Balance {
		feeRatio = float64(fromAcc.Balance) / float64(totalFees)
	} else {
		feeRatio = 1
	}
	if !commitBlock.Beneficiary.IsZeroAddress() {
		cbAcc, err = db.FetchAccount(block.StateRoot, commitBlock.Beneficiary)
		if err != nil {
			return
		}
		cfee := uint64(feeRatio * float64(commitBlock.CommitFee))
		if cfee > fromAcc.Balance {
			panic("invalid commit fee")
		}
		cbAcc.Balance += uint64(cfee)
		fromAcc.Balance -= uint64(cfee)

	}
	if !block.Beneficiary.IsZeroAddress() {
		txbAcc, err = db.FetchAccount(block.StateRoot, block.Beneficiary)
		if err != nil {
			return
		}
		txfee := uint64(feeRatio * float64(block.TxFee))
		if txfee > fromAcc.Balance {
			panic("invalid tx fee")
		}
		txbAcc.Balance += uint64(txfee)
		fromAcc.Balance -= uint64(txfee)
	}
	newroot = block.StateRoot
	newroot = db.SaveAccount(newroot, fromAcc)
	newroot = db.SaveAccount(newroot, toAcc)
	newroot = db.SaveAccount(newroot, cbAcc)
	newroot = db.SaveAccount(newroot, txbAcc)
	var h TxHash
	h, err = tx.ID()
	if err != nil {
		return
	}
	receipt = &Receipt{
		TxID:          h,
		BlockNum:      block.Number,
		Index:         index,
		PrevStateRoot: block.StateRoot,
		StateRoot:     newroot,
		TotalFees:     totalFees,
		Success:       false,
	}
	return
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
