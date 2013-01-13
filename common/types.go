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
    "github.com/dpx-infinity/bridge-server/common/msg"
)

// ================== BRIDGE ==================

type Bridge interface {
    Start() error
    Stop()
}

type BridgeAPI interface {
    Comm() Communicator
}

// ================== LISTENERS ==================

// Handler is a function which is able to handle standard connection.
// It is supposed that the handler itself does not close the connection.
type Handler func(net.Conn)

type Listener interface {
    Start() error
    Stop()
    SetHandler(handler Handler)
}

// ================== COMMUNICATORS ==================

type Node interface {
    Name() string
    Addr() net.IPAddr
}

type Communicator interface {
    Send(node Node, msg *msg.Message)
}

// ================== PORTS ==================

type Chan chan interface{}
type SourceChan <-chan interface{}
type SinkChan chan<- interface{}

type ChanPair interface {
    Source() SourceChan
    Sink() SinkChan
}

type Port interface {
    First() ChanPair
    Second() ChanPair
}

// ================== PLUGINS ==================

type Plugin interface {
    Name() string
    Config(conf *conf.PluginConf) error
    Port(name string) Port
    SupportsMessage(name string) bool
    HandleMessage(msg *msg.Message, api BridgeAPI) *msg.Message
    Term()
}

// ================== MEDIATORS ==================

type Mediator interface {
    Name() string
    Config(conf *conf.MediatorConf) error
    HasEndpoint(endpoint string) bool
    Connect(endpoint string, port Port) error
    Term()
}

