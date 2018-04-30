package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	p2p "github.com/figaro-tech/figaro/pkg/figp2p"
)

type parameters struct {
	port        int
	peerAddress string
}

func main() {
	params := paramsFromFlags()

	// First Node
	masterNode := p2p.NewNode(params.port)
	logMessages(masterNode)

	// Listener Nodes
	for i := 1; i <= 5; i++ {
		listenerNode := p2p.NewNode(params.port + i)
		listenerNode.AddPeer(masterNode.Address())
		logMessages(listenerNode)
	}

	// Listen for Input
	sendUserInput(masterNode)
}

func paramsFromFlags() parameters {
	port := flag.Int("port", 3000, "Listen Port Number")
	peerAddress := flag.String("conn", "", "Connect directly to a peer")
	flag.Parse()

	return parameters{
		peerAddress: *peerAddress,
		port:        *port,
	}
}

func logMessages(node *p2p.Node) {
	node.Listen(func(message string) {
		fmt.Println(fmt.Sprintf("Node %s received a message: %s", node.PeerID(), string(message)))
	})
}

func sendUserInput(node *p2p.Node) {
	stdReader := bufio.NewReader(os.Stdin)
	for {
		time.Sleep(2 * time.Second)
		fmt.Print(fmt.Sprintf("Send from %s > ", node.PeerID()))
		userInput, err := stdReader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		node.Broadcast(userInput)
	}
}
