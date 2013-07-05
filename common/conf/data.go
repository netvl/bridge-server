/**
 * Date: 27.08.12
 * Time: 21:42
 *
 * @author Vladimir Matveev
 */
package conf

import (
    "net"
)

// PortType designates a network of given port. It has the same values
// as standard net functions accept.
type PortType string

// PluginType is an alias for textual name of plugin.
type PluginType string

const (
    PortTypeTCP4 PortType = "tcp4"
    PortTypeUDP4 PortType = "udp4"
    PortTypeTCP6 PortType = "tcp6"
    PortTypeUDP6 PortType = "udp6"
    PortTypeUnix PortType = "unix"
)

var (
    PortTypes = map[PortType]bool{
        PortTypeTCP4: true,
        PortTypeUDP4: true,
        PortTypeTCP6: true,
        PortTypeUDP6: true,
        PortTypeUnix: true,
    }
    TCPPortTypes = map[PortType]bool{
        PortTypeTCP4: true,
        PortTypeTCP6: true,
    }
    UDPPortTypes = map[PortType]bool{
        PortTypeUDP4: true,
        PortTypeUDP6: true,
    }
    UnixPortTypes = map[PortType]bool{
        PortTypeUnix: true,
    }
)

func ValidPortType(s string) (PortType, bool) {
    if pt := PortType(s); PortTypes[pt] {
        return pt, true
    }
    return PortType(""), false
}

type PortConf struct {
    Type PortType
    Addr net.Addr
}

type ListenerConf struct {
    Name  string
    Ports map[PortType]*PortConf
}

type PluginConf struct {
    Name    string
    Plugin  PluginType
    Options map[string][]string
}

type CommonConf struct {
    Discoverable     []uint16
    PresentServices  bool
    StartDebugPlugin []string
}

type PeerConf struct {
    Name   string
    Socket string
}

type LinkConf struct {
    EndA *PeerConf
    EndZ *PeerConf
}

type Conf struct {
    Common    *CommonConf
    Listeners map[string]*ListenerConf
    Plugins   map[string]*PluginConf
    Links     []*LinkConf
}
