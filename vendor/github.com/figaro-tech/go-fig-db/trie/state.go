// Package trie implements merkle tries over the database
package trie

import (
	"bytes"
	"errors"
	"sync"

	"github.com/figaro-tech/go-fig-buf"
	"github.com/figaro-tech/go-fig-crypto/hasher"
	"github.com/figaro-tech/go-fig-db/types"
)

var (
	// ErrInvalidCompactEncoding is a self-explantory error
	ErrInvalidCompactEncoding = errors.New("figdb trie state: invalid compact encoding for path")
)

// State implements a Merkle Patricia trie over a DB
// for data that is updated often
type State struct {
	lock     sync.RWMutex
	KeyStore types.KeyStore
	Cache    types.Cache
	branches [32][17][]byte
	nodes    [8][2][]byte
	brat     int
	noat     int
}

func (tr *State) getNewBranch() [][]byte {
	if tr.brat == len(tr.branches) {
		// If we're at the end of the pool,
		// we allocate some more and let the GC
		// collect the old ones
		var fl [32][17][]byte
		tr.branches = fl
		tr.brat = 0
	}
	// Grab one from the pool and set the values
	bs := tr.branches[tr.brat]
	tr.brat++
	return bs[:]
}

func (tr *State) getNewNode(a, b []byte) [][]byte {
	if tr.noat == len(tr.nodes) {
		// If we're at the end of the pool,
		// we allocate some more and let the GC
		// collect the old ones
		var fl [8][2][]byte
		tr.nodes = fl
		tr.noat = 0
	}
	// Grab one from the pool and set the values
	bs := tr.nodes[tr.noat]
	bs[0] = a
	bs[1] = b
	tr.noat++
	return bs[:]
}

// Set updates a key/value pair from a given Merkle root,
// returning the new Merkle root containing all state
func (tr *State) Set(root, key, value []byte) ([]byte, error) {
	tr.lock.Lock()
	defer tr.lock.Unlock()

	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	tr.KeyStore.Batch()
	defer tr.KeyStore.Write()

	path := nibbles(key)
	if len(root) == 0 {
		return tr.setNilRoot(enc, path, value)
	}
	return tr.set(enc, dec, root, path, value)
}

// SetInBatch updates a key/value pair from a given Merkle root,
// returning the new Merkle root containing all state
//
// It does not batch the inserts, because it assumes this is part
// of a larger batch of updates
func (tr *State) SetInBatch(root, key, value []byte) ([]byte, error) {
	tr.lock.Lock()
	defer tr.lock.Unlock()

	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	path := nibbles(key)
	if len(root) == 0 {
		return tr.setNilRoot(enc, path, value)
	}

	return tr.set(enc, dec, root, path, value)
}

func (tr *State) set(enc *figbuf.Encoder, dec *figbuf.Decoder, root []byte, path []uint8, value []byte) ([]byte, error) {
	node, err := tr.getNode(dec, root)
	if err != nil {
		return nil, err
	}
	if len(node) == 0 {
		return tr.setNilNode(enc, path, value)
	}
	if len(node) == 17 {
		return tr.setBranchNode(enc, dec, node, path, value)
	}
	return tr.setLeafOrExtension(enc, dec, node, path, value)
}

func (tr *State) setNilRoot(enc *figbuf.Encoder, path []uint8, value []byte) ([]byte, error) {
	if len(value) == 0 {
		return nil, nil
	}
	return tr.setNode(enc, tr.getNewNode(compactEncode(path, true), value))
}

func (tr *State) setNilNode(enc *figbuf.Encoder, path []uint8, value []byte) ([]byte, error) {
	return tr.setNode(enc, tr.getNewNode(compactEncode(path, true), value))
}

func (tr *State) setBranchNode(enc *figbuf.Encoder, dec *figbuf.Decoder, node [][]byte, path []uint8, value []byte) ([]byte, error) {
	/*
		<branch
		<1234> -> value
		....
		<branch [1]=recurse(234, value)
		------------------------------
		<branch
		<1234> -> value
		...
		<branch [1]=nil [?]=hash
		...
		<?> : hash
		------------------------------
		<branch
		<1234> -> value
		...
		<branch [1]=nil [16]=value
		...
		<16> : value
	*/
	if len(path) == 0 {
		node[16] = value
	} else {
		k, err := tr.set(enc, dec, node[path[0]], path[1:], value)
		if err != nil {
			return nil, err
		}
		node[path[0]] = k
	}
	if i, ok := singleNode(node); ok {
		return tr.setSingleBranchNode(enc, dec, node, nil, i)
	}
	return tr.setNode(enc, node)
}

func (tr *State) setSingleBranchNode(enc *figbuf.Encoder, dec *figbuf.Decoder, node [][]byte, path []uint8, i uint8) ([]byte, error) {
	if i == 16 {
		return tr.setNode(enc, tr.getNewNode(compactEncode(path, true), node[i]))
	}
	child, err := tr.getNode(dec, node[i])
	if err != nil {
		return nil, err
	}
	var newpath []uint8
	if len(child) == 16 {
		newpath = make([]uint8, 0, len(path)+1)
		newpath = append(newpath, path...)
		newpath = append(newpath, i)
		return tr.setNode(enc, tr.getNewNode(compactEncode(newpath, false), node[i]))
	}
	short, term, err := compactDecode(child[0])
	if err != nil {
		return nil, err
	}
	newpath = make([]uint8, 0, len(path)+1+len(short))
	newpath = append(newpath, path...)
	newpath = append(newpath, i)
	newpath = append(newpath, short...)
	return tr.setNode(enc, tr.getNewNode(compactEncode(newpath, term), child[1]))
}

func (tr *State) setLeafOrExtension(enc *figbuf.Encoder, dec *figbuf.Decoder, node [][]byte, path []uint8, value []byte) ([]byte, error) {
	short, term, err := compactDecode(node[0])
	if err != nil {
		return nil, err
	}
	if term {
		return tr.setLeaf(enc, dec, node, short, path, value)
	}
	return tr.setExtension(enc, dec, node, short, path, value)
}

func (tr *State) setLeaf(enc *figbuf.Encoder, dec *figbuf.Decoder, node [][]byte, short []uint8, path []uint8, value []byte) ([]byte, error) {
	o, rs, rp := overlap(short, path)
	if len(rs) == 0 && len(rp) == 0 {
		return tr.setLeafIsPath(enc, node, value)
	}
	if len(rs) == 0 {
		return tr.setLeafNoShort(enc, dec, node, o, rp, value)
	}
	if len(rp) == 0 {
		return tr.setLeafNoPath(enc, dec, node, o, rs, value)
	}
	return tr.setLeafShortAndPath(enc, dec, node, o, rs, rp, value)
}

func (tr *State) setLeafIsPath(enc *figbuf.Encoder, node [][]byte, value []byte) ([]byte, error) {
	/*
		<1234> : <value>
		<1234> -> value
	*/
	node[1] = value
	return tr.setNode(enc, node)
}

func (tr *State) setLeafNoShort(enc *figbuf.Encoder, dec *figbuf.Decoder, node [][]byte, overlap, path []uint8, value []byte) ([]byte, error) {
	/*
		<1234> : <value>
		<123456> -> value
		...
			   <1234> : hashA
		hashA: <branch [5]=>hashB, [16]=<value>
		hashB: <6> : value
		------------------------------
		<1234> : <value>
		<12345> -> value
		...
			   <1234> : hashA
		hashA: <branch [5]=>hashB, [16]=<value>
		hashB: <> : value
	*/
	branch := tr.getNewBranch()
	branch[16] = node[1]
	k, err := tr.setNode(enc, tr.getNewNode(compactEncode(path[1:], true), value))
	if err != nil {
		return nil, err
	}
	branch[path[0]] = k
	if i, ok := singleNode(branch); ok {
		return tr.setSingleBranchNode(enc, dec, branch, overlap, i)
	}
	k, err = tr.setNode(enc, branch)
	if err != nil {
		return nil, err
	}
	if len(overlap) == 0 {
		return k, nil
	}
	return tr.setNode(enc, tr.getNewNode(compactEncode(overlap, false), k))
}

func (tr *State) setLeafNoPath(enc *figbuf.Encoder, dec *figbuf.Decoder, node [][]byte, overlap, path []uint8, value []byte) ([]byte, error) {
	/*
		<123456> : <value>
		<1234> -> value
		...
			   <1234> : hashA
		hashA: <branch [5]=>hashB, [16]=value
		hashB: <6> : <value>
		------------------------------
		<12345> : <value>
		<1234> -> value
		...
			   <1234> : hashA
		hashA: <branch [5]=>hashB, [16]=value
		hashB: <> : <value>
	*/
	branch := tr.getNewBranch()
	branch[16] = value
	k, err := tr.setNode(enc, tr.getNewNode(compactEncode(path[1:], true), node[1]))
	if err != nil {
		return nil, err
	}
	branch[path[0]] = k
	if i, ok := singleNode(branch); ok {
		return tr.setSingleBranchNode(enc, dec, branch, overlap, i)
	}
	k, err = tr.setNode(enc, branch)
	if err != nil {
		return nil, err
	}
	if len(overlap) == 0 {
		return k, nil
	}
	return tr.setNode(enc, tr.getNewNode(compactEncode(overlap, false), k))
}

func (tr *State) setLeafShortAndPath(enc *figbuf.Encoder, dec *figbuf.Decoder, node [][]byte, overlap, short, path []uint8, value []byte) ([]byte, error) {
	/*
		<123456> : <value>
		<123478> -> value
		...
			   <1234> : hashA
		hashA: <branch [5]=>hashB, [7]=hashC
		hashB: <6> : <value>
		hashC: <8> : value
		------------------------------
		<12345> : <value>
		<12347> -> value
		...
			   <1234> : hashA
		hashA: <branch [5]=>hashB, [7]=hashC
		hashB: <> : <value>
		hashC: <> : value
	*/
	branch := tr.getNewBranch()
	k, err := tr.setNode(enc, tr.getNewNode(compactEncode(short[1:], true), node[1]))
	if err != nil {
		return nil, err
	}
	branch[short[0]] = k
	k, err = tr.setNode(enc, tr.getNewNode(compactEncode(path[1:], true), value))
	if err != nil {
		return nil, err
	}
	branch[path[0]] = k
	if i, ok := singleNode(branch); ok {
		return tr.setSingleBranchNode(enc, dec, branch, overlap, i)
	}
	k, err = tr.setNode(enc, branch)
	if err != nil {
		return nil, err
	}
	if len(overlap) == 0 {
		return k, nil
	}
	return tr.setNode(enc, tr.getNewNode(compactEncode(overlap, false), k))
}

func (tr *State) setExtension(enc *figbuf.Encoder, dec *figbuf.Decoder, node [][]byte, short []uint8, path []uint8, value []byte) ([]byte, error) {
	o, rs, rp := overlap(short, path)
	if len(rs) == 0 && len(rp) == 0 {
		return tr.setExtensionIsPath(enc, dec, node, path, value)
	}
	if len(rs) == 0 {
		return tr.setExtensionNoShort(enc, dec, node, o, rp, value)
	}
	if len(rp) == 0 {
		return tr.setExtensionNoPath(enc, dec, node, o, rs, value)
	}
	return tr.setExtensionShortAndPath(enc, dec, node, o, rs, rp, value)
}

func (tr *State) setExtensionIsPath(enc *figbuf.Encoder, dec *figbuf.Decoder, node [][]byte, path []uint8, value []byte) ([]byte, error) {
	/*
		<1234> : hashA
		<1234> -> value
		...
			   <1234> : hashA
		hashA: <branch ... [16]=value
		------------------------------
		<1234> : hashA
		<1234> -> value
		...
			   <1234> : hashA
		hashA: <branch [?]=hashB [16]=nil
		...
			   <1234?> : hashB
		------------------------------
		<1234> : hashA
		<1234> -> value
		...
			   <1234> : hashA
		hashA: <branch ...nil [16]=value
		...
			IMPOSSIBLE (set 16 on an existing branch, so must be at least 1 non-nil key)

	*/
	branch, err := tr.getNode(dec, node[1])
	if err != nil {
		return nil, err
	}
	branch[16] = value
	if i, ok := singleNode(branch); ok {
		return tr.setSingleBranchNode(enc, dec, branch, path, i)
	}
	k, err := tr.setNode(enc, branch)
	if err != nil {
		return nil, err
	}
	node[1] = k
	return tr.setNode(enc, node)
}

func (tr *State) setExtensionNoShort(enc *figbuf.Encoder, dec *figbuf.Decoder, node [][]byte, overlap, path []uint8, value []byte) ([]byte, error) {
	/*
		<1234> : hashA
		<123456> -> value
		...
			   <1234> : hashA
		hashA: <branch [5]=recurse(6, value)
		------------------------------
		<1234> : hashA
		<123456> -> value
		...
			   <1234> : hashA
		hashA: <branch [5]=nil [?]=?
		...
			   <1234?> : hashB
		------------------------------
		<1234> : hashA
		<1234> -> value
		...
			   <1234> : hashA
		hashA: <branch ...nil [16]=value
		...
			   <1234?> : value

	*/
	branch, err := tr.getNode(dec, node[1])
	if err != nil {
		return nil, err
	}
	k, err := tr.set(enc, dec, branch[path[0]], path[1:], value)
	if err != nil {
		return nil, err
	}
	branch[path[0]] = k
	if i, ok := singleNode(branch); ok {
		return tr.setSingleBranchNode(enc, dec, branch, overlap, i)
	}
	k, err = tr.setNode(enc, branch)
	if err != nil {
		return nil, err
	}
	node[1] = k
	return tr.setNode(enc, node)
}

func (tr *State) setExtensionNoPath(enc *figbuf.Encoder, dec *figbuf.Decoder, node [][]byte, overlap, short []uint8, value []byte) ([]byte, error) {
	/*
		<123456> : hashA
		<1234> -> value
		...
			   <1234> : hashB
		hashB: <branch [5]=hashC [16]=value
		hashC: <6> : hashA
		------------------------------
		<12345> : hashA
		<1234> -> value
		...
			   <1234> : hashB
		hashB: <branch [5]=hashA [16]=value
		------------------------------
		<123456> : hashA
		<1234> -> value
		...
			   <1234> : hashB
		hashB: <branch [?]=hashB [16]=nil
		...
			   <1234?> : hashB
		------------------------------
		<123456> : hashA
		<1234> -> value
		...
			   <1234> : hashA
		hashB: <branch ...nil [16]=value
		...
			IMPOSSIBLE (set 16 on an existing branch, so must be at least 1 non-nil key)

	*/
	branch := tr.getNewBranch()
	var k []byte
	var err error
	if len(short) == 1 {
		branch[short[0]] = node[1]
	} else {
		k, err = tr.setNode(enc, tr.getNewNode(compactEncode(short[1:], false), node[1]))
		if err != nil {
			return nil, err
		}
		branch[short[0]] = k
	}
	branch[16] = value
	if i, ok := singleNode(branch); ok {
		return tr.setSingleBranchNode(enc, dec, branch, overlap, i)
	}
	k, err = tr.setNode(enc, branch)
	if err != nil {
		return nil, err
	}
	if len(overlap) == 0 {
		return k, nil
	}
	return tr.setNode(enc, tr.getNewNode(compactEncode(overlap, false), k))
}

func (tr *State) setExtensionShortAndPath(enc *figbuf.Encoder, dec *figbuf.Decoder, node [][]byte, overlap, short, path []uint8, value []byte) ([]byte, error) {
	/*
		<123456> : hashA
		<123478> -> value
		...
			   <1234> : hashB
		hashB: <branch [5]=>hashC, [7]=hashD
		hashC: <6> : hashA
		hashD: <8> : value
		------------------------------
		<12345> : hashA
		<12347> -> value
		...
			   <1234> : hashB
		hashB: <branch [5]=>hashA, [7]=hashD
		hashD: <> : value
	*/
	branch := tr.getNewBranch()
	var k []byte
	var err error
	if len(short) == 1 {
		branch[short[0]] = node[1]
	} else {
		k, err = tr.setNode(enc, tr.getNewNode(compactEncode(short[1:], false), node[1]))
		if err != nil {
			return nil, err
		}
		branch[short[0]] = k
	}
	k, err = tr.setNode(enc, tr.getNewNode(compactEncode(path[1:], true), value))
	if err != nil {
		return nil, err
	}
	branch[path[0]] = k
	if i, ok := singleNode(branch); ok {
		return tr.setSingleBranchNode(enc, dec, branch, overlap, i)
	}
	k, err = tr.setNode(enc, branch)
	if err != nil {
		return nil, err
	}
	if len(overlap) == 0 {
		return k, nil
	}
	return tr.setNode(enc, tr.getNewNode(compactEncode(overlap, false), k))
}

func (tr *State) setNode(enc *figbuf.Encoder, node [][]byte) ([]byte, error) {
	if nullNode(node) {
		return nil, nil
	}
	v := enc.EncodeBytesSlice(node)
	if len(v) < 32 {
		return enc.Copy(v), nil
	}
	k := hasher.Hash256(node...)
	if tr.Cache != nil {
		tr.Cache.Add(k, v)
	}
	err := tr.KeyStore.Set(k, v)
	if err != nil {
		return nil, err
	}
	return k, nil
}

// Get returns the value stored at a key under a given Merkle root
func (tr *State) Get(root, key []byte) ([]byte, error) {
	tr.lock.RLock()
	defer tr.lock.RUnlock()

	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	if len(root) == 0 {
		return nil, nil
	}
	if len(key) == 0 {
		return nil, nil
	}

	path := nibbles(key)
	return tr.get(dec, root, path)
}

func (tr *State) get(dec *figbuf.Decoder, root []byte, path []uint8) ([]byte, error) {
	node, err := tr.getNode(dec, root)
	if err != nil {
		return nil, err
	}
	if len(node) == 0 {
		return nil, nil
	}
	if len(node) == 17 {
		if len(path) == 0 {
			return node[16], nil
		}
		return tr.get(dec, node[path[0]], path[1:])
	}
	short, term, err := compactDecode(node[0])
	if err != nil {
		return nil, err
	}
	if term {
		if pathEqual(short, path) {
			return node[1], nil
		}
		return nil, nil
	}
	_, rs, rp := overlap(short, path)
	if len(rs) != 0 {
		return nil, nil
	}
	return tr.get(dec, node[1], rp)
}

// GetAndProve returns the value stored at a key under a given Merkle root,
// along with a Merkle proof that the value resides at key under root
func (tr *State) GetAndProve(root, key []byte) ([]byte, [][][]byte, error) {
	tr.lock.RLock()
	defer tr.lock.RUnlock()

	dec := figbuf.DecoderPool.Get().(*figbuf.Decoder)
	defer figbuf.DecoderPool.Put(dec)

	if len(root) == 0 {
		return nil, nil, nil
	}
	if len(key) == 0 {
		return nil, nil, nil
	}

	path := nibbles(key)
	return tr.getAndProve(dec, root, path)
}

func (tr *State) getAndProve(dec *figbuf.Decoder, root []byte, path []uint8) ([]byte, [][][]byte, error) {
	var proof [][][]byte
	node, err := tr.getNode(dec, root)
	if err != nil {
		return nil, nil, err
	}
	if len(node) == 0 {
		return nil, nil, nil
	}
	if len(node) == 17 {
		if len(path) == 0 {
			proof = append(proof, node)
			return node[16], proof, nil
		}
		r, p, err := tr.getAndProve(dec, node[path[0]], path[1:])
		if err != nil {
			return nil, nil, err
		}
		proof = append(proof, node)
		proof = append(proof, p...)
		return r, proof, nil
	}
	short, term, err := compactDecode(node[0])
	if err != nil {
		return nil, nil, err
	}
	if term {
		if pathEqual(short, path) {
			proof = append(proof, node)
			return node[1], proof, nil
		}
		return nil, nil, nil
	}
	_, rs, rp := overlap(short, path)
	if len(rs) != 0 {
		return nil, nil, nil
	}
	r, p, err := tr.getAndProve(dec, node[1], rp)
	if err != nil {
		return nil, nil, err
	}
	proof = append(proof, node)
	proof = append(proof, p...)
	return r, proof, nil
}

func (tr *State) getNode(dec *figbuf.Decoder, k []byte) ([][]byte, error) {
	var v []byte
	var err error
	if len(k) == 0 {
		return nil, nil
	}
	if len(k) < 32 {
		v = dec.Copy(k)
	} else if tr.Cache != nil {
		if c, ok := tr.Cache.Get(k); ok {
			v = c
		}
	}
	if v == nil {
		v, err = tr.KeyStore.Get(k)

		if err != nil {
			return nil, err
		}
	}
	var node [][]byte
	node, err = dec.DecodeBytesSlice(v)
	return node, err
}

// ValidateMPT confirms whether a merkle patricia trie proof is valid for a given root, key, and value
func ValidateMPT(root []byte, key, value []byte, proof [][][]byte) bool {
	enc := figbuf.EncoderPool.Get().(*figbuf.Encoder)
	defer figbuf.EncoderPool.Put(enc)

	path := nibbles(key)

	nh, ok := validate(enc, path, value, proof)
	if ok && bytes.Equal(root, nh) {
		return true
	}
	return false
}

func validate(enc *figbuf.Encoder, path []uint8, value []byte, proof [][][]byte) ([]byte, bool) {
	if len(proof) == 0 {
		return nil, false
	}
	node := proof[0]
	if len(node) == 17 {
		if len(path) == 0 {
			if bytes.Equal(value, node[16]) {
				return hashNode(enc, node), true
			}
			return nil, false
		}
		nh, ok := validate(enc, path[1:], value, proof[1:])
		if ok && bytes.Equal(nh, node[path[0]]) {
			return hashNode(enc, node), true
		}
		return nil, false
	}
	short, term, err := compactDecode(node[0])
	if err != nil {
		return nil, false
	}
	if term {
		if pathEqual(short, path) {
			if bytes.Equal(value, node[1]) {
				return hashNode(enc, node), true
			}
		}
		return nil, false
	}
	_, rs, rp := overlap(short, path)
	if len(rs) != 0 {
		return nil, false
	}
	nh, ok := validate(enc, rp, value, proof[1:])
	if ok && bytes.Equal(nh, node[1]) {
		return hashNode(enc, node), true
	}
	return nil, false
}

func hashNode(enc *figbuf.Encoder, node [][]byte) []byte {
	if nullNode(node) {
		return nil
	}
	v := enc.EncodeBytesSlice(node)
	if len(v) < 32 {
		return enc.Copy(v)
	}
	return hasher.Hash256(node...)
}

// Helper Functions

func compactEncode(path []uint8, term bool) []byte {
	var termSet int
	if term {
		termSet = 1
	}
	flags := uint8(2*termSet + len(path)&1)
	var fpath []uint8
	if len(path)&1 == 1 {
		fpath = make([]uint8, 0, len(path)+1)
		fpath = append(fpath, flags)
		fpath = append(fpath, path...)
	} else {
		fpath = make([]uint8, 0, len(path)+2)
		fpath = append(fpath, flags)
		fpath = append(fpath, 0)
		fpath = append(fpath, path...)
	}
	return nibbleBytes(fpath)
}

func compactDecode(bytes []byte) ([]uint8, bool, error) {
	if len(bytes) == 0 {
		return nil, false, ErrInvalidCompactEncoding
	}
	nibs := nibbles(bytes)
	flags := nibs[0]
	var short []uint8
	var term bool
	switch flags {
	case 0:
		short, term = nibs[2:], false
	case 1:
		short, term = nibs[1:], false
	case 2:
		short, term = nibs[2:], true
	case 3:
		short, term = nibs[1:], true
	default:
		return nil, false, ErrInvalidCompactEncoding
	}
	return short, term, nil
}

func nullNode(bb [][]byte) bool {
	if len(bb) == 2 {
		return len(bb[1]) == 0
	}
	for _, b := range bb {
		if len(b) > 0 {
			return false
		}
	}
	return true
}

func singleNode(bb [][]byte) (uint8, bool) {
	var count int
	var hit uint8
	for i, b := range bb {
		if len(b) > 0 {
			hit = uint8(i)
			count++
		}
		if count > 1 {
			return 0, false
		}
	}
	return hit, true
}

func overlap(short, path []uint8) (overlap []uint8, sremainder []uint8, premainder []uint8) {
	for i, v := range short {
		if i > len(path)-1 || v != path[i] {
			break
		}
		overlap = append(overlap, v)
	}
	sremainder = make([]byte, len(short)-len(overlap))
	copy(sremainder, short[len(overlap):])
	premainder = make([]byte, len(path)-len(overlap))
	copy(premainder, path[len(overlap):])
	return
}

func pathEqual(a, b []uint8) bool {
	for i, v := range a {
		if i > len(b)-1 {
			return false
		}
		if v != b[i] {
			return false
		}
	}
	return true
}

func nibbles(bytes []byte) []uint8 {
	nibbles := make([]uint8, len(bytes)*2)
	for i, b := range bytes {
		nibbles[i*2] = uint8(b >> 4)
		nibbles[i*2+1] = uint8(b & 0xf)
	}
	return nibbles
}

func nibbleBytes(nibbles []uint8) []byte {
	var nbytes []byte
	if len(nibbles)&1 == 1 {
		nbytes = make([]byte, (len(nibbles)+1)/2)
		nbytes[0] = byte(0) + byte(nibbles[1]&0xf)
		for i := 2; i < len(nibbles); i += 2 {
			nbytes[i/2] = byte(nibbles[i]<<4) + byte(nibbles[i+1]&0xf)
		}
	} else {
		nbytes = make([]byte, len(nibbles)/2)
		for i := 0; i < len(nibbles); i += 2 {
			nbytes[i/2] = byte(nibbles[i]<<4) + byte(nibbles[i+1]&0xf)
		}
	}
	return nbytes
}
