/**
 * Date: 15.09.12
 * Time: 23:46
 *
 * @author Vladimir Matveev
 */
package repo

import (
    . "github.com/dpx-infinity/bridge-server/common"
    "github.com/dpx-infinity/bridge-server/common/plugins"
)

type PluginMaker func () Plugin

var pluginsRepo = map[string]PluginMaker {
    "echo": func () Plugin { return new(plugins.EchoPlugin) },
}

func GetPlugin(name string) Plugin {
    pmaker := pluginsRepo[name]
    if pmaker == nil {
        return nil
    }
    return pmaker()
}
