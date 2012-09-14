/**
 * Date: 01.09.12
 * Time: 23:11
 *
 * @author Vladimir Matveev
 */
package conf

import (
    "code.google.com/p/gelo"
    "net"
)

// buildConfig parses given Gelo dictionary into actual config structure *Conf. Returns a *Conf pointer
// and a list of config errors.
func buildConfig(d *gelo.Dict) (*Conf, *ConfigErrors) {
    errs := makeConfigErrors()

    cfg := &Conf{
        Listeners: make(map[string]*ListenerConf),
    }

    cfgmap, err := convertDictSafe(d);
    if err != nil {
        return nil, errs.addError(err.(*ConfigError))
    }

    if glmap, ok := cfgmap["listeners"]; ok {
        loadListeners(errs, cfg, glmap.AsDict())
    }

    return cfg, errs
}

// loadListeners fills given *Conf value with listeners configuration from the provided configDict,
// adding errors, if any, to the errs list.
func loadListeners(errs *ConfigErrors, cfg *Conf, lmap map[string]configElement) {
    for lname, gpmap := range lmap {
        lcfg := &ListenerConf{Name: lname, Ports: make(map[PortType]*PortConf)}
        perrs := makeConfigErrors()

        loadListenerPorts(perrs, lcfg, gpmap.AsDict())

        perrs.prependLocation("listener " + lname)
        errs.merge(perrs)

        cfg.Listeners[lname] = lcfg
    }
}

// loadListenerPorts fills given *ListenerConf value with listener ports configuration from
// the provided configDict, adding errors, if any, to the errs list.
func loadListenerPorts(errs *ConfigErrors, lcfg *ListenerConf, pmap map[string]configElement) {
    for pname, gpcfgmap := range pmap {
        pcfgmap := gpcfgmap.AsDict()
        pcfg := new(PortConf)
        location := "port " + pname

        // Try to load port address
        var paddr net.Addr
        switch pname {
        case "tcp4", "tcp6", "udp4", "udp6":
            var phost, pport string
            if gphost, ok := pcfgmap["host"]; !ok {
                errs.add(location, "No 'host' entry present")
                continue
            } else {
                phost = gphost.AsEntry()[0]
                if phost == "*" {
                    switch pname {
                    case "tcp4", "udp4":
                        phost = net.IPv4zero.String()
                    case "tcp6", "udp6":
                        phost = net.IPv6zero.String()
                    }
                }
            }
            if gpport, ok := pcfgmap["port"]; ok {
                pport = gpport.AsEntry()[0]
            } else {
                errs.add(location, "No 'port' entry present")
                continue
            }
            switch pname {
            case "tcp4", "tcp6":
                paddr, _ = net.ResolveTCPAddr(pname, phost + ":" + pport)
            case "udp4", "udp6":
                paddr, _ = net.ResolveUDPAddr(pname, phost + ":" + pport)
            }
        case "unix":
            var ppath string
            if gppath, ok := pcfgmap["path"]; ok {
                ppath = gppath.AsEntry()[0]
            } else {
                errs.add(location, "No 'path' entry present")
                continue
            }
            paddr, _ = net.ResolveUnixAddr(pname, ppath)
        default:
            errs.add("", "Unknown port type '%s'", pname)
            continue
        }

        pcfg.Type = PortType(pname)
        pcfg.Addr = paddr

        lcfg.Ports[pcfg.Type] = pcfg
    }
}

// loadPlugins fills given *Conf value with plugins configuration from the provided configDict,
// adding errors, if any, to the errs list.
func loadPlugins(errs *ConfigErrors, cfg *Conf, pmap map[string]configElement) {
    for pname, gpmap := range pmap {
        pcfg := &PluginConf{
            Name: pname,
            Listeners: make([]string, 0, 2),
            Mediators: make([]MediatorMap, 0, 2),
            Options: make(map[string][]string),
        }
        perrs := makeConfigErrors()

        loadPluginConfig(perrs, pcfg, gpmap.AsDict())

        perrs.prependLocation("plugin " + pname)
        errs.merge(perrs)

        cfg.Plugins[pname] = pcfg
    }
}

// loadPluginConfig fills given *PluginConf value with plugin configuration from provided configDict,
// adding errors, if any, to the errs list.
func loadPluginConfig(errs *ConfigErrors, pcfg *PluginConf, pcfgmap map[string]configElement) {
    if gpplugin, ok := pcfgmap["plugin"]; ok {
        pcfg.Plugin = PluginName(gpplugin.AsEntry()[0])
        delete(pcfgmap, "plugin")
    } else {
        errs.add("", "Plugin name does not set")
        return
    }

    if gplisteners, ok := pcfgmap["listeners"]; ok {
        pcfg.Listeners = gplisteners.AsEntry()
        delete(pcfgmap, "listeners")
    }

    if gpmediators, ok := pcfgmap["mediators"]; ok {
        pmediators := gpmediators.AsEntry()
        if len(pmediators) % 2 != 0 {
            errs.add("", "Mediators config length is odd: %d", len(pmediators))
            return
        }
        for i := range pmediators {
            pcfg.Mediators = append(pcfg.Mediators, MediatorMap{pmediators[i], pmediators[i+1]})
        }
        delete(pcfgmap, "mediators")
    }

    for k, e := range pcfgmap {
        pcfg.Options[k] = e.AsEntry()
    }
}
