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

// MediatorName is an alias for textual name of mediator.
type MediatorName string

// PluginName is an alias for textual name of plugin.
type PluginName string

const (
    PortTypeTCP4 PortType = "tcp4"
    PortTypeUDP4 PortType = "udp4"
    PortTypeTCP6 PortType = "tcp6"
    PortTypeUDP6 PortType = "udp6"
    PortTypeUnix PortType = "unix"
)

type PortConf struct {
    Type PortType
    Addr net.Addr
}

type ListenerConf struct {
    Name  string
    Ports map[PortType]*PortConf
}

type MediatorConf struct {
    Name string
    Mediator MediatorName
    EndpointNames []string
    Options map[string][]string
}

type MediatorMap struct {
    Mediator string
    Endpoint string
}

type PluginConf struct {
    Name string
    Plugin PluginName
    Listeners []string
    Mediators []MediatorMap
    Options map[string][]string
}

type Conf struct {
    Listeners map[string]*ListenerConf
    Plugins map[string]*PluginConf
    Mediators map[string]*MediatorConf
}
