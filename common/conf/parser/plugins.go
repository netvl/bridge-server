/**
 * Date: 15.01.13
 * Time: 19:27
 *
 * @author Vladimir Matveev
 */
package parser

import (
    "code.google.com/p/gelo"
    . "github.com/dpx-infinity/bridge-server/common"
    "github.com/dpx-infinity/bridge-server/common/conf"
)

func (p *ConfigParser) plugins(vm *gelo.VM, args *gelo.List, ac uint) gelo.Word {
    if ac != 1 {
        gelo.ArgumentError(vm, "plugins", "{body}", args)
    }
    body := vm.API.QuoteOrElse(args.Value)

    checkNotInSection(vm, "plugins")

    insideNamelessSection(vm, "plugins",
        func() {
            p.conf.Plugins = make(map[string]*conf.PluginConf)
            vm.API.InvokeCmdOrElse(body, nil)
        },
    )

    return nil
}

func (p *ConfigParser) plugin(vm *gelo.VM, args *gelo.List, ac uint) gelo.Word {
    if ac != 2 {
        gelo.ArgumentError(vm, "plugin", "<name> {body}", args)
    }
    name := args.Value.Ser().String()
    body := vm.API.QuoteOrElse(args.Next.Value)

    checkInSection(vm, "plugin", "plugins")

    insideSection(vm, "plugin", name,
        func() {
            p.conf.Plugins[name] = new(conf.PluginConf)
            p.conf.Plugins[name].Name = name
            vm.API.InvokeCmdOrElse(body, args)

            p.conf.Plugins[name].Options = make(map[string][]string)
            d := getOrMakeDict(vm, "data")
            for k, v := range d.Map() {
                var elements []string
                args.Slice()
                for e := vm.API.ListOrElse(v); e != nil; e = e.Next {
                    elements = append(elements, e.Value.Ser())
                }
                p.conf.Plugins[name].Options[k] = elements
            }
        },
    )

    return nil
}

func (p *ConfigParser) pluginType(vm *gelo.VM, args *gelo.List, ac uint) gelo.Word {
    if ac != 1 {
        gelo.ArgumentError(vm, "type", "<type>", args)
    }
    pluginType := args.Value.Ser().String()

    pluginName := checkInSection(vm, "type", "plugin")

    p.conf.Plugins[pluginName].Plugin = conf.PluginType(pluginType)

    return nil
}
