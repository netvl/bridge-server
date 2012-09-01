/**
 * Date: 31.08.12
 * Time: 1:27
 *
 * @author Vladimir Matveev
 */
package plugins

import (
    "bridge/common/bridge"
    "bridge/common/msg"
    "bridge/common/net/comm"
)

type EchoPlugin struct{}

func (_ *EchoPlugin) Name() string {
    return "echo"
}

func (_ *EchoPlugin) PluginTypes() map[bridge.PluginType]bool {
    return bridge.AllPluginTypes
}

func (_ *EchoPlugin) SupportsMessage(name string) bool {
    return true
}

func (_ *EchoPlugin) DeserializeHook() msg.DeserializeHook {
    return msg.EmptyHook
}

func (_ *EchoPlugin) HandleMessage(msg *msg.Message, c *comm.Communicator) *msg.Message {
    return msg
}
