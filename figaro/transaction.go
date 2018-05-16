// Package figaro is the main package for go-figaro
package figaro

import (
	"container/heap"
	"math/big"
	"time"
)

// A TxCommit must be mined into a block.
type TxCommit struct {
	// A TxHash is a hash of valuable transaction fields. Concretely `H(Enc(Nonce, To, Stake, Value, Data))`.
	TxHash []byte
}

// A PendingTxCommit implements the commit phase of MPTx.
type PendingTxCommit struct {
	TxCommit

	Received time.Time
}

// A Transaction must be mined into a block.
type Transaction struct {
	Signature []byte
	To        string
	Nonce     *big.Int
	Stake     *big.Int
	Value     *big.Int
	Data      []byte
}

// A PendingTransaction is a transaction that has been mined into a block.
type PendingTransaction struct {
	Transaction

	Received time.Time
	// The (compressed) public key (33 bytes) should be provided in a transaction when
	// possible; however, it is not required in low-bandwidth situations and is not saved
	// to the database as the sender can be identified from the signature and tx hash alone.
	PublicKey []byte
}

// TransactionEncodingService should implement deterministic encoding/encoding of an account
type TransactionEncodingService interface {
	EncodeTxCommit(tx *TxCommit) ([]byte, error)
	DecodeTxCommit(buf []byte) (*TxCommit, error)

	EncodeTransaction(tx *Transaction) ([]byte, error)
	DecodeTransaction(buf []byte) (*Transaction, error)
}

// TransactionDataService provides merkelized data services for TxCommit and Transaction histories
type TransactionDataService interface {
	// TxCommit data services
	ArchiveTxCommits(ed TransactionEncodingService, commits ...*TxCommit) ([]byte, error)
	RetrieveTxCommits(ed TransactionEncodingService, root []byte) ([]*TxCommit, error)
	GetTxCommit(ed TransactionEncodingService, root []byte, index int) (*TxCommit, error)
	GetAndProveTxCommit(ed TransactionEncodingService, root []byte, index int) (*TxCommit, [][]byte, error)
	ValidateTxCommit(ed TransactionEncodingService, root []byte, index int, commit *TxCommit, proof [][]byte) bool

	// Transaction data services
	ArchiveTransactions(ed TransactionEncodingService, commits ...*Transaction) ([]byte, error)
	RetrieveTransactions(ed TransactionEncodingService, root []byte) ([]*Transaction, error)
	GetTransaction(ed TransactionEncodingService, root []byte, index int) (*Transaction, error)
	GetAndProveTransaction(ed TransactionEncodingService, root []byte, index int) (*Transaction, [][]byte, error)
	ValidateTransaction(ed TransactionEncodingService, root []byte, index int, commit *Transaction, proof [][]byte) bool
}

// NewTxCommitQueue returns a TxCommitQueue, ready to use.
func NewTxCommitQueue() *TxCommitQueue {
	pq := &TxCommitQueue{}
	heap.Init(pq)
	return pq
}

// NewTransactionQueue returns a TransactionQueue, ready to use.
func NewTransactionQueue() *TransactionQueue {
	pq := &TransactionQueue{}
	heap.Init(pq)
	return pq
}

// TxCommitQueue is a priority queue of pending tx commits
type TxCommitQueue []*PendingTxCommit

// Add adds a PendingTxCommit to the priority queue
func (pq *TxCommitQueue) Add(commit *PendingTxCommit) {
	heap.Push(pq, commit)
}

// Next gets the highest priority PendingTxCommit from the queue
func (pq *TxCommitQueue) Next() *PendingTxCommit {
	return heap.Pop(pq).(*PendingTxCommit)
}

// Len implements sort.Interface
func (pq TxCommitQueue) Len() int { return len(pq) }

// Less implements sort.Interface. Should not be called directly.
func (pq TxCommitQueue) Less(i, j int) bool {
	return pq[i].Received.After(pq[j].Received)
}

// Swap implements sort.Interface. Should not be called directly.
func (pq TxCommitQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// Push implements "container/heap". Should not be called directly.
func (pq *TxCommitQueue) Push(x interface{}) {
	item := x.(*PendingTxCommit)
	*pq = append(*pq, item)
}

// Pop implements "container/heap". Should not be called directly.
func (pq *TxCommitQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// TransactionQueue is a priority queue of pending transactions
type TransactionQueue []*PendingTransaction

// Add adds a PendingTransaction to the priority queue
func (pq *TransactionQueue) Add(commit *PendingTransaction) {
	heap.Push(pq, commit)
}

// Next gets the highest priority PendingTransaction from the queue
func (pq *TransactionQueue) Next() *PendingTransaction {
	return heap.Pop(pq).(*PendingTransaction)
}

// Len implements sort.Interface
func (pq TransactionQueue) Len() int { return len(pq) }

// Less implements sort.Interface. Should not be called directly.
func (pq TransactionQueue) Less(i, j int) bool {
	return pq[i].Received.After(pq[j].Received)
}

// Swap implements sort.Interface. Should not be called directly.
func (pq TransactionQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// Push implements "container/heap". Should not be called directly.
func (pq *TransactionQueue) Push(x interface{}) {
	item := x.(*PendingTransaction)
	*pq = append(*pq, item)
}

// Pop implements "container/heap". Should not be called directly.
func (pq *TransactionQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
