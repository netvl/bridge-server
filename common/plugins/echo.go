/**
 * Date: 31.08.12
 * Time: 1:27
 *
 * @author Vladimir Matveev
 */
package plugins

import (
    . "github.com/dpx-infinity/bridge-server/common"
    "github.com/dpx-infinity/bridge-server/common/conf"
    "github.com/dpx-infinity/bridge-server/common/msg"
    "log"
)

type EchoPlugin struct{}

func (_ *EchoPlugin) Name() string {
    return "echo"
}

func (_ *EchoPlugin) Config(cfg *conf.PluginConf) error {
    return nil
}

func (_ *EchoPlugin) SupportsMessage(_ string) bool {
    return true
}

func (_ *EchoPlugin) HandleMessage(msg *msg.Message, _ BridgeAPI) *msg.Message {
    log.Printf("Received a message: %v", msg)
    return msg
}

func (_ *EchoPlugin) Subscriber(_ string) Subscriber {
    return EmptySubscriber
}

func (_ *EchoPlugin) Term() {
    return
}
