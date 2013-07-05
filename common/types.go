/**
 * Date: 15.09.12
 * Time: 23:48
 *
 * @author Vladimir Matveev
 */
package common

import (
    "github.com/dpx-infinity/bridge-server/common/conf"
    "net"
)

// ============================================
// ================== BRIDGE ==================
// ============================================

type Bridge interface {
    Start() error
    Stop()
}

type BridgeAPI interface {
}

// ===================================================
// ================== COMMUNICATORS ==================
// ===================================================

type Node interface {
    Name() string
    Addr() net.IPAddr
}

// ===========================================
// ================== PORTS ==================
// ===========================================

// Chan represents a channel which can transfer objects of any type in both ways.
type Chan chan interface{}

// SourceChan represents a channel which can only receive objects of any type.
type SourceChan <-chan interface{}

// SinkChan represents a channel which can only send objecs of any type.
type SinkChan chan<- interface{}

// ChanPair interface represents a pair of unidirectional chans with opposite directions.
// It is intended for bidirectional communication between Peers. Each Peer will get one side of Link
// represented with ChanPair, and the Link will be set up in such way that sending to sink part from one
// end will result in message arriving to source part on the other side and vice versa.
type ChanPair interface {
    Source() SourceChan
    Sink() SinkChan
}

// Link represents bidirectional connection between Peers. Link has two sides, A and Z, which
// will be assigned to different Peers and used for bidirectional communication.
type Link interface {
    // EndA returns first link end, opposite to EndZ.
    EndA() ChanPair
    // EndZ returns second link end, opposite to EndA.
    EndZ() ChanPair
}

// Peer represents one side of communication. Peer has a number of sockets (named endpoints)
// to which one side of a Link can be connected. Any number of links can be connected to the single socket;
// when more than one link is connected to single socket, its sinks and sources are aggregated.
type Peer interface {
    // Sockets returns a "set" of socket names this peer have. Note that implementors may return
    // dynamic set of socket names, i.e. the set returned by this method may vary from call to call.
    Sockets() map[string]bool
    // Connect attaches given link end (represented by a ChanPair) to the specified socket.
    Connect(socket string, linkEnd ChanPair)
}

// ===================================================
// ================== COMMUNICATORS ==================
// ===================================================

// Handler is a function which is able to handle standard connection.
// It is supposed that the handler itself does not close the connection.
type Handler func(net.Conn)

// Communicator exposes networking interface to the plugins. When started on host network interface (or on all
// interfaces at once) it accepts messages incoming through network, routing them through its Peer interface
// to all attached sockets.
// It also accepts messages from any of its sockets, decoding value of 'Destination' header and trying
// to send it through the interface it is configured on. If 'Destination' header is not present,
// tries to send the message in some appropriate default way.
type Communicator interface {
    Peer
    Start() error
    Stop()
}

// =============================================
// ================== PLUGINS ==================
// =============================================

// Plugin represents an entity which handles messages incoming from listeners
// and performs some useful work on it. This interface is the main extension point
// of bridge. Different plugins can do virtually anything.
type Plugin interface {
    Peer
    Name() string
    Init(conf *conf.PluginConf, api BridgeAPI) error
    Term()
}
