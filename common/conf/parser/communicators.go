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
    "net"
)

func (p *ConfigParser) communicators(vm *gelo.VM, args *gelo.List, ac uint) gelo.Word {
    if ac != 1 {
        gelo.ArgumentError(vm, "communicators", "{body}", args)
    }
    body := vm.API.QuoteOrElse(args.Value)

    checkNotInSection(vm, "communicators")

    insideNamelessSection(vm, "communicators",
        func() {
            p.conf.Listeners = make(map[string]*conf.ListenerConf)
            vm.API.InvokeCmdOrElse(body, nil)
        },
    )

    return nil
}

func (p *ConfigParser) communicator(vm *gelo.VM, args *gelo.List, ac uint) gelo.Word {
    if ac != 2 {
        gelo.ArgumentError(vm, "communicator", "<name> {body}", args)
    }
    name := args.Value.Ser().String()
    body := vm.API.QuoteOrElse(args.Next.Value)

    checkInSection(vm, "communicator", "communicators")

    insideSection(vm, "communicator", name,
        func() {
            p.conf.Listeners[name] = new(conf.ListenerConf)
            p.conf.Listeners[name].Name = name
            p.conf.Listeners[name].Ports = make(map[conf.PortType]*conf.PortConf)
            vm.API.InvokeCmdOrElse(body, nil)
        },
    )

    return nil
}

func (p *ConfigParser) port(vm *gelo.VM, args *gelo.List, ac uint) gelo.Word {
    if ac != 2 {
        gelo.ArgumentError(vm, "port", "<type> {body}", args)
    }
    portType := args.Value.Ser().String()
    body := vm.API.QuoteOrElse(args.Next.Value)

    listenerName := checkInSection(vm, "port", "communicator")

    insideSection(vm, "port", portType,
        func() {
            vm.API.InvokeCmdOrElse(body, nil)
            d := getOrMakeDict(vm, "data")

            if pt, ok := conf.ValidPortType(portType); ok {
                var addr net.Addr

                // Try to resolve the address specified in the port configuration
                if conf.TCPPortTypes[pt] || conf.UDPPortTypes[pt] { // If port is network-related
                    portHost, ok := vm.API.ListOrElse(d.StrGet("host")).Value.Ser().String()
                    if !ok {
                        runtimeError(vm,
                            "Port host address is not specified: communicator %s, port %s",
                            listenerName, portType,
                        )
                    }

                    portNum, ok := vm.API.ListOrElse(d.StrGet("port")).Value.Ser().String()
                    if !ok {
                        runtimeError(vm,
                            "Port number is not specified: communicator %s, port %s",
                            listenerName, portType,
                        )
                    }

                    if conf.TCPPortTypes[pt] {
                        if addr, err := net.ResolveTCPAddr(portType, portHost+":"+portNum); err != nil {
                            runtimeError(vm,
                                "Unable to resolve TCP address: %s: communicator %s, port %s",
                                err.Error(), listenerName, portType,
                            )
                        }
                    } else if conf.UDPPortTypes[pt] {
                        if addr, err := net.ResolveUDPAddr(portType, portHost+":"+portNum); err != nil {
                            runtimeError(vm,
                                "Unable to resolve UDP address: %s: communicator %s, port %s",
                                err.Error(), listenerName, portType,
                            )
                        }
                    }

                } else if conf.UnixPortTypes[pt] { // If port is local
                    portPath, ok := vm.API.ListOrElse(d.StrGet("path")).Value.Ser().String()
                    if !ok {
                        runtimeError(vm,
                            "Port file path is not specified: communicator %s, port %s",
                            listenerName, portType,
                        )
                    }

                    if addr, err := net.ResolveUnixAddr(portType, portPath); err != nil {
                        runtimeError(vm,
                            "Unable to resolve Unix address: %s: communicator %s, port %s",
                            err.Error(), listenerName, portType,
                        )
                    }
                }

                // Resolution was successful, assing port configuration
                p.conf.Listeners[listenerName].Ports[pt] = &conf.PortConf{
                    Type: pt,
                    Addr: addr,
                }
            } else {
                runtimeError(vm, "Invalid port type: %s; communicator %s", portType, listenerName)
            }
        },
    )

    return nil
}
