/**
 * Date: 14.01.13
 * Time: 23:34
 *
 * @author Vladimir Matveev
 */
package parser

import (
    "code.google.com/p/gelo"
    "github.com/dpx-infinity/bridge-server/common/conf"
)

type ConfigParser struct {
    conf *conf.Conf
    port gelo.Port
    vm *gelo.VM
}

func (p *ConfigParser) Init() {
    p.conf = new(conf.Conf)
    p.port = gelo.NewChan()
    p.vm = gelo.NewVM(p.port)

    p.vm.RegisterBundle(p.commands())
}

func (p *ConfigParser) Term() {
    p.vm.Destroy()
}

func (p *ConfigParser) commands() map[string]interface{} {
    return map[string]interface{}{
        "common": p.common,
        "discoverable": p.discoverable,
        "present-services": p.presentServices,
        "start-debug-plugin": p.startDebugPlugin,

        "listeners": p.listeners,
        "listener": p.listener,
        "port": p.port,

        "plugins": p.plugins,
        "plugin": p.plugin,

        "mediators": p.mediators,
        "mediator": p.mediator,
        "endpoints": p.endpoints,

        "links": p.links,
        "connect": p.connect,

        "type": p.typeCmd,
        "set":  p.set,
    }
}

