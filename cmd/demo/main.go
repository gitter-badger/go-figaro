package main

// type encdec struct{}

// func (e *encdec) Encode(i interface{}) ([]byte, error) {
// 	return rlp_old.Encode(i)
// }
// func (e *encdec) Decode(i interface{}, b []byte) error {
// 	return rlp_old.Decode(i, b)
// }

func main() {
	// demoArchiveTrie()
	// demoStateTrie()
	// h := &sha256.Hasher{}
	// e := &rlp.EncoderDecoder{}
	// db := figdb.New("/tmp/figaro", h, e)

}

// func demoArchiveTrie() {
// 	store := memory.NewStore()
// 	cypher := &crypto.Cypher{}
// 	encoder := &encdec{}
// 	tr := trie.NewArchiveTrie(store, cypher, encoder)
// 	tp := trie.NewArchiveProver(cypher)

// 	data := [][]byte{
// 		[]byte("dog"),
// 		[]byte("cat"),
// 		[]byte("bird"),
// 		[]byte("wallabeye"),
// 	}
// 	fmt.Println("\nArchiving data...")
// 	root, err := tr.Archive(data)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Printf("Root is %#x\n", root)

// 	fmt.Println("\nRetrieving value...")
// 	value := tr.Retrieve(root, 2)
// 	fmt.Printf("Value is % s\n", string(value))

// 	fmt.Println("\nProving value...")
// 	proof, err := tr.Prove(root, 2, []byte("bird"))
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	result := tp.VerifyProof(root, 2, []byte("bird"), proof)
// 	fmt.Printf("Proof is %t\n", result)
// }

// func demoStateTrie() {
// 	store := memory.NewStore()
// 	cypher := &crypto.Cypher{}
// 	encoder := &encdec{}
// 	tr := trie.NewStateTrie(store, cypher, encoder)
// 	tp := trie.NewStateProver(cypher)

// 	fmt.Println("\nInitialize with a bunch of values...")
// 	root := tr.Update(nil, []byte("one"), []byte("apple"))
// 	root = tr.Update(nil, []byte("two"), []byte("banana"))
// 	root = tr.Update(nil, []byte("three"), []byte("carrot"))
// 	root = tr.Update(nil, []byte("four"), []byte("drumstick"))
// 	root = tr.Update(nil, []byte("five"), []byte("eggplant"))
// 	root = tr.Update(nil, []byte("six"), []byte("fig"))
// 	root = tr.Update(nil, []byte("seven"), []byte("ginger"))
// 	root = tr.Update(nil, []byte("eight"), []byte("horticorts"))
// 	root = tr.Update(nil, []byte("nine"), []byte("icecream"))
// 	root = tr.Update(nil, []byte("ten"), []byte("jalepeno"))

// 	fmt.Println("\nSetting value of 'dog'...")
// 	root = tr.Update(root, []byte("dog"), []byte("puppy"))
// 	fmt.Printf("Root is now: %#x\n", root)
// 	value := tr.Get(root, []byte("dog"))
// 	fmt.Printf("Dog is now: %s\n", string(value))
// 	value = tr.Get(root, []byte("doge"))
// 	fmt.Printf("Doge is now: %s\n", string(value))

// 	fmt.Println("\nSetting value of 'doge'...")
// 	root = tr.Update(root, []byte("doge"), []byte("coin"))
// 	fmt.Printf("Root is now: %#x\n", root)
// 	value = tr.Get(root, []byte("dog"))
// 	fmt.Printf("Dog is now: %s\n", string(value))
// 	value = tr.Get(root, []byte("doge"))
// 	fmt.Printf("Doge is now: %s\n", string(value))

// 	checkpoint := make([]byte, len(root))
// 	copy(checkpoint, root)

// 	fmt.Println("\nProving value of 'dog'...")
// 	proof, err := tr.Prove(root, []byte("dog"), []byte("puppy"))
// 	fmt.Printf("Size proof is %d, size store is %d\n", len(proof), store.Len())
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	result := tp.VerifyProof(root, []byte("dog"), []byte("puppy"), proof)
// 	fmt.Printf("Proof is %t\n", result)
// 	value = tr.Get(root, []byte("dog"))
// 	fmt.Printf("Dog is now: %s\n", string(value))
// 	value = tr.Get(root, []byte("doge"))
// 	fmt.Printf("Doge is now: %s\n", string(value))

// 	fmt.Println("\nProving value of 'doge'...")
// 	proof, err = tr.Prove(root, []byte("doge"), []byte("coin"))
// 	fmt.Printf("Size proof is %d, size store is %d\n", len(proof), store.Len())
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	result = tp.VerifyProof(root, []byte("doge"), []byte("coin"), proof)
// 	fmt.Printf("Proof is %t\n", result)
// 	value = tr.Get(root, []byte("dog"))
// 	fmt.Printf("Dog is now: %s\n", string(value))
// 	value = tr.Get(root, []byte("doge"))
// 	fmt.Printf("Doge is now: %s\n", string(value))

// 	fmt.Println("\nDeleting value of 'doge'...")
// 	root = tr.Delete(root, []byte("doge"))
// 	fmt.Printf("Root is now: %#x\n", root)
// 	value = tr.Get(root, []byte("dog"))
// 	fmt.Printf("Dog is now: %s\n", string(value))
// 	value = tr.Get(root, []byte("doge"))
// 	fmt.Printf("Doge is now: %s\n", string(value))

// 	fmt.Println("\nDeleting value of 'dog'...")
// 	root = tr.Delete(root, []byte("dog"))
// 	fmt.Printf("Root is now: %#x\n", root)
// 	value = tr.Get(root, []byte("dog"))
// 	fmt.Printf("Dog is now: %s\n", string(value))
// 	value = tr.Get(root, []byte("doge"))
// 	fmt.Printf("Doge is now: %s\n", string(value))

// 	fmt.Println("\nEvaluating checkpoint...")
// 	value = tr.Get(checkpoint, []byte("dog"))
// 	fmt.Printf("Checkpoint: Dog is now: %s\n", string(value))
// 	value = tr.Get(checkpoint, []byte("doge"))
// 	fmt.Printf("Checkpoint: Doge is now: %s\n", string(value))
// }

// func logBlock(b *core.Block) {
// 	fmt.Println("\n=================================================")
// 	fmt.Println(fmt.Sprintf("\nIndex: %d", b.Index))
// 	fmt.Println(fmt.Sprintf("TimeStamp: %s", b.TimeStamp))
// 	fmt.Println(fmt.Sprintf("Data: %s", b.Data))
// 	fmt.Println(fmt.Sprintf("PreviousHash: %s", b.PreviousHash))
// 	fmt.Println(fmt.Sprintf("Hash: %s", b.Hash))
// }
