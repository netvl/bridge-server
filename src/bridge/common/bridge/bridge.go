/**
 * Date: 24.08.12
 * Time: 0:03
 *
 * @author Vladimir Matveev
 */
package bridge

import (
    . "bridge/common"
    "bridge/common/conf"
    "bridge/common/net/listener"
    "bridge/common/net/comm"
    "bridge/common/repo"
    "log"
)

type bridge struct {
    listeners    map[string]Listener
    plugins      map[string]Plugin
    mediators    map[string]Mediator
    communicator Communicator
}

func New(cfg *conf.Conf) Bridge {
    b := &bridge{
        listeners:    make(map[string]Listener),
        plugins:      make(map[string]Plugin),
        mediators:    make(map[string]Mediator),
        communicator: comm.NewComm(),
    }

    for lname, lconf := range cfg.Listeners {
        b.listeners[lname] = listener.NewListener(lconf)
    }

    for pname, pconf := range cfg.Plugins {
        plugin := repo.GetPlugin(string(pconf.Plugin))
        if plugin == nil {
            log.Printf("Plugin '%s' of type '%s' not found", pname, pconf.Plugin)
            continue
        }
        if err := plugin.Config(pconf); err != nil {
            log.Printf("Failed to configure '%s' plugin: %v", pname, err)
            continue
        }
        b.plugins[pname] = plugin
    }
    
    for mname, mconf := range cfg.Mediators {
        mediator := repo.GetMediator(string(mconf.Mediator))
        if mediator == nil {
            log.Printf("Mediator '%s' of type '%s' not found", mname, mconf.Mediator)
            continue
        }
        if err := mediator.Config(mconf); err != nil {
            log.Printf("Failed to configure '%s' mediator: %v", mname, err)
            continue
        }
        b.mediators[mname] = mediator
    }

    for pname, pconf := range cfg.Plugins {
        plugin, ok := b.plugins[pname]
        if !ok {
            continue
        }
        for _, mmap := range pconf.Mediators {
            mediator, ok := b.mediators[mmap.Mediator]
            if !ok {
                continue
            }
            endpoint := mmap.Endpoint
            if err := mediator.Subscribe(endpoint, plugin.Subscriber(endpoint)); err != nil {
                log.Printf("Error subscribing plugin '%s' to endpoint '%s' at mediator '%s': %v",
                    pname, endpoint, mediator.Name(), err)
                continue
            }
        }
    }

    return b
}

func (b *bridge) Start() error {
    h := makePluginsHandler(b.plugins, b.communicator)
    for _, l := range b.listeners {
        l.SetHandler(h)
        if err := l.Start(); err != nil {
            return err
        }
    }

    return nil
}

func (b *bridge) Stop() {
    for _, l := range b.listeners {
        l.Stop()
    }
}
