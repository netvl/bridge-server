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
    "io"
    "os"
)

type ConfigParser struct {
    conf *conf.Conf
    port gelo.Port
    vm   *gelo.VM
}

func (p *ConfigParser) Init() {
    p.conf = new(conf.Conf)
    p.port = gelo.NewChan()
    p.vm = gelo.NewVM(p.port)

    p.vm.RegisterBundle(p.commands())
}

func (p *ConfigParser) Load(src io.Reader) error {
    _, err := p.vm.Run(src, nil)
    if err != nil {
        return err
    }

    return nil
}

func (p *ConfigParser) LoadFromFile(file string) error {
    f, err := os.Open(file)
    if err != nil {
        return err
    }
    defer f.Close()

    return p.Load(f)
}

func (p *ConfigParser) Term() {
    p.vm.Destroy()
}

func (p *ConfigParser) commands() map[string]interface{} {
    return map[string]interface{}{
        "common":             p.common,
        "discoverable":       p.discoverable,
        "present-services":   p.presentServices,
        "start-debug-plugin": p.startDebugPlugin,

        "communicators": p.communicators,
        "communicator":  p.communicator,
        "port":      p.port,

        "plugins": p.plugins,
        "plugin":  p.plugin,
        "type":    p.pluginType,

        "links":   p.links,
        "connect": p.connect,

        "set": p.set,
    }
}
