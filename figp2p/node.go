// Package figp2p implements a peer-to-peer network
package figp2p

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"

	datastore "github.com/ipfs/go-datastore"
	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	swarm "github.com/libp2p/go-libp2p-swarm"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// Node handles communication with Peers on the network
type Node struct {
	addrs            []multiaddr.Multiaddr
	bootAddrs        []multiaddr.Multiaddr
	connManager      *ConnManager
	dht              *dht.IpfsDHT
	findPeersTrigger chan int
	host             host.Host
}

// NewBootstrapNode creates a new Node that is ready to receive peers
// Unlike Nodes, BootstrapNodes do no work other than routing
func NewBootstrapNode(ctx context.Context) (*Node, error) {
	node, err := newNode(ctx, []multiaddr.Multiaddr{})
	if err != nil {
		return nil, err
	}

	// BootstrapNode DHT applies stream handlers, waiting for clients to connect
	node.dht = dht.NewDHT(ctx, node.host, datastore.NewMapDatastore())

	return node, nil
}

// NewNode creates a new Node that immediately connects to a Bootstrap Node
func NewNode(ctx context.Context, bootAddrs []multiaddr.Multiaddr) (*Node, error) {
	node, err := newNode(ctx, bootAddrs)
	if err != nil {
		return nil, err
	}

	// Node DHTs will open streams with remote hosts once a connection is made
	node.dht = dht.NewDHTClient(ctx, node.host, datastore.NewMapDatastore())

	return node, nil
}

// newNode creates a generalized node for use with Node or BootstrapNode
func newNode(ctx context.Context, bootAddrs []multiaddr.Multiaddr) (*Node, error) {
	node := &Node{}

	// Not a huge fan of having the ConnManager know about a node, but the
	// manager needs to be able to create connections and periodically find the
	// closest peers. Another approach would be to send in a func or a channel
	// that would be told to find new peers.
	connManager := NewConnManager(node, bootAddrs)
	node.connManager = connManager

	host, err := newHost(ctx, connManager)
	if err != nil {
		return nil, err
	}
	node.host = host

	// Create a multiaddr.Multiaddr with the ID of the node appended on. This is
	// only used when you'd like a Node to connect diretly to anothernode by
	// that address, as in the case of a Node connecting to a Bootstrap Node.
	nodeID := host.ID().Pretty()
	var addrs []multiaddr.Multiaddr
	for _, hostAddress := range host.Addrs() {
		addr, err := multiaddr.NewMultiaddr(
			fmt.Sprintf("%s/ipfs/%s", hostAddress, nodeID))

		if err != nil {
			log.Println("Error creating host address", err)
			return nil, err
		}
		addrs = append(addrs, addr)
	}
	node.addrs = addrs

	return node, nil
}

// newHost returns a libp2p host.Host that wraps a libp2p Swarm Network. The
// libp2p swarm manages groups of connections to peers, and handles incoming and
// outgoing streams
func newHost(ctx context.Context, connManager *ConnManager) (host.Host, error) {
	addr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/0")
	if err != nil {
		return nil, err
	}

	listenAddrs := []multiaddr.Multiaddr{addr}

	// TODO: switch to figcrypto/fastsig.GenerateKey()
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		return nil, err
	}

	pid, err := peer.IDFromPublicKey(priv.GetPublic())
	if err != nil {
		return nil, err
	}

	ps := peerstore.NewPeerstore()
	ps.AddPrivKey(pid, priv)
	ps.AddPubKey(pid, priv.GetPublic())

	muxer := libp2p.DefaultMuxer()
	swrm, err := swarm.NewSwarmWithProtector(ctx, listenAddrs, pid, ps, nil, muxer, nil)
	if err != nil {
		return nil, err
	}

	netw := (*swarm.Network)(swrm)

	return bhost.NewHost(ctx, netw, &bhost.HostOpts{
		ConnManager: connManager,
	})
}

// Addrs returns a slice of fully formed connection strings for use when
// newly connecting a Node to a network.
func (n *Node) Addrs() []multiaddr.Multiaddr {
	return n.addrs
}

// Broadcast NEEDSTOBEIMPLEMENTED
func (n *Node) Broadcast(message []byte) {
}

// ConnectToPeer connects to a peer by the PeerID
func (n *Node) ConnectToPeer(ctx context.Context, peerID peer.ID) error {
	peerInfo := n.Host().Peerstore().PeerInfo(peerID)
	return n.host.Connect(ctx, peerInfo)
}

// ConnectToAddr connects to a peer by supplied addr
func (n *Node) ConnectToAddr(ctx context.Context, addr multiaddr.Multiaddr) error {
	peerInfo, err := peerstore.InfoFromP2pAddr(addr)
	if err != nil {
		return err
	}
	return n.host.Connect(ctx, *peerInfo)
}

// GetClosestPeers returns a list of the closest peers to the node
func (n *Node) GetClosestPeers(ctx context.Context) ([]peer.ID, error) {
	closestPeersChan, err := n.dht.GetClosestPeers(ctx, n.ID().Pretty())
	if err != nil {
		return nil, err
	}

	closestPeers := make([]peer.ID, 0)
	for peer := range closestPeersChan {
		closestPeers = append(closestPeers, peer)
	}
	return closestPeers, nil
}

// Host returns the libp2p Host, which handles Connecting and opening new
// Streams with peers.
func (n *Node) Host() host.Host {
	return n.host
}

// ID returns the unique identifier of a Node on a network.
func (n *Node) ID() peer.ID {
	return n.host.ID()
}

// Send NEEDSTOBEIMPLEMENTED
func (n *Node) Send(message []byte, peerID peer.ID) {
}

// Start forms a connection to a BootstrapNode, based on the bootAddrs that were
// supplied on initialization. Upon connection, the Node will find the closest
// peers with which it should form a permanent connection.
func (n *Node) Start(ctx context.Context) {
	n.dht.Bootstrap(ctx)

	go n.connManager.Start(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}
