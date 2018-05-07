package main

import (
	"context"
	"log"
	"time"

	"github.com/figaro-tech/go-figaro/figp2p"
	"github.com/libp2p/go-floodsub"
)

var rendezvousName = "figaro-demo-v0.0.1"
var topicName = "figaro-message"

func main() {
	ctx := context.Background()

	// Create Boot Node
	bootNode := newBootNode(ctx)
	go bootNode.Listen(ctx)

	// Create Other Nodes
	nodes := []*figp2p.Node{}
	for i := 0; i < 5; i++ {
		node := newNode(ctx, bootNode)
		nodes = append(nodes, node)
		go node.Listen(ctx)
	}

	// Have them talk
	time.Sleep(1 * time.Second)
	bootNode.Broadcast(ctx, topicName, []byte("Hey everyone!"))
	time.Sleep(1 * time.Second)
	bootNode.Send(nodes[2].PeerID(), []byte("Yo, sup??"))
	time.Sleep(3 * time.Second)
}

func newBootNode(ctx context.Context) *figp2p.Node {
	bootNode, err := figp2p.NewBootstrapNode(ctx, demoMuxer())
	if err != nil {
		log.Panic(err)
	}
	return bootNode
}

func newNode(ctx context.Context, bootNode *figp2p.Node) *figp2p.Node {
	newNode, err := figp2p.NewNode(ctx, bootNode.FullAddresses(), rendezvousName, demoMuxer())
	if err != nil {
		log.Panic(err)
	}
	return newNode
}

func demoMuxer() *figp2p.FanoutMux {
	fanoutMux := figp2p.NewFanoutMux()
	fanoutMux.Handle(topicName, func(node *figp2p.Node, msg *floodsub.Message) {
		log.Println(node.ID(), "Got a message: ", string(msg.GetData()))
	})
	fanoutMux.HandleDirectMessage(func(node *figp2p.Node, msg *floodsub.Message) {
		log.Println(node.ID(), "Got a direct message: ", string(msg.GetData()))
	})
	return fanoutMux
}
