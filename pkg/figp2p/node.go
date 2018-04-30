package p2p

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/libp2p/go-libp2p-swarm"
	"github.com/libp2p/go-libp2p/p2p/host/basic"
	"github.com/multiformats/go-multiaddr"
)

const figProtocolDiscriptor = "/figaro/0.0.1-beta1"

type receiverFunc func(string)
type messageChannel chan string
type senderChannelMap map[string]messageChannel

// Node implements figaro.Node
type Node struct {
	address         string
	host            *basichost.BasicHost
	receiverChannel messageChannel
	senderChannels  senderChannelMap
	store           peerstore.Peerstore
}

// AddPeer adds a Peer Node to a Node's Store
func (n *Node) AddPeer(peerAddress string) {
	peerID, targetAddr := peerIDFromAddress(peerAddress)
	n.store.AddAddr(peerID, targetAddr, peerstore.PermanentAddrTTL)

	ctx := context.Background()
	stream, err := n.host.NewStream(ctx, peerID, figProtocolDiscriptor)
	if err != nil {
		panic(err)
	}

	makeStreamHandler(n)(stream)
}

// Address returns the address of the Node
func (n *Node) Address() string {
	return n.address
}

// Broadcast sends a message to all connected Nodes
func (n *Node) Broadcast(message string) {
	for _, nodeID := range n.store.Peers() {
		n.Send(message, nodeID.Pretty())
	}
}

// Listen adds a new receiver to the Node
func (n *Node) Listen(onMessageReceived receiverFunc) {
	go func() {
		for {
			message := <-n.receiverChannel
			onMessageReceived(message)
		}
	}()
}

// PeerID returns the Nodes Peer ID
func (n *Node) PeerID() string {
	return n.host.ID().Pretty()
}

// Send transmits a message to a single peer Node
func (n *Node) Send(message string, nodeID string) {
	if n.senderChannels[nodeID] != nil {
		n.senderChannels[nodeID] <- message
	}
}

// NewNode returns a new Node
func NewNode(port int) *Node {
	nodeID, privKey, pubKey := newNodeID(port)
	networkAddress := newNetworkAddress(port)
	nodeAddress := fmt.Sprintf("%s/ipfs/%s", networkAddress, nodeID.Pretty())

	store := peerstore.NewPeerstore()
	store.AddPrivKey(nodeID, privKey)
	store.AddPubKey(nodeID, pubKey)

	swarm, err := swarm.NewNetwork(
		context.Background(),
		[]multiaddr.Multiaddr{networkAddress},
		nodeID,
		store,
		nil,
	)
	if err != nil {
		panic(err)
	}

	node := &Node{
		address:         nodeAddress,
		host:            basichost.New(swarm),
		receiverChannel: make(chan string),
		senderChannels:  senderChannelMap{},
		store:           store,
	}

	node.host.SetStreamHandler(figProtocolDiscriptor, makeStreamHandler(node))

	return node
}

func newNodeID(seed int) (peer.ID, crypto.PrivKey, crypto.PubKey) {
	r := rand.New(rand.NewSource(int64(seed)))
	privKey, pubKey, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	nodeID, err := peer.IDFromPublicKey(pubKey)
	if err != nil {
		panic(err)
	}

	return nodeID, privKey, pubKey
}

func newNetworkAddress(port int) multiaddr.Multiaddr {
	addressString := fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)
	address, err := multiaddr.NewMultiaddr(addressString)
	if err != nil {
		panic(err)
	}
	return address
}

func peerIDFromAddress(addr string) (peer.ID, multiaddr.Multiaddr) {
	ipfsAddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		log.Fatalln(err)
	}

	pid, err := ipfsAddr.ValueForProtocol(multiaddr.P_IPFS)
	if err != nil {
		log.Fatalln(err)
	}

	peerID, err := peer.IDB58Decode(pid)
	if err != nil {
		log.Fatalln(err)
	}

	ipfsAddress := fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerID))
	targetPeerAddr, _ := multiaddr.NewMultiaddr(ipfsAddress)
	targetAddr := ipfsAddr.Decapsulate(targetPeerAddr)

	return peerID, targetAddr
}
