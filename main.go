package main

import (
	"fmt"

	"github.com/figaro-tech/go-figaro/core"
)

func main() {
	chain := core.NewBlockChain()
	logBlock(chain.Head())
	for i := 1; i <= 20; i++ {
		h := chain.Head()
		d := fmt.Sprintf("This is block number %d", h.Index+1)
		n := chain.CreateBlock(d)
		logBlock(n)
	}
}

func logBlock(b *core.Block) {
	fmt.Println("\n=================================================")
	fmt.Println(fmt.Sprintf("\nIndex: %d", b.Index))
	fmt.Println(fmt.Sprintf("TimeStamp: %s", b.TimeStamp))
	fmt.Println(fmt.Sprintf("Data: %s", b.Data))
	fmt.Println(fmt.Sprintf("PreviousHash: %s", b.PreviousHash))
	fmt.Println(fmt.Sprintf("Hash: %s", b.Hash))
}
