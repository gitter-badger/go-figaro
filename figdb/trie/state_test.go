package trie_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/figaro-tech/go-figaro/figdb/mock"

	"github.com/figaro-tech/go-figaro/figdb/trie"
)

func ExampleState_Set() {
	state := trie.State{
		KeyStore: mock.NewKeyStore(),
	}
	do := [][]byte{[]byte("do"), []byte("verb")}
	dog := [][]byte{[]byte("dog"), []byte("puppy")}
	doge := [][]byte{[]byte("doge"), []byte("coin")}
	horse := [][]byte{[]byte("horse"), []byte("stallion")}

	root, err := state.Set(nil, do[0], do[1])
	if err != nil {
		log.Print(err)
		return
	}

	fmt.Printf("%#x\n", root)

	root, err = state.Set(root, dog[0], dog[1])
	if err != nil {
		log.Print(err)
		return
	}

	fmt.Printf("%#x\n", root)

	root, err = state.Set(root, dog[0], nil)
	if err != nil {
		log.Print(err)
		return
	}

	fmt.Printf("%#x\n", root)

	root, err = state.Set(root, dog[0], dog[1])
	if err != nil {
		log.Print(err)
		return
	}

	fmt.Printf("%#x\n", root)

	root, err = state.Set(root, doge[0], doge[1])
	if err != nil {
		log.Print(err)
		return
	}

	fmt.Printf("%#x\n", root)

	root, err = state.Set(root, horse[0], horse[1])
	if err != nil {
		log.Print(err)
		return
	}

	fmt.Printf("%#x\n", root)

	root, err = state.Set(root, horse[0], nil)
	if err != nil {
		log.Print(err)
		return
	}

	fmt.Printf("%#x\n", root)

	root, err = state.Set(root, horse[0], horse[1])
	if err != nil {
		log.Print(err)
		return
	}

	fmt.Printf("%#x\n", root)

	// Output:
	// 0xc98320646f8476657262
	// 0x30a2e29e52a440b4184d5118ef67207bc4408277c03e40a244f02ce9125b79d7
	// 0xc98320646f8476657262
	// 0x30a2e29e52a440b4184d5118ef67207bc4408277c03e40a244f02ce9125b79d7
	// 0x0897dc12f0a9fdfdfdc58aebf798d303454917409816bb0361a4264359c58168
	// 0x5888437106063661c31dda57f26bf76c63ee1682a1421cee7fe45e8809d67c9a
	// 0x0897dc12f0a9fdfdfdc58aebf798d303454917409816bb0361a4264359c58168
	// 0x5888437106063661c31dda57f26bf76c63ee1682a1421cee7fe45e8809d67c9a
}

func BenchmarkState_Set(b *testing.B) {
	do := [][]byte{[]byte("do"), []byte("verb")}
	dog := [][]byte{[]byte("dog"), []byte("puppy")}
	doge := [][]byte{[]byte("doge"), []byte("coin")}
	horse := [][]byte{[]byte("horse"), []byte("stallion")}

	state := trie.State{
		KeyStore: mock.NewKeyStore(),
	}
	root, err := state.Set(nil, do[0], do[1])
	if err != nil {
		log.Fatal(err)
		return
	}
	root, err = state.Set(root, dog[0], dog[1])
	if err != nil {
		log.Fatal(err)
		return
	}
	root, err = state.Set(root, doge[0], doge[1])
	if err != nil {
		log.Fatal(err)
		return
	}
	for i := 0; i < b.N; i++ {
		state.Set(root, horse[0], horse[1])
	}
}
