package main

import (
	"fmt"

	"github.com/figaro-tech/go-figaro/core"
	"github.com/figaro-tech/go-figaro/encoding/rlp"
)

func main() {
	// chain := core.NewBlockChain()
	// logBlock(chain.Head())
	// for i := 1; i <= 20; i++ {
	// 	h := chain.Head()
	// 	d := fmt.Sprintf("This is block number %d", h.Index+1)
	// 	n := chain.CreateBlock(d)
	// 	logBlock(n)
	// }

	// t1, _ := rlp.Encode("hello my name is inigo montoya")
	// fmt.Printf("% x\n", t1)
	// var t2 string
	// err := rlp.Decode(&t2, t1)
	// fmt.Printf("% s\n", t2)
	// fmt.Println(err)

	// one := "hello my name is"
	// two := "inigo montoya"
	// t1, _ := rlp.Encode([]interface{}{&one, &two})
	// fmt.Printf("% x\n", t1)
	// var t2 []string
	// err := rlp.Decode(&t2, t1)
	// fmt.Printf("%#v\n", t2)
	// fmt.Println(err)

	// one := 1
	// two := 2
	// t1, _ := rlp.Encode(map[string]*int{"one": &one, "two": &two})
	// fmt.Printf("% x\n", t1)
	// var t2 map[string]*int
	// err := rlp.Decode(&t2, t1)
	// fmt.Printf("%#v\n", t2)
	// fmt.Println(err)

	type test struct {
		One   int
		Two   int
		Three *int `rlp:"-"`
		four  int
	}
	three := 3
	t1, _ := rlp.Encode(test{One: 1, Two: 2, Three: &three, four: 5})
	fmt.Printf("% x\n", t1)
	var t2 test
	err := rlp.Decode(&t2, t1)
	fmt.Printf("%#v\n", t2)
	fmt.Println(err)
}

func logBlock(b *core.Block) {
	fmt.Println("\n=================================================")
	fmt.Println(fmt.Sprintf("\nIndex: %d", b.Index))
	fmt.Println(fmt.Sprintf("TimeStamp: %s", b.TimeStamp))
	fmt.Println(fmt.Sprintf("Data: %s", b.Data))
	fmt.Println(fmt.Sprintf("PreviousHash: %s", b.PreviousHash))
	fmt.Println(fmt.Sprintf("Hash: %s", b.Hash))
}
