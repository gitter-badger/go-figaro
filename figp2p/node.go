package figp2p

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-peer"

	"github.com/multiformats/go-multihash"

	cid "github.com/ipfs/go-cid"
	datastore "github.com/ipfs/go-datastore"
	"github.com/libp2p/go-floodsub"
	"github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
)

const connectTimeout time.Duration = time.Second * 10

// ErrConnectingToBootstrap indicate a node was unable to connect to any of the specified bootstrap addresses
var ErrConnectingToBootstrap = errors.New("Unable to connect to a bootstrap node")

// Node implements figaro.Node
type Node struct {
	addresses     []multiaddr.Multiaddr
	dht           *dht.IpfsDHT
	fsub          *floodsub.PubSub
	handlers      *FanoutMux
	host          *host.Host
	newMessage    chan *floodsub.Message
	peers         []peerstore.PeerInfo
	subscriptions map[string]*floodsub.Subscription
}

// NewBootstrapNode creates a new Node that is ready to receive peers
func NewBootstrapNode(ctx context.Context, handlers *FanoutMux) (*Node, error) {
	node, err := newNode(ctx, handlers)
	if err != nil {
		return nil, err
	}

	log.Println(node.ID(), "Setting as a bootstrap node")
	node.dht = dht.NewDHT(ctx, *node.host, datastore.NewMapDatastore())

	return node, nil
}

// NewNode creates a new Node that connects to a bootstrap Node
func NewNode(ctx context.Context, bootstrapAddresses []string, rendezvousName string, handlers *FanoutMux) (*Node, error) {
	node, err := newNode(ctx, handlers)
	if err != nil {
		return nil, err
	}

	node.dht = dht.NewDHTClient(ctx, *node.host, datastore.NewMapDatastore())

	connectSuccess := false
	for _, bootstrapAddress := range bootstrapAddresses {
		bootMultiAddr, err := multiaddr.NewMultiaddr(bootstrapAddress)
		if err != nil {
			return nil, err
		}

		pinfo, _ := peerstore.InfoFromP2pAddr(bootMultiAddr)
		if err := (*node.host).Connect(ctx, *pinfo); err != nil {
			log.Println(node.ID(), "Error connecting to peer", bootstrapAddress, err)
		}
		connectSuccess = true
		break
	}

	if connectSuccess != true {
		return nil, ErrConnectingToBootstrap
	}

	rdvAddr, _ := cid.NewPrefixV1(cid.Raw, multihash.SHA2_256).Sum([]byte(rendezvousName))

	tctx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	log.Println(node.ID(), "Announcing self")
	if err := node.dht.Provide(tctx, rdvAddr, true); err != nil {
		log.Println(node.ID(), "Error announcing self", err)
		return nil, err
	}

	tctx, cancel = context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	log.Println(node.ID(), "Finding peers")
	peers, err := node.dht.FindProviders(tctx, rdvAddr)
	if err != nil {
		log.Println(node.ID(), "Error findind peers", err)
		return nil, err
	}
	log.Println(node.ID(), fmt.Sprintf("Found %d peers", len(peers)))

	for _, peer := range peers {
		node.peers = append(node.peers, peer)
	}

	return node, nil
}

func newNode(ctx context.Context, handlers *FanoutMux) (*Node, error) {
	// Create a Host
	log.Println("Creating a new host")
	host, err := libp2p.New(ctx, libp2p.Defaults)
	if err != nil {
		log.Println("Error creating a new host", err)
		return nil, err
	}
	nodeID := host.ID().Pretty()
	log.Println(nodeID, "Created a new host")

	// Create the Node's Addresses
	var addresses []multiaddr.Multiaddr
	for _, hostAddress := range host.Addrs() {
		addr, err := multiaddr.NewMultiaddr(fmt.Sprintf("%s/ipfs/%s", hostAddress, nodeID))

		if err != nil {
			log.Println("Error creating host address", err)
			return nil, err
		}
		addresses = append(addresses, addr)
	}

	// Create a PubSub, attaching the host
	log.Println(nodeID, "Creating a new pubsub")
	fsub, err := floodsub.NewFloodSub(ctx, host)
	if err != nil {
		log.Println("Error creating a new pubsub", err)
		return nil, err
	}

	// Replace the nodeIDTopic placeholder with the actual nodeID
	handlers.Set(nodeID, handlers.Get(nodeIDTopic))
	handlers.Delete(nodeIDTopic)

	return &Node{
		addresses:     addresses,
		fsub:          fsub,
		host:          &host,
		handlers:      handlers,
		newMessage:    make(chan *floodsub.Message),
		subscriptions: make(map[string]*floodsub.Subscription),
	}, nil
}

// Bootstrap starts a standalone network that will be connected to by other Nodes
func (n *Node) Bootstrap(ctx context.Context) {
	log.Println(n.ID(), "Setting as a bootstrap node")
	n.dht = dht.NewDHT(ctx, *n.host, datastore.NewMapDatastore())
}

// Broadcast sends out a message to all peer Nodes that are subscribed to the given topic
func (n *Node) Broadcast(ctx context.Context, topicName string, msg []byte) error {
	log.Println(n.ID(), "Broadcasting message: ", string(msg))
	return n.fsub.Publish(topicName, msg)
}

// FullAddresses returns a string version of the Node's transport address
func (n *Node) FullAddresses() []string {
	var addresses []string
	for _, addr := range n.addresses {
		addresses = append(addresses, addr.String())
	}
	return addresses
	// return (*n.addr).String()
}

// ID returns a string version of a Node's ID
func (n *Node) ID() string {
	return (*n.host).ID().Pretty()
}

// Listen subscribes to all topics preregisted through the FanoutMux, starts
// listening to incoming messages, and routes them to the appropriate Handler
func (n *Node) Listen(ctx context.Context) {
	n.subscribeToRegisteredHandlers()

	for _, sub := range n.subscriptions {
		go n.listenToSubscription(ctx, sub)
	}

	n.listenForIncomingMessages(ctx)
}

// listenForIncomingMessages reads messages in priority order and calls the
// appropriate Handler, optionally closing when the context ends
func (n *Node) listenForIncomingMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			break
		// case newConnection<-newConnectionChan:
		//  DOSOMETHING
		// 	break
		case msg := <-n.newMessage:
			topicName := msg.GetTopicIDs()[0]
			msgHandlers := n.handlers.Get(topicName)
			for _, handler := range msgHandlers {
				// Should we be calling these handlers asynchronously?
				// go handler(msg)
				handler(n, msg)
			}
		}
	}
}

// listenToSubscription waits for new messages on a subscription
// and sends the messages into the Node's Message channel
func (n *Node) listenToSubscription(ctx context.Context, sub *floodsub.Subscription) {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			break
		}
		n.newMessage <- msg
	}
}

// PeerID returns the peer.ID of a Node
func (n *Node) PeerID() peer.ID {
	return (*n.host).ID()
}

// Send sends a message directly to a peer Node
func (n *Node) Send(peerID peer.ID, data []byte) {
	log.Println(n.ID(), "Sending message", string(data))
	n.fsub.Publish(peerID.Pretty(), data)
}

// subscribe subscribes a Node to a particular topic
func (n *Node) subscribe(topicName string) error {
	if sub := n.subscriptions[topicName]; sub != nil {
		log.Println(n.ID(), "Already subscribed to", topicName)
		return nil
	}

	log.Println(n.ID(), "Subscribing to", topicName)
	sub, err := n.fsub.Subscribe(topicName)
	if err != nil {
		return err
	}
	n.subscriptions[topicName] = sub

	return nil
}

// subscribeToRegisteredHandlers runs through all registred
// handlers and subscribes the node to the appropriate topic
func (n *Node) subscribeToRegisteredHandlers() error {
	for _, key := range n.handlers.Keys() {
		err := n.subscribe(key)
		if err != nil {
			return err
		}
	}
	return nil
}
