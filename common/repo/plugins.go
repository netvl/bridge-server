/**
 * Date: 15.09.12
 * Time: 23:46
 *
 * @author Vladimir Matveev
 */
package repo

import (
    . "github.com/dpx-infinity/bridge-server/common"
    "log"
)

type PluginMaker func() Plugin

var pluginsRepo = map[string]PluginMaker{
}

func AddPlugin(name string, maker PluginMaker) {
    if _, present := pluginsRepo[name]; present {
        log.Printf("Plugin with name %s is already in the repo, replacing it", name)
    }
    pluginsRepo[name] = maker
}

func GetPlugin(name string) Plugin {
    pmaker := pluginsRepo[name]
    if pmaker == nil {
        return nil
    }
    return pmaker()
}
