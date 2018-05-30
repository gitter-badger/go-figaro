package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/figaro-tech/go-figaro/figp2p"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	ctx := context.Background()

	bootAddFlag := flag.String("bootaddr", "", "Bootstrap Node Addrs")
	flag.Parse()

	bootAddr, err := multiaddr.NewMultiaddr(*bootAddFlag)
	if err != nil {
		log.Panic(err)
	}

	node, err := figp2p.NewNode(ctx, []multiaddr.Multiaddr{bootAddr})
	if err != nil {
		log.Panic(err)
	}

	go node.Start(ctx)

	for {
		time.Sleep(1 * time.Second)
		log.Println("Connections:", len(node.Host().Network().Conns()))
	}
}
