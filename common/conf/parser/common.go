/**
 * Date: 14.01.13
 * Time: 23:33
 *
 * @author Vladimir Matveev
 */
package parser

import (
    "code.google.com/p/gelo"
    "github.com/dpx-infinity/bridge-server/common/conf"
)

func (p *ConfigParser) common(vm *gelo.VM, args *gelo.List, ac uint) gelo.Word {
    if ac != 1 {
        gelo.ArgumentError(vm, "common", "{body}", args)
    }
    body := vm.API.QuoteOrElse(args.Value)

    checkNotInSection(vm, "common")

    enterNamelessSection(vm, "common")

    p.conf.Common = new(conf.CommonConf)
    vm.API.InvokeCmdOrElse(body, nil)

    leaveNamelessSection(vm, "common")

    return nil
}

func (p *ConfigParser) discoverable(vm *gelo.VM, args *gelo.List, ac uint) gelo.Word {
    if ac < 2 || !gelo.StrEqualsSym("at", args.Value.Ser()) {
        gelo.ArgumentError(vm, "discoverable", "at <port> [<port> ...]", args)
    }

    // Check whether we're in common section
    checkInSection(vm, "discoverable", "common")

    args = args.Next // Skip first 'at' symbol

    // Collect all port numbers from the arguments
    var ports []uint16
    for p := args; p != nil; p = p.Next {
        port, ok := p.Value.(*gelo.Number)
        if !ok {
            gelo.RuntimeError(vm, "discoverable arguments should be numbers")
        }
        ports = append(ports, uint16(port.Int()))
    }
    p.conf.Common.Discoverable = ports

    return nil
}

func (p *ConfigParser) presentServices(vm *gelo.VM, args *gelo.List, ac uint) gelo.Word {
    booleanOption(vm, args, ac, &p.conf.Common.PresentServices, "present-services", "common")
    return nil
}

func (p *ConfigParser) startDebugPlugin(vm *gelo.VM, args *gelo.List, ac uint) gelo.Word {
    if ac < 2 || !gelo.StrEqualsSym("on", args.Value.Ser()) {
        gelo.ArgumentError(vm, "start-debug-plugin", "on <listener> [<listener> ...]", args)
    }

    // Check whether we're in common section
    checkInSection(vm, "start-debug-plugin", "common")

    // Collect listener names from arguments
    var names []string
    for p := args; p != nil; p = p.Next {
        name := p.Value.Ser()
        names = append(names, name.String())
    }
    p.conf.Common.StartDebugPlugin = names

    return nil
}
