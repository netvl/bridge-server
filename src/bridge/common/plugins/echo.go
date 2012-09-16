/**
 * Date: 31.08.12
 * Time: 1:27
 *
 * @author Vladimir Matveev
 */
package plugins

import (
    . "bridge/common"
    "bridge/common/conf"
    "bridge/common/msg"
)

type EchoPlugin struct{}

func (_ *EchoPlugin) Name() string {
    return "echo"
}

func (_ *EchoPlugin) Config(cfg *conf.PluginConf) error {
    return nil
}

func (_ *EchoPlugin) PluginTypes() map[PluginType]bool {
    return AllPluginTypes
}

func (_ *EchoPlugin) SupportsMessage(name string) bool {
    return true
}

func (_ *EchoPlugin) DeserializeHook() msg.DeserializeHook {
    return msg.EmptyHook
}

func (_ *EchoPlugin) HandleMessage(msg *msg.Message, c Communicator) *msg.Message {
    return msg
}

func (_ *EchoPlugin) Term() {
    return
}
