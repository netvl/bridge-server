/**
 * Date: 24.08.12
 * Time: 0:03
 *
 * @author Vladimir Matveev
 */
package common

import (
    "net"
)

type Node interface {
    GetName() string
    GetAddress() net.IPAddr
}

type Bridge interface {
    DiscoverNodes() []Node
    RegisterLocalPlugin(plugin LocalPlugin)
    RegisterRemotePlugin(plugin RemotePlugin)
}

type bridge struct {
    nodes []Node
    localPlugins []LocalPlugin
    remotePlugins []RemotePlugin
}
