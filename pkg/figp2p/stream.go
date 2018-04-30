package p2p

import (
	"bufio"
	"fmt"

	"github.com/libp2p/go-libp2p-net"
)

func makeStreamHandler(node *Node) net.StreamHandler {
	return func(s net.Stream) {
		remotePeerID := s.Conn().RemotePeer().Pretty()
		senderChannel := make(chan string)
		receiverChannel := node.receiverChannel
		node.senderChannels[remotePeerID] = senderChannel

		fmt.Println(node.PeerID(), "got a new stream from ", remotePeerID)
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		go readData(rw, receiverChannel)
		go writeData(rw, senderChannel)
	}
}

func readData(rw *bufio.ReadWriter, receiverChannel chan string) {
	for {
		str, _ := rw.ReadString('\n')
		if str == "" {
			return
		}
		if str != "\n" {
			receiverChannel <- str
		}
	}
}

func writeData(rw *bufio.ReadWriter, senderChannel chan string) {
	for {
		message := <-senderChannel
		rw.WriteString(fmt.Sprintf("%s\n", message))
		rw.Flush()
	}
}
