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

// SocketOwner represents kind of socket container, either plugin or communicator.
type SocketOwner string

const (
    SocketOwnerPlugin SocketOwner = "plugin"
    SocketOwnerCommunicator SocketOwner = "communicator"
)

// PortType designates a network of given port. It has the same values
// as standard net functions accept.
type PortType string

const (
    PortTypeTCP4 PortType = "tcp4"
    PortTypeUDP4 PortType = "udp4"
    PortTypeTCP6 PortType = "tcp6"
    PortTypeUDP6 PortType = "udp6"
    PortTypeUnix PortType = "unix"
)

var (
    PortTypes     = map[PortType]bool{
        PortTypeTCP4: true,
        PortTypeUDP4: true,
        PortTypeTCP6: true,
        PortTypeUDP6: true,
        PortTypeUnix: true,
    }
    TCPPortTypes  = map[PortType]bool{
        PortTypeTCP4: true,
        PortTypeTCP6: true,
    }
    UDPPortTypes  = map[PortType]bool{
        PortTypeUDP4: true,
        PortTypeUDP6: true,
    }
    UnixPortTypes = map[PortType]bool{
        PortTypeUnix: true,
    }
)

// PluginType is an alias for textual name of plugin.
type PluginType string

// PortTypeFromString returns a pair (portType, ok). If port type designated by the given string exists,
// then ok is true, and portType is that port type. Otherwise, ok is false and portType value is undefined.
func PortTypeFromString(s string) (PortType, bool) {
    if pt := PortType(s); PortTypes[pt] {
        return pt, true
    }
    return PortType(""), false
}

// PortConf represents configuration of communicator port.
type PortConf struct {
    Type PortType
    Addr net.Addr
}

// CommunicatorConf represents configuration of a communicator.
type CommunicatorConf struct {
    Name  string
    Ports map[PortType]*PortConf
}

// PluginConf represents configuration of a plugin.
type PluginConf struct {
    Name    string
    Plugin  PluginType
    Options map[string][]string
}

// CommonConf represents various common options.
type CommonConf struct {
    Discoverable     []uint16
    PresentServices  bool
    StartDebugPlugin []string
}

// PeerConf represents one side of a link, i.e. owner type (plugin or communicator), name of the owner
// and name of the socket.
type PeerConf struct {
    Name   string
    Owner  SocketOwner
    Socket string
}

// LinkConf represents link configuration, i.e. a pair of peers.
type LinkConf struct {
    EndA *PeerConf
    EndZ *PeerConf
}

// Conf represents whole bridge configuration.
type Conf struct {
    Common        *CommonConf
    Communicators map[string]*CommunicatorConf
    Plugins       map[string]*PluginConf
    Links         []*LinkConf
}
