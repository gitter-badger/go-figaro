// Package figaro is the main package for go-figaro
package figaro

import (
	"container/heap"
	"time"
)

// A ReceivedCommit is a Commit that is waiting to be mined into a block.
type ReceivedCommit struct {
	Commit

	Received time.Time
}

// A Commit is sent prior to sending a transaction, according to a Wait/TTL block scheme. It must
// be mined into a block by creating set membership data that can be queried. It consists of a hash
// and is validated against a future transaction by comparing the hash to the Commit, which should
// be identical.
type Commit TxHash

// CommitLDataService implements limited local data commits.
type CommitLDataService interface {
	RetrieveCommits(root Root) ([]Commit, error)
	GetCommit(root Root, index int) (Commit, error)
	GetAndProveCommit(root Root, index int) (commit Commit, proof [][]byte, err error)
	ValidateCommit(root Root, index int, commit Commit, proof [][]byte) bool
}

// CommitDataService provides archive data service for commits.
type CommitDataService interface {
	ArchiveCommits(commits []Commit) (root Root, err error)
	CommitLDataService
}

// NewCommitHeap returns a CommitHeap, ready to use.
func NewCommitHeap() *CommitHeap {
	h := &CommitHeap{}
	heap.Init(h)
	return h
}

// CommitHeap is a min heap of pending tx commits. It implements `heap.Interface`. It
// implements a number of functions for sorting and managing.
type CommitHeap []*ReceivedCommit

func (h CommitHeap) Len() int           { return len(h) }
func (h CommitHeap) Less(i, j int) bool { return h[i].Received.Before(h[j].Received) }
func (h CommitHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push implements a heap.Interface. Use `heap.Push, etc`.
func (h *CommitHeap) Push(x interface{}) {
	*h = append(*h, x.(*ReceivedCommit))
}

// Pop implements a heap.Interface. Use `heap.Pop, etc`.
func (h *CommitHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
