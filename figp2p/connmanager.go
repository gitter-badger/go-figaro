// Package figp2p implements a peer-to-peer network.
package figp2p

// This file is heavily based on https://github.com/libp2p/go-libp2p-connmgr.

import (
	"context"
	"errors"
	"log"
	"sort"
	"sync"
	"time"

	ifconnmgr "github.com/libp2p/go-libp2p-interface-connmgr"
	kb "github.com/libp2p/go-libp2p-kbucket"
	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// ErrConnToBootstrap happens when failing to connect to bootstrap nodes
var ErrConnToBootstrap = errors.New("Unable to connect to boot addresses")

// ErrUntrackedDisconnect happens with an untracked peer disconnects
var ErrUntrackedDisconnect = errors.New("Received disconnected notification for peer we are not tracking")

// ErrFindingNearPeers happens when failing to find nearest peers
var ErrFindingNearPeers = errors.New("Unable to find nearest peers")

// lowWaterConnCount indicates when to connect to more peers
var lowWaterConnCount = 3

// targetConnCount is the optimal number of permanent connections to keep open
var targetConnCount = 7

// highWaterConnCount indicates when to start pruning connections to peers
var highWaterConnCount = 9

// connGracePeriod indicates how long to allow connections from peers that are
// not considered "close" to this peer.
var connGracePeriod = 900 * time.Millisecond

// findNearestPeriod indicates how often to find closer peers than the
// current set of permanent connections
var findNearestPeriod = 180000 * time.Millisecond

// retryAttempts is the number of attempts at connecting to nearest peers
var retryAttempts = 3

// retryConnectPeriod is the time between connection attempts
var retryConnectPeriod = 1 * time.Second

// ConnManager handles connecting to peers that are 'near' to this peer, and
// prunes connections to peers that are 'far'
type ConnManager struct {
	bootAddrs     []multiaddr.Multiaddr
	connCount     int
	connMutex     sync.RWMutex
	peerConnector PeerConnecter
	peerConns     map[peer.ID]map[inet.Conn]time.Time
}

// PeerConnecter handles connecting to and finding peers
type PeerConnecter interface {
	ID() peer.ID
	ConnectToPeer(context.Context, peer.ID) error
	ConnectToAddr(ctx context.Context, addr multiaddr.Multiaddr) error
	GetClosestPeers(ctx context.Context) ([]peer.ID, error)
}

// NewConnManager returns a new ConnManager
func NewConnManager(peerConnector PeerConnecter, bootAddrs []multiaddr.Multiaddr) *ConnManager {
	return &ConnManager{
		bootAddrs:     bootAddrs,
		peerConnector: peerConnector,
		peerConns:     make(map[peer.ID]map[inet.Conn]time.Time),
	}
}

// CMInfo is not implemented
type CMInfo struct{}

// TagInfo is not implemented
type TagInfo struct{}

// checkHighWater removes connections if the number of connections is too high
func (cm *ConnManager) checkHighWater(ctx context.Context) {
	go func() {
		select {
		case <-time.After(connGracePeriod):
			if cm.ConnCount() >= highWaterConnCount {
				cm.TrimOpenConns(ctx)
			}
		case <-ctx.Done():
			return
		}
	}()
}

// checkLowWater adds connections if the number of connections is too low
func (cm *ConnManager) checkLowWater(ctx context.Context) {
	go func() {
		select {
		case <-time.After(connGracePeriod):
			if cm.ConnCount() <= lowWaterConnCount {
				cm.connectToNearestPeers(ctx, retryAttempts, "lowwater")
			}
		case <-ctx.Done():
			return
		}
	}()
}

// ConnCount returns the number of connections
func (cm *ConnManager) ConnCount() int {
	cm.connMutex.RLock()
	defer cm.connMutex.RUnlock()
	return cm.connCount
}

// connectToBootstrapNodes forms a conection with the initialized bootstrap addrs
func (cm *ConnManager) connectToBootstrapNodes(ctx context.Context) {
	if len(cm.bootAddrs) == 0 {
		return
	}

	// Some of the bootAddrs may not be reachable.
	// Allow all but one of the following attempts to fail.
	var errs []error
	for _, bootAddr := range cm.bootAddrs {
		if err := cm.peerConnector.ConnectToAddr(ctx, bootAddr); err != nil {
			errs = append(errs, err)
		}
	}
	if len(cm.bootAddrs) > 0 && len(errs) == len(cm.bootAddrs) {
		log.Panicln(ErrConnToBootstrap, errs)
	}
}

// connectToNearestPeers forms a conection with the nearest peers
func (cm *ConnManager) connectToNearestPeers(ctx context.Context, attempsRemaining int, reason string) {
	if len(cm.bootAddrs) == 0 {
		return
	}

	// Briefly reconnect to the bootstrap nodes to refresh DHT information
	cm.connectToBootstrapNodes(ctx)

	closestPeers, err := cm.peerConnector.GetClosestPeers(ctx)
	if err != nil {
		cm.retryConnectToNearestPeers(ctx, attempsRemaining-1, reason, err)
	}

	for i, peerID := range closestPeers {
		if i < targetConnCount {
			cm.peerConnector.ConnectToPeer(ctx, peerID)
		}
	}

	cm.TrimOpenConns(ctx)
}

// getConnsToClose returns a list of connections that are deemed closable
func (cm *ConnManager) getConnsToClose(ctx context.Context) []inet.Conn {
	cm.connMutex.RLock()
	defer cm.connMutex.RUnlock()

	if cm.ConnCount() < targetConnCount {
		return nil
	}

	closable := make([]inet.Conn, 0)

	// Only add connections outside of the grace period as candidate to be closed
	for _, conns := range cm.peerConns {
		for conn, timeInitiated := range conns {
			if timeInitiated.Add(connGracePeriod).Before(time.Now()) {
				closable = append(closable, conn)
			}
		}
	}

	// Sort the slice so the farthest is at the beginning
	sort.Slice(closable, func(i, j int) bool {
		return kb.Closer(closable[j].RemotePeer(), closable[i].RemotePeer(), cm.peerConnector.ID().Pretty())
	})

	closeCount := len(closable) - targetConnCount

	if closeCount < 1 {
		return []inet.Conn{}
	}

	return closable[:closeCount]
}

// GetInfo is not implemented
func (cm *ConnManager) GetInfo() CMInfo {
	return CMInfo{}
}

// GetTagInfo is not implemented
func (cm *ConnManager) GetTagInfo(p peer.ID) *ifconnmgr.TagInfo {
	return &ifconnmgr.TagInfo{}
}

// Notifee returns an inet.Notifiee that will be notified about any changes
// in connections or streams for a Node
func (cm *ConnManager) Notifee() inet.Notifiee {
	return (*cmNotifee)(cm)
}

// retryConnectToNearestPeers will retry the ConnectToNearestPeers after a certain time period
func (cm *ConnManager) retryConnectToNearestPeers(ctx context.Context, attempsRemaining int, reason string, err error) {
	if attempsRemaining == 0 {
		log.Panicln(ErrFindingNearPeers, err)
	}

	go func() {
		select {
		case <-time.After(retryConnectPeriod):
			cm.connectToNearestPeers(ctx, attempsRemaining-1, reason)
		case <-ctx.Done():
			return
		}
	}()
}

// Start begins the process of opening connections to 'near' peers and pruning
// connections of peers that are 'far'
func (cm *ConnManager) Start(ctx context.Context) {
	cm.connectToNearestPeers(ctx, retryAttempts, "boot")

	for {
		select {
		case <-time.After(findNearestPeriod):
			cm.connectToNearestPeers(ctx, retryAttempts, "periodic")
		case <-ctx.Done():
			return
		}
	}
}

// TagPeer is not implemented
func (cm *ConnManager) TagPeer(p peer.ID, tag string, val int) {}

// trackConnection tracks a connection to a peer
func (cm *ConnManager) trackConnection(conn inet.Conn) {
	cm.connMutex.Lock()
	defer cm.connMutex.Unlock()

	if _, ok := cm.peerConns[conn.RemotePeer()]; !ok {
		cm.peerConns[conn.RemotePeer()] = make(map[inet.Conn]time.Time)
	}

	cm.peerConns[conn.RemotePeer()][conn] = time.Now()
	cm.connCount++
}

// TrimOpenConns closes connections to peers that are not considered 'near' enough
func (cm *ConnManager) TrimOpenConns(ctx context.Context) {
	connsToTrim := cm.getConnsToClose(ctx)

	for _, c := range connsToTrim {
		c.Close()
	}
}

// UntagPeer is not implemented
func (cm *ConnManager) UntagPeer(p peer.ID, tag string) {}

// untrackConnection untracks a connection to a peer
func (cm *ConnManager) untrackConnection(conn inet.Conn) {
	cm.connMutex.Lock()
	defer cm.connMutex.Unlock()

	delete(cm.peerConns[conn.RemotePeer()], conn)

	if len(cm.peerConns[conn.RemotePeer()]) == 0 {
		delete(cm.peerConns, conn.RemotePeer())
	}

	cm.connCount--
}

// cmNotifee is an instance of a ConnManager
type cmNotifee ConnManager

// cm transforms a cmNotifee into a ConnManager
func (nn *cmNotifee) cm() *ConnManager {
	return (*ConnManager)(nn)
}

// ClosedStream is not implemented
func (nn *cmNotifee) ClosedStream(inet.Network, inet.Stream) {}

// Connected is called when a new connection to a peer is formed
func (nn *cmNotifee) Connected(n inet.Network, conn inet.Conn) {
	cm := nn.cm()
	cm.trackConnection(conn)
	cm.checkHighWater(context.Background())
}

// Disconnected is called with a connection to a peer is dropped
func (nn *cmNotifee) Disconnected(n inet.Network, conn inet.Conn) {
	cm := nn.cm()
	cm.untrackConnection(conn)
	cm.checkLowWater(context.Background())
}

// Listen is not implemented
func (nn *cmNotifee) Listen(n inet.Network, addr multiaddr.Multiaddr) {}

// ListenClose is not implemented
func (nn *cmNotifee) ListenClose(n inet.Network, addr multiaddr.Multiaddr) {}

// OpenedStream is not implemented
func (nn *cmNotifee) OpenedStream(inet.Network, inet.Stream) {}
