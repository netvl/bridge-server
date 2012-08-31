/**
 * Date: 24.08.12
 * Time: 0:03
 *
 * @author Vladimir Matveev
 */
package bridge

import (
    "bridge/common/conf"
    "bridge/common/net/comm"
    "bridge/common/net/listener"
)

type Bridge struct {
    localPlugins   map[string]LocalPlugin
    remotePlugins  map[string]RemotePlugin
    localListener  listener.Listener
    remoteListener listener.Listener
    communicator   *comm.Communicator
}

func New() *Bridge {
    return &Bridge{
        localPlugins:   make(map[string]LocalPlugin),
        remotePlugins:  make(map[string]RemotePlugin),
        localListener:  listener.NewLocalListener(),
        remoteListener: listener.NewRemoteListener(),
        communicator:   comm.NewCommunicator(),
    }
}

func (b *Bridge) AddLocalPlugin(id string, lp LocalPlugin) {
    b.localPlugins[id] = lp
}

func (b *Bridge) AddRemotePlugin(id string, rp RemotePlugin) {
    b.remotePlugins[id] = rp
}

func (b *Bridge) Start(conf *conf.Conf) error {
    b.localListener.SetHandler(makeLocalPluginsHandler(b.localPlugins, b.communicator))
    b.localListener.Start(conf)

    b.remoteListener.SetHandler(makeRemotePluginsHandler(b.remotePlugins, b.communicator))
    b.remoteListener.Start(conf)

    return nil
}

func (b *Bridge) Stop() {
    b.localListener.Stop()
    b.remoteListener.Stop()
}
