// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package p2p implements the Ethereum p2p network protocols.
package p2p

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"net"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/mclock"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/extdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/discv5"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/p2p/netutil"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	defaultDialTimeout = 15 * time.Second

	// Connectivity defaults.
	maxActiveDialTasks     = 16
	defaultMaxPendingPeers = 50
	defaultDialRatio       = 3

	// Maximum time allowed for reading a complete message.
	// This is effectively the amount of time a connection can be idle.
	frameReadTimeout = 30 * time.Second

	// Maximum amount of time allowed for writing a complete message.
	frameWriteTimeout = 20 * time.Second
)

var errServerStopped = errors.New("server stopped")

// Config holds Server options.
type Config struct {
	// This field must be set to a valid secp256k1 private key.
	PrivateKey *ecdsa.PrivateKey `toml:"-"`

	// MaxPeers is the maximum number of peers that can be
	// connected. It must be greater than zero.
	MaxPeers int

	// MaxPendingPeers is the maximum number of peers that can be pending in the
	// handshake phase, counted separately for inbound and outbound connections.
	// Zero defaults to preset values.
	MaxPendingPeers int `toml:",omitempty"`

	// DialRatio controls the ratio of inbound to dialed connections.
	// Example: a DialRatio of 2 allows 1/2 of connections to be dialed.
	// Setting DialRatio to zero defaults it to 3.
	DialRatio int `toml:",omitempty"`

	// NoDiscovery can be used to disable the peer discovery mechanism.
	// Disabling is useful for protocol debugging (manual topology).
	NoDiscovery bool

	// DiscoveryV5 specifies whether the new topic-discovery based V5 discovery
	// protocol should be started or not.
	DiscoveryV5 bool `toml:",omitempty"`

	// Name sets the node name of this server.
	// Use common.MakeName to create a name that follows existing conventions.
	Name string `toml:"-"`

	// BootstrapNodes are used to establish connectivity
	// with the rest of the network.
	BootstrapNodes []*enode.Node

	// BootstrapNodesV5 are used to establish connectivity
	// with the rest of the network using the V5 discovery
	// protocol.
	BootstrapNodesV5 []*discv5.Node `toml:",omitempty"`

	// Static nodes are used as pre-configured connections which are always
	// maintained and re-connected on disconnects.
	StaticNodes []*enode.Node

	// Trusted nodes are used as pre-configured connections which are always
	// allowed to connect, even above the peer limit.
	TrustedNodes []*enode.Node

	// Connectivity can be restricted to certain IP networks.
	// If this option is set to a non-nil value, only hosts which match one of the
	// IP networks contained in the list are considered.
	NetRestrict *netutil.Netlist `toml:",omitempty"`

	// NodeDatabase is the path to the database containing the previously seen
	// live nodes in the network.
	NodeDatabase string `toml:",omitempty"`

	// Protocols should contain the protocols supported
	// by the server. Matching protocols are launched for
	// each peer.
	Protocols []Protocol `toml:"-"`

	// If ListenAddr is set to a non-nil address, the server
	// will listen for incoming connections.
	//
	// If the port is zero, the operating system will pick a port. The
	// ListenAddr field will be updated with the actual address when
	// the server is started.
	ListenAddr string

	// If set to a non-nil value, the given NAT port mapper
	// is used to make the listening port available to the
	// Internet.
	NAT nat.Interface `toml:",omitempty"`

	// If Dialer is set to a non-nil value, the given Dialer
	// is used to dial outbound peer connections.
	Dialer NodeDialer `toml:"-"`

	// If NoDial is true, the server will not dial any peers.
	NoDial bool `toml:",omitempty"`

	// If EnableMsgEvents is set then the server will emit PeerEvents
	// whenever a message is sent to or received from a peer
	EnableMsgEvents bool

	// Logger is a custom logger to use with the p2p.Server.
	Logger log.Logger `toml:",omitempty"`
}

// Server manages all peer connections.
type Server struct {
	// Config fields may not be modified while the server is running.
	Config

	// Hooks for testing. These are useful because we can inhibit
	// the whole protocol stack.
	newTransport func(net.Conn) transport
	newPeerHook  func(*Peer)

	lock    sync.Mutex // protects running
	running bool

	nodedb       *enode.DB
	localnode    *enode.LocalNode
	ntab         discoverTable
	listener     net.Listener
	ourHandshake *protoHandshake
	lastLookup   time.Time
	DiscV5       *discv5.Network

	// These are for Peers, PeerCount (and nothing else).
	peerOp     chan peerOpFunc
	peerOpDone chan struct{}

	quit          chan struct{}
	addstatic     chan *enode.Node
	removestatic  chan *enode.Node
	addtrusted    chan *enode.Node
	removetrusted chan *enode.Node
	posthandshake chan *conn
	addpeer       chan *conn
	delpeer       chan peerDrop
	loopWG        sync.WaitGroup // loop, listenLoop
	peerFeed      event.Feed
	log           log.Logger
}

type peerOpFunc func(map[enode.ID]*Peer)

type peerDrop struct {
	*Peer
	err       error
	requested bool // true if signaled by the peer
}

type connFlag int32

const (
	dynDialedConn connFlag = 1 << iota
	staticDialedConn
	inboundConn
	trustedConn
)

// conn wraps a network connection with information gathered
// during the two handshakes.
type conn struct {
	fd net.Conn
	transport
	node  *enode.Node
	flags connFlag
	cont  chan error // The run loop uses cont to signal errors to SetupConn.
	caps  []Cap      // valid after the protocol handshake
	name  string     // valid after the protocol handshake
}

type transport interface {
	// The two handshakes.
	doEncHandshake(prv *ecdsa.PrivateKey, dialDest *ecdsa.PublicKey) (*ecdsa.PublicKey, error)
	doProtoHandshake(our *protoHandshake) (*protoHandshake, error)
	// The MsgReadWriter can only be used after the encryption
	// handshake has completed. The code uses conn.id to track this
	// by setting it to a non-nil value after the encryption handshake.
	MsgReadWriter
	// transports must provide Close because we use MsgPipe in some of
	// the tests. Closing the actual network connection doesn't do
	// anything in those tests because MsgPipe doesn't use it.
	close(err error)
}

func (c *conn) String() string {
	s := c.flags.String()
	if (c.node.ID() != enode.ID{}) {
		s += " " + c.node.ID().String()
	}
	s += " " + c.fd.RemoteAddr().String()
	return s
}

func (f connFlag) String() string {
	s := ""
	if f&trustedConn != 0 {
		s += "-trusted"
	}
	if f&dynDialedConn != 0 {
		s += "-dyndial"
	}
	if f&staticDialedConn != 0 {
		s += "-staticdial"
	}
	if f&inboundConn != 0 {
		s += "-inbound"
	}
	if s != "" {
		s = s[1:]
	}
	return s
}

func (c *conn) is(f connFlag) bool {
	flags := connFlag(atomic.LoadInt32((*int32)(&c.flags)))
	return flags&f != 0
}

func (c *conn) set(f connFlag, val bool) {
	for {
		oldFlags := connFlag(atomic.LoadInt32((*int32)(&c.flags)))
		flags := oldFlags
		if val {
			flags |= f
		} else {
			flags &= ^f
		}
		if atomic.CompareAndSwapInt32((*int32)(&c.flags), int32(oldFlags), int32(flags)) {
			return
		}
	}
}

// Peers returns all connected peers.
func (srv *Server) Peers() []*Peer {
	var ps []*Peer
	select {
	// Note: We'd love to put this function into a variable but
	// that seems to cause a weird compiler error in some
	// environments.
	case srv.peerOp <- func(peers map[enode.ID]*Peer) {
		for _, p := range peers {
			ps = append(ps, p)
		}
	}:
		<-srv.peerOpDone
	case <-srv.quit:
	}
	return ps
}

// PeerCount returns the number of connected peers.
func (srv *Server) PeerCount() int {
	var count int
	select {
	case srv.peerOp <- func(ps map[enode.ID]*Peer) { count = len(ps) }:
		<-srv.peerOpDone
	case <-srv.quit:
	}
	return count
}

// AddPeer connects to the given node and maintains the connection until the
// server is shut down. If the connection fails for any reason, the server will
// attempt to reconnect the peer.
func (srv *Server) AddPeer(node *enode.Node) {
	select {
	case srv.addstatic <- node:
	case <-srv.quit:
	}
}

// RemovePeer disconnects from the given node
func (srv *Server) RemovePeer(node *enode.Node) {
	select {
	case srv.removestatic <- node:
	case <-srv.quit:
	}
}

// AddTrustedPeer adds the given node to a reserved whitelist which allows the
// node to always connect, even if the slot are full.
func (srv *Server) AddTrustedPeer(node *enode.Node) {
	select {
	case srv.addtrusted <- node:
	case <-srv.quit:
	}
}

// RemoveTrustedPeer removes the given node from the trusted peer set.
func (srv *Server) RemoveTrustedPeer(node *enode.Node) {
	select {
	case srv.removetrusted <- node:
	case <-srv.quit:
	}
}

// SubscribePeers subscribes the given channel to peer events
func (srv *Server) SubscribeEvents(ch chan *PeerEvent) event.Subscription {
	return srv.peerFeed.Subscribe(ch)
}

// Self returns the local node's endpoint information.
func (srv *Server) Self() *enode.Node {
	srv.lock.Lock()
	ln := srv.localnode
	srv.lock.Unlock()

	if ln == nil {
		return enode.NewV4(&srv.PrivateKey.PublicKey, net.ParseIP("0.0.0.0"), 0, 0)
	}
	return ln.Node()
}

// Stop terminates the server and all active peer connections.
// It blocks until all active connections have been closed.
func (srv *Server) Stop() {
	srv.lock.Lock()
	if !srv.running {
		srv.lock.Unlock()
		return
	}
	srv.running = false
	if srv.listener != nil {
		// this unblocks listener Accept
		srv.listener.Close()
	}
	close(srv.quit)
	srv.lock.Unlock()
	srv.loopWG.Wait()
}

// sharedUDPConn implements a shared connection. Write sends messages to the underlying connection while read returns
// messages that were found unprocessable and sent to the unhandled channel by the primary listener.
type sharedUDPConn struct {
	*net.UDPConn
	unhandled chan discover.ReadPacket
}

// ReadFromUDP implements discv5.conn
func (s *sharedUDPConn) ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error) {
	packet, ok := <-s.unhandled
	if !ok {
		return 0, nil, errors.New("Connection was closed")
	}
	l := len(packet.Data)
	if l > len(b) {
		l = len(b)
	}
	copy(b[:l], packet.Data[:l])
	return l, packet.Addr, nil
}

// Close implements discv5.conn
func (s *sharedUDPConn) Close() error {
	return nil
}

// Start starts running the server.
// Servers can not be re-used after stopping.
func (srv *Server) Start() (err error) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if srv.running {
		return errors.New("server already running")
	}
	srv.running = true
	srv.log = srv.Config.Logger
	if srv.log == nil {
		srv.log = log.New()
	}
	if srv.NoDial && srv.ListenAddr == "" {
		srv.log.Warn("P2P server will be useless, neither dialing nor listening")
	}

	// static fields
	if srv.PrivateKey == nil {
		return errors.New("Server.PrivateKey must be set to a non-nil key")
	}
	if srv.newTransport == nil {
		srv.newTransport = newRLPX
	}
	if srv.Dialer == nil {
		srv.Dialer = TCPDialer{&net.Dialer{Timeout: defaultDialTimeout}}
	}
	srv.quit = make(chan struct{})
	srv.addpeer = make(chan *conn)
	srv.delpeer = make(chan peerDrop)
	srv.posthandshake = make(chan *conn)
	srv.addstatic = make(chan *enode.Node)
	srv.removestatic = make(chan *enode.Node)
	srv.addtrusted = make(chan *enode.Node)
	srv.removetrusted = make(chan *enode.Node)
	srv.peerOp = make(chan peerOpFunc)
	srv.peerOpDone = make(chan struct{})

	if err := srv.setupLocalNode(); err != nil {
		return err
	}
	if srv.ListenAddr != "" {
		if err := srv.setupListening(); err != nil {
			return err
		}
	}
	if err := srv.setupDiscovery(); err != nil {
		return err
	}

	dynPeers := srv.maxDialedConns()
	dialer := newDialState(srv.localnode.ID(), srv.StaticNodes, srv.BootstrapNodes, srv.ntab, dynPeers, srv.NetRestrict)
	srv.loopWG.Add(1)
	go srv.run(dialer)
	return nil
}

func (srv *Server) setupLocalNode() error {
	// Create the devp2p handshake.
	pubkey := crypto.FromECDSAPub(&srv.PrivateKey.PublicKey)
	srv.ourHandshake = &protoHandshake{Version: baseProtocolVersion, Name: srv.Name, ID: pubkey[1:]}
	for _, p := range srv.Protocols {
		srv.ourHandshake.Caps = append(srv.ourHandshake.Caps, p.cap())
	}
	sort.Sort(capsByNameAndVersion(srv.ourHandshake.Caps))

	// Create the local node.
	db, err := enode.OpenDB(srv.Config.NodeDatabase)
	if err != nil {
		return err
	}
	srv.nodedb = db
	srv.localnode = enode.NewLocalNode(db, srv.PrivateKey)
	srv.localnode.SetFallbackIP(net.IP{127, 0, 0, 1})
	srv.localnode.Set(capsByNameAndVersion(srv.ourHandshake.Caps))
	// TODO: check conflicts
	for _, p := range srv.Protocols {
		for _, e := range p.Attributes {
			srv.localnode.Set(e)
		}
	}
	switch srv.NAT.(type) {
	case nil:
		// No NAT interface, do nothing.
	case nat.ExtIP:
		// ExtIP doesn't block, set the IP right away.
		ip, _ := srv.NAT.ExternalIP()
		srv.localnode.SetStaticIP(ip)
	default:
		// Ask the router about the IP. This takes a while and blocks startup,
		// do it in the background.
		srv.loopWG.Add(1)
		go func() {
			defer srv.loopWG.Done()
			if ip, err := srv.NAT.ExternalIP(); err == nil {
				srv.localnode.SetStaticIP(ip)
			}
		}()
	}
	return nil
}

func (srv *Server) setupDiscovery() error {
	if srv.NoDiscovery && !srv.DiscoveryV5 {
		return nil
	}

	addr, err := net.ResolveUDPAddr("udp", srv.ListenAddr)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	realaddr := conn.LocalAddr().(*net.UDPAddr)
	srv.log.Debug("UDP listener up", "addr", realaddr)
	if srv.NAT != nil {
		if !realaddr.IP.IsLoopback() {
			go nat.Map(srv.NAT, srv.quit, "udp", realaddr.Port, realaddr.Port, "ethereum discovery")
		}
	}
	srv.localnode.SetFallbackUDP(realaddr.Port)

	// Discovery V4
	var unhandled chan discover.ReadPacket
	var sconn *sharedUDPConn
	if !srv.NoDiscovery {
		if srv.DiscoveryV5 {
			unhandled = make(chan discover.ReadPacket, 100)
			sconn = &sharedUDPConn{conn, unhandled}
		}
		cfg := discover.Config{
			PrivateKey:  srv.PrivateKey,
			NetRestrict: srv.NetRestrict,
			Bootnodes:   srv.BootstrapNodes,
			Unhandled:   unhandled,
		}
		ntab, err := discover.ListenUDP(conn, srv.localnode, cfg)
		if err != nil {
			return err
		}
		srv.ntab = ntab
	}
	// Discovery V5
	if srv.DiscoveryV5 {
		var ntab *discv5.Network
		var err error
		if sconn != nil {
			ntab, err = discv5.ListenUDP(srv.PrivateKey, sconn, "", srv.NetRestrict)
		} else {
			ntab, err = discv5.ListenUDP(srv.PrivateKey, conn, "", srv.NetRestrict)
		}
		if err != nil {
			return err
		}
		if err := ntab.SetFallbackNodes(srv.BootstrapNodesV5); err != nil {
			return err
		}
		srv.DiscV5 = ntab
	}
	return nil
}

func (srv *Server) setupListening() error {
	// Launch the TCP listener.
	listener, err := net.Listen("tcp", srv.ListenAddr)
	if err != nil {
		return err
	}
	laddr := listener.Addr().(*net.TCPAddr)
	srv.ListenAddr = laddr.String()
	srv.listener = listener
	srv.localnode.Set(enr.TCP(laddr.Port))

	srv.loopWG.Add(1)
	go srv.listenLoop()

	// Map the TCP listening port if NAT is configured.
	if !laddr.IP.IsLoopback() && srv.NAT != nil {
		srv.loopWG.Add(1)
		go func() {
			nat.Map(srv.NAT, srv.quit, "tcp", laddr.Port, laddr.Port, "ethereum p2p")
			srv.loopWG.Done()
		}()
	}
	return nil
}

type dialer interface {
	newTasks(running int, peers map[enode.ID]*Peer, now time.Time) []task
	taskDone(task, time.Time)
	addStatic(*enode.Node)
	removeStatic(*enode.Node)
}

func (srv *Server) run(dialstate dialer) {
	srv.log.Info("Started P2P networking", "self", srv.localnode.Node())
	extdb.SetNodeId(srv.localnode.ID())
	defer srv.loopWG.Done()
	defer srv.nodedb.Close()

	var (
		peers        = make(map[enode.ID]*Peer)
		inboundCount = 0
		trusted      = make(map[enode.ID]bool, len(srv.TrustedNodes))
		taskdone     = make(chan task, maxActiveDialTasks)
		runningTasks []task
		queuedTasks  []task // tasks that can't run yet
	)
	// Put trusted nodes into a map to speed up checks.
	// Trusted peers are loaded on startup or added via AddTrustedPeer RPC.
	for _, n := range srv.TrustedNodes {
		trusted[n.ID()] = true
	}

	// removes t from runningTasks
	delTask := func(t task) {
		for i := range runningTasks {
			if runningTasks[i] == t {
				runningTasks = append(runningTasks[:i], runningTasks[i+1:]...)
				break
			}
		}
	}
	// starts until max number of active tasks is satisfied
	startTasks := func(ts []task) (rest []task) {
		i := 0
		for ; len(runningTasks) < maxActiveDialTasks && i < len(ts); i++ {
			t := ts[i]
			srv.log.Trace("New dial task", "task", t)
			go func() { t.Do(srv); taskdone <- t }()
			runningTasks = append(runningTasks, t)
		}
		return ts[i:]
	}
	scheduleTasks := func() {
		// Start from queue first.
		queuedTasks = append(queuedTasks[:0], startTasks(queuedTasks)...)
		// Query dialer for new tasks and start as many as possible now.
		if len(runningTasks) < maxActiveDialTasks {
			nt := dialstate.newTasks(len(runningTasks)+len(queuedTasks), peers, time.Now())
			queuedTasks = append(queuedTasks, startTasks(nt)...)
		}
	}

running:
	for {
		scheduleTasks()

		select {
		case <-srv.quit:
			// The server was stopped. Run the cleanup logic.
			break running
		case n := <-srv.addstatic:
			// This channel is used by AddPeer to add to the
			// ephemeral static peer list. Add it to the dialer,
			// it will keep the node connected.
			srv.log.Trace("Adding static node", "node", n)
			dialstate.addStatic(n)
		case n := <-srv.removestatic:
			// This channel is used by RemovePeer to send a
			// disconnect request to a peer and begin the
			// stop keeping the node connected.
			srv.log.Trace("Removing static node", "node", n)
			dialstate.removeStatic(n)
			if p, ok := peers[n.ID()]; ok {
				p.Disconnect(DiscRequested)
			}
		case n := <-srv.addtrusted:
			// This channel is used by AddTrustedPeer to add an enode
			// to the trusted node set.
			srv.log.Trace("Adding trusted node", "node", n)
			trusted[n.ID()] = true
			// Mark any already-connected peer as trusted
			if p, ok := peers[n.ID()]; ok {
				p.rw.set(trustedConn, true)
			}
		case n := <-srv.removetrusted:
			// This channel is used by RemoveTrustedPeer to remove an enode
			// from the trusted node set.
			srv.log.Trace("Removing trusted node", "node", n)
			if _, ok := trusted[n.ID()]; ok {
				delete(trusted, n.ID())
			}
			// Unmark any already-connected peer as trusted
			if p, ok := peers[n.ID()]; ok {
				p.rw.set(trustedConn, false)
			}
		case op := <-srv.peerOp:
			// This channel is used by Peers and PeerCount.
			op(peers)
			srv.peerOpDone <- struct{}{}
		case t := <-taskdone:
			// A task got done. Tell dialstate about it so it
			// can update its state and remove it from the active
			// tasks list.
			srv.log.Trace("Dial task done", "task", t)
			dialstate.taskDone(t, time.Now())
			delTask(t)
		case c := <-srv.posthandshake:
			// A connection has passed the encryption handshake so
			// the remote identity is known (but hasn't been verified yet).
			if trusted[c.node.ID()] {
				// Ensure that the trusted flag is set before checking against MaxPeers.
				c.flags |= trustedConn
			}
			// TODO: track in-progress inbound node IDs (pre-Peer) to avoid dialing them.
			select {
			case c.cont <- srv.encHandshakeChecks(peers, inboundCount, c):
			case <-srv.quit:
				break running
			}
		case c := <-srv.addpeer:
			// At this point the connection is past the protocol handshake.
			// Its capabilities are known and the remote identity is verified.
			err := srv.protoHandshakeChecks(peers, inboundCount, c)
			if err == nil {
				// The handshakes are done and it passed all checks.
				p := newPeer(c, srv.Protocols)
				// If message events are enabled, pass the peerFeed
				// to the peer
				if srv.EnableMsgEvents {
					p.events = &srv.peerFeed
				}
				name := truncateName(c.name)
				srv.log.Debug("Adding p2p peer", "name", name, "addr", c.fd.RemoteAddr(), "peers", len(peers)+1)
				go srv.runPeer(p)
				peers[c.node.ID()] = p
				if p.Inbound() {
					inboundCount++
				}
			}
			// The dialer logic relies on the assumption that
			// dial tasks complete after the peer has been added or
			// discarded. Unblock the task last.
			select {
			case c.cont <- err:
			case <-srv.quit:
				break running
			}
		case pd := <-srv.delpeer:
			// A peer disconnected.
			d := common.PrettyDuration(mclock.Now() - pd.created)
			pd.log.Debug("Removing p2p peer", "duration", d, "peers", len(peers)-1, "req", pd.requested, "err", pd.err)
			delete(peers, pd.ID())
			if pd.Inbound() {
				inboundCount--
			}
		}
	}

	srv.log.Trace("P2P networking is spinning down")

	// Terminate discovery. If there is a running lookup it will terminate soon.
	if srv.ntab != nil {
		srv.ntab.Close()
	}
	if srv.DiscV5 != nil {
		srv.DiscV5.Close()
	}
	// Disconnect all peers.
	for _, p := range peers {
		p.Disconnect(DiscQuitting)
	}
	// Wait for peers to shut down. Pending connections and tasks are
	// not handled here and will terminate soon-ish because srv.quit
	// is closed.
	for len(peers) > 0 {
		p := <-srv.delpeer
		p.log.Trace("<-delpeer (spindown)", "remainingTasks", len(runningTasks))
		delete(peers, p.ID())
	}
}

func (srv *Server) protoHandshakeChecks(peers map[enode.ID]*Peer, inboundCount int, c *conn) error {
	// Drop connections with no matching protocols.
	if len(srv.Protocols) > 0 && countMatchingProtocols(srv.Protocols, c.caps) == 0 {
		return DiscUselessPeer
	}
	// Repeat the encryption handshake checks because the
	// peer set might have changed between the handshakes.
	return srv.encHandshakeChecks(peers, inboundCount, c)
}

func (srv *Server) encHandshakeChecks(peers map[enode.ID]*Peer, inboundCount int, c *conn) error {
	switch {
	case !c.is(trustedConn|staticDialedConn) && len(peers) >= srv.MaxPeers:
		return DiscTooManyPeers
	case !c.is(trustedConn) && c.is(inboundConn) && inboundCount >= srv.maxInboundConns():
		return DiscTooManyPeers
	case peers[c.node.ID()] != nil:
		return DiscAlreadyConnected
	case c.node.ID() == srv.localnode.ID():
		return DiscSelf
	default:
		return nil
	}
}

func (srv *Server) maxInboundConns() int {
	return srv.MaxPeers - srv.maxDialedConns()
}
func (srv *Server) maxDialedConns() int {
	if srv.NoDiscovery || srv.NoDial {
		return 0
	}
	r := srv.DialRatio
	if r == 0 {
		r = defaultDialRatio
	}
	return srv.MaxPeers / r
}

// listenLoop runs in its own goroutine and accepts
// inbound connections.
func (srv *Server) listenLoop() {
	defer srv.loopWG.Done()
	srv.log.Debug("TCP listener up", "addr", srv.listener.Addr())

	tokens := defaultMaxPendingPeers
	if srv.MaxPendingPeers > 0 {
		tokens = srv.MaxPendingPeers
	}
	slots := make(chan struct{}, tokens)
	for i := 0; i < tokens; i++ {
		slots <- struct{}{}
	}

	for {
		// Wait for a handshake slot before accepting.
		<-slots

		var (
			fd  net.Conn
			err error
		)
		for {
			fd, err = srv.listener.Accept()
			if netutil.IsTemporaryError(err) {
				srv.log.Debug("Temporary read error", "err", err)
				continue
			} else if err != nil {
				srv.log.Debug("Read error", "err", err)
				return
			}
			break
		}

		// Reject connections that do not match NetRestrict.
		if srv.NetRestrict != nil {
			if tcp, ok := fd.RemoteAddr().(*net.TCPAddr); ok && !srv.NetRestrict.Contains(tcp.IP) {
				srv.log.Debug("Rejected conn (not whitelisted in NetRestrict)", "addr", fd.RemoteAddr())
				fd.Close()
				slots <- struct{}{}
				continue
			}
		}

		var ip net.IP
		if tcp, ok := fd.RemoteAddr().(*net.TCPAddr); ok {
			ip = tcp.IP
		}
		fd = newMeteredConn(fd, true, ip)
		srv.log.Trace("Accepted connection", "addr", fd.RemoteAddr())
		go func() {
			srv.SetupConn(fd, inboundConn, nil)
			slots <- struct{}{}
		}()
	}
}

// SetupConn runs the handshakes and attempts to add the connection
// as a peer. It returns when the connection has been added as a peer
// or the handshakes have failed.
func (srv *Server) SetupConn(fd net.Conn, flags connFlag, dialDest *enode.Node) error {
	c := &conn{fd: fd, transport: srv.newTransport(fd), flags: flags, cont: make(chan error)}
	err := srv.setupConn(c, flags, dialDest)
	if err != nil {
		c.close(err)
		srv.log.Trace("Setting up connection failed", "addr", fd.RemoteAddr(), "err", err)
	}
	return err
}

func (srv *Server) setupConn(c *conn, flags connFlag, dialDest *enode.Node) error {
	// Prevent leftover pending conns from entering the handshake.
	srv.lock.Lock()
	running := srv.running
	srv.lock.Unlock()
	if !running {
		return errServerStopped
	}
	// If dialing, figure out the remote public key.
	var dialPubkey *ecdsa.PublicKey
	if dialDest != nil {
		dialPubkey = new(ecdsa.PublicKey)
		if err := dialDest.Load((*enode.Secp256k1)(dialPubkey)); err != nil {
			return errors.New("dial destination doesn't have a secp256k1 public key")
		}
	}
	// Run the encryption handshake.
	remotePubkey, err := c.doEncHandshake(srv.PrivateKey, dialPubkey)
	if err != nil {
		srv.log.Trace("Failed RLPx handshake", "addr", c.fd.RemoteAddr(), "conn", c.flags, "err", err)
		return err
	}
	if dialDest != nil {
		// For dialed connections, check that the remote public key matches.
		if dialPubkey.X.Cmp(remotePubkey.X) != 0 || dialPubkey.Y.Cmp(remotePubkey.Y) != 0 {
			return DiscUnexpectedIdentity
		}
		c.node = dialDest
	} else {
		c.node = nodeFromConn(remotePubkey, c.fd)
	}
	if conn, ok := c.fd.(*meteredConn); ok {
		conn.handshakeDone(c.node.ID())
	}
	clog := srv.log.New("id", c.node.ID(), "addr", c.fd.RemoteAddr(), "conn", c.flags)
	err = srv.checkpoint(c, srv.posthandshake)
	if err != nil {
		clog.Trace("Rejected peer before protocol handshake", "err", err)
		return err
	}
	// Run the protocol handshake
	phs, err := c.doProtoHandshake(srv.ourHandshake)
	if err != nil {
		clog.Trace("Failed proto handshake", "err", err)
		return err
	}
	if id := c.node.ID(); !bytes.Equal(crypto.Keccak256(phs.ID), id[:]) {
		clog.Trace("Wrong devp2p handshake identity", "phsid", hex.EncodeToString(phs.ID))
		return DiscUnexpectedIdentity
	}
	c.caps, c.name = phs.Caps, phs.Name
	err = srv.checkpoint(c, srv.addpeer)
	if err != nil {
		clog.Trace("Rejected peer", "err", err)
		return err
	}
	// If the checks completed successfully, runPeer has now been
	// launched by run.
	clog.Trace("connection set up", "inbound", dialDest == nil)
	return nil
}

func nodeFromConn(pubkey *ecdsa.PublicKey, conn net.Conn) *enode.Node {
	var ip net.IP
	var port int
	if tcp, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		ip = tcp.IP
		port = tcp.Port
	}
	return enode.NewV4(pubkey, ip, port, port)
}

func truncateName(s string) string {
	if len(s) > 20 {
		return s[:20] + "..."
	}
	return s
}

// checkpoint sends the conn to run, which performs the
// post-handshake checks for the stage (posthandshake, addpeer).
func (srv *Server) checkpoint(c *conn, stage chan<- *conn) error {
	select {
	case stage <- c:
	case <-srv.quit:
		return errServerStopped
	}
	select {
	case err := <-c.cont:
		return err
	case <-srv.quit:
		return errServerStopped
	}
}

// runPeer runs in its own goroutine for each peer.
// it waits until the Peer logic returns and removes
// the peer.
func (srv *Server) runPeer(p *Peer) {
	if srv.newPeerHook != nil {
		srv.newPeerHook(p)
	}

	// broadcast peer add
	srv.peerFeed.Send(&PeerEvent{
		Type: PeerEventTypeAdd,
		Peer: p.ID(),
	})

	// run the protocol
	remoteRequested, err := p.run()

	// broadcast peer drop
	srv.peerFeed.Send(&PeerEvent{
		Type:  PeerEventTypeDrop,
		Peer:  p.ID(),
		Error: err.Error(),
	})

	// Note: run waits for existing peers to be sent on srv.delpeer
	// before returning, so this send should not select on srv.quit.
	srv.delpeer <- peerDrop{p, err, remoteRequested}
}

// NodeInfo represents a short summary of the information known about the host.
type NodeInfo struct {
	ID    string `json:"id"`    // Unique node identifier (also the encryption key)
	Name  string `json:"name"`  // Name of the node, including client type, version, OS, custom data
	Enode string `json:"enode"` // Enode URL for adding this peer from remote peers
	ENR   string `json:"enr"`   // Ethereum Node Record
	IP    string `json:"ip"`    // IP address of the node
	Ports struct {
		Discovery int `json:"discovery"` // UDP listening port for discovery protocol
		Listener  int `json:"listener"`  // TCP listening port for RLPx
	} `json:"ports"`
	ListenAddr string                 `json:"listenAddr"`
	Protocols  map[string]interface{} `json:"protocols"`
}

// NodeInfo gathers and returns a collection of metadata known about the host.
func (srv *Server) NodeInfo() *NodeInfo {
	// Gather and assemble the generic node infos
	node := srv.Self()
	info := &NodeInfo{
		Name:       srv.Name,
		Enode:      node.String(),
		ID:         node.ID().String(),
		IP:         node.IP().String(),
		ListenAddr: srv.ListenAddr,
		Protocols:  make(map[string]interface{}),
	}
	info.Ports.Discovery = node.UDP()
	info.Ports.Listener = node.TCP()
	if enc, err := rlp.EncodeToBytes(node.Record()); err == nil {
		info.ENR = "0x" + hex.EncodeToString(enc)
	}

	// Gather all the running protocol infos (only once per protocol type)
	for _, proto := range srv.Protocols {
		if _, ok := info.Protocols[proto.Name]; !ok {
			nodeInfo := interface{}("unknown")
			if query := proto.NodeInfo; query != nil {
				nodeInfo = proto.NodeInfo()
			}
			info.Protocols[proto.Name] = nodeInfo
		}
	}
	return info
}

// PeersInfo returns an array of metadata objects describing connected peers.
func (srv *Server) PeersInfo() []*PeerInfo {
	// Gather all the generic and sub-protocol specific infos
	infos := make([]*PeerInfo, 0, srv.PeerCount())
	for _, peer := range srv.Peers() {
		if peer != nil {
			infos = append(infos, peer.Info())
		}
	}
	// Sort the result array alphabetically by node identifier
	for i := 0; i < len(infos); i++ {
		for j := i + 1; j < len(infos); j++ {
			if infos[i].ID > infos[j].ID {
				infos[i], infos[j] = infos[j], infos[i]
			}
		}
	}
	return infos
}
