/**
 * Date: 23.08.12
 * Time: 23:47
 *
 * @author Vladimir Matveev
 */
package bridge

import (
    "bridge/common/msg"
    "bridge/common/net/comm"
)

type PluginType int

const (
    PluginTypeTCP PluginType = iota
    PluginTypeUDP
    PluginTypeUnix
)

var AllPluginTypes = map[PluginType]bool {
    PluginTypeTCP: true,
    PluginTypeUDP: true,
    PluginTypeUnix: true,
}

type LocalPlugin interface {
    Name() string
    PluginTypes() map[PluginType]bool
    SupportsMessage(name string) bool
    DeserializeHook() msg.DeserializeHook
    HandleMessage(msg *msg.Message, c *comm.Communicator) *msg.Message
}

type RemotePlugin interface {
    Name() string
    PluginTypes() map[PluginType]bool
    SupportsMessage(name string) bool
    DeserializeHook() msg.DeserializeHook
    HandleMessage(msg *msg.Message, c *comm.Communicator) *msg.Message
}
