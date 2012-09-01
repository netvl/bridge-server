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
    plugins      map[string]Plugin
    listeners    []*listener.Listener
    communicator *comm.Communicator
}

func New(cfg *conf.Conf) *Bridge {
    b := &Bridge{
        plugins:      make(map[string]Plugin),
        listeners:    make([]*listener.Listener, 0, 2),
        communicator: comm.NewCommunicator(),
    }

    for _, lconf := range cfg.Listeners {
        b.listeners = append(b.listeners, listener.NewListener(lconf))
    }

    return b
}

func (b *Bridge) AddPlugin(id string, lp Plugin) {
    b.plugins[id] = lp
}

func (b *Bridge) Start() error {
    h := makePluginsHandler(b.plugins, b.communicator)
    for _, l := range b.listeners {
        l.SetHandler(h)
        l.Start()
    }

    return nil
}

func (b *Bridge) Stop() {
    for _, l := range b.listeners {
        l.Stop()
    }
}
