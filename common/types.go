/**
 * Date: 15.09.12
 * Time: 23:48
 *
 * @author Vladimir Matveev
 */
package common

import (
    "bridge/common/conf"
    "net"
    "bridge/common/msg"
)

// ================== BRIDGE ==================

type Bridge interface {
    Start() error
    Stop()
}

type BridgeAPI interface {
    Mediator(name string) Mediator
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
    Send(node Node, msg *msg.Message) error
    Receive(node Node) (*msg.Message, error)
}

// ================== PLUGINS ==================

type PluginType int

const (
    PluginTypeTCP PluginType = iota
    PluginTypeUDP
    PluginTypeUnix
)

var AllPluginTypes = map[PluginType]bool{
    PluginTypeTCP:  true,
    PluginTypeUDP:  true,
    PluginTypeUnix: true,
}

type Plugin interface {
    Name() string
    Config(conf *conf.PluginConf) error
    PluginTypes() map[PluginType]bool
    SupportsMessage(name string) bool
    DeserializeHook() msg.DeserializeHook
    HandleMessage(msg *msg.Message, api BridgeAPI) *msg.Message
    Subscriber(endpoint string) Subscriber
    Term()
}

// ================== MEDIATORS ==================

type Subscriber func (msg interface{})

var EmptySubscriber = func (_ interface{}) {}

type Mediator interface {
    Name() string
    Config(mcfg *conf.MediatorConf) error
    Submit(endpoint string, msg interface{}) error
    Subscribe(endpoint string, s Subscriber) error
    Term()
}

