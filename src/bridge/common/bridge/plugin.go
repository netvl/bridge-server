/**
 * Date: 23.08.12
 * Time: 23:47
 *
 * @author Vladimir Matveev
 */
package bridge

import (
    msg "bridge/common/msg"
)

type PluginType int

const (
    PluginTypeTCP PluginType = iota
    PluginTypeUDP
)

type LocalPlugin interface {
    Name() string
    SupportsMessage(name string) bool
    DeserializeHook() msg.DeserializeHook
    HandleMessage(msg msg.Message) msg.Message
}

type RemotePlugin interface {
    Name() string
    Type() PluginType
    SupportsMessage(name string) bool
    DeserializeHook() msg.DeserializeHook
    HandleMessage(msg msg.Message) msg.Message
}
