/**
 * Date: 01.09.12
 * Time: 23:11
 *
 * @author Vladimir Matveev
 */
package conf

import (
    "code.google.com/p/gelo"
    "fmt"
    "net"
)

type configElement interface {
    AsDict() map[string]configElement
    AsEntry() []string
}

type configDict map[string]configElement

func (d configDict) AsDict() map[string]configElement {
    return map[string]configElement(d)
}

func (d configDict) AsEntry() []string {
    panic(fmt.Sprintf("Attempt to use dict as config entry: %v", d))
    return nil
}

type configEntry []string

func (e configEntry) AsDict() map[string]configElement {
    panic(fmt.Sprintf("Attempt to use config entry as dict: %v", e))
    return nil
}

func (e configEntry) AsEntry() []string {
    return []string(e)
}

func convertDict(d *gelo.Dict) map[string]configElement {
    rm := make(map[string]configElement)

    m := d.Map()
    for k, w := range m {
        switch cw := w.(type) {
        case *gelo.Dict:
            cm := convertDict(cw)
            rm[k] = configDict(cm)
        case *gelo.List:
            ce := convertList(cw)
            rm[k] = configEntry(ce)
        default:
            panic(fmt.Sprintf("Illegal object %v encountered at key %v inside %v dict", w, k, d))
        }
    }

    return rm
}

func convertList(l *gelo.List) []string {
    rl := make([]string, 0, 1)
    for ; l != nil; l = l.Next {
        rl = append(rl, l.Value.Ser().String())
    }
    return rl
}

func buildConfig(d *gelo.Dict) *Conf {
    cfg := &Conf{
        Listeners: make(map[string]*ListenerConf),
    }

    cfgmap := convertDict(d)

    if glmap, ok := cfgmap["listeners"]; ok {
        loadListeners(cfg, glmap.AsDict())
    }

    return cfg
}

func loadListeners(cfg *Conf, lmap map[string]configElement) {
    for lname, gpmap := range lmap {
        lcfg := &ListenerConf{Name: lname}
        loadListenerPorts(lcfg, gpmap.AsDict())
        cfg.Listeners[lname] = lcfg
    }
}

func loadListenerPorts(lcfg *ListenerConf, pmap map[string]configElement) {
    for pname, gpcfgmap := range pmap {
        pcfgmap := gpcfgmap.AsDict()
        pcfg := new(PortConf)

        // Try to load port address
        var paddr net.Addr
        switch pname {
        case "tcp4", "tcp6", "udp4", "udp6":
            var phost, pport string
            if gphost, ok := pcfgmap["host"]; !ok {
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
            if gpport, ok := pcfgmap["port"]; !ok {
                continue
            } else {
                pport = gpport.AsEntry()[0]
            }
            switch pname {
            case "tcp4", "tcp6":
                paddr, _ = net.ResolveTCPAddr(pname, phost + ":" + pport)
            case "udp4", "udp6":
                paddr, _ = net.ResolveUDPAddr(pname, phost + ":" + pport)
            }
        case "unix":
            var ppath string
            if gppath, ok := pcfgmap["path"]; !ok {
                continue
            } else {
                ppath = gppath.AsEntry()[0]
            }
            paddr, _ = net.ResolveUnixAddr(pname, ppath)
        default:
            continue
        }

        pcfg.Type = PortType(pname)
        pcfg.Addr = paddr

        lcfg.Ports = append(lcfg.Ports, pcfg)
    }
}
