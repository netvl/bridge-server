/**
 * Date: 24.08.12
 * Time: 0:03
 *
 * @author Vladimir Matveev
 */
package bridge

import (
    "net"
    "bridge/common/conf"
)

type Bridge struct {
    localPlugins  map[string]LocalPlugin
    remotePlugins map[string]RemotePlugin
    tcpConn net.Listener
}

func New() *Bridge {
    return &Bridge{make(map[string]LocalPlugin), make(map[string]RemotePlugin)}
}

func (b *Bridge) AddLocalPlugin(id string, lp LocalPlugin) {
    b.localPlugins[id] = lp
}

func (b *Bridge) AddRemotePlugin(id string, rp RemotePlugin) {
    b.remotePlugins[id] = rp
}

func (b *Bridge) Start(conf *conf.Conf) error {

}

func (b *Bridge) Stop() {

}
