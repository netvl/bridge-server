/**
 * Date: 24.08.12
 * Time: 0:03
 *
 * @author Vladimir Matveev
 */
package bridge

import (
    "bridge/common/conf"
    "bridge/common/net/listener"
    "bridge/common/net/comm"
)

type Bridge struct {
    localPlugins  map[string]LocalPlugin
    remotePlugins map[string]RemotePlugin
    localListener listener.Listener
    remoteListener listener.Listener
    communicator *comm.Communicator
}

func New() *Bridge {
    return &Bridge{
        make(map[string]LocalPlugin),
        make(map[string]RemotePlugin),
        listener.NewLocalListener(),
        listener.NewRemoteListener(),
        comm.NewCommunicator(),
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
