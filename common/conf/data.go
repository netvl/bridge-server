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

// ===================================================
// ================== COMMUNICATORS ==================
// ===================================================

// CommunicatorConf represents configuration of a communicator.
type CommunicatorConf struct {
    Name      string
    Addresses []string
}

// =============================================
// ================== PLUGINS ==================
// =============================================

// PluginType is an alias for textual name of plugin.
type PluginType string

// PluginConf represents configuration of a plugin.
type PluginConf struct {
    Name    string
    Plugin  PluginType
    Options map[string][]string
}

// ===============================================
// ================== DISCOVERY ==================
// ===============================================

// DiscoveryConf represents configuration of the discovery module.
type DiscoveryConf struct {
    Ports             []int
    DiscoveryInterval uint
    DiscoveryIfaces   []string
    DiscoveryNetworks []*net.IPNet
    ExposeIfaces      []string
    ExposeNetworks    []*net.IPNet
    Statics           []*net.TCPAddr
}

// ===========================================
// ================== LINKS ==================
// ===========================================

// SocketOwner represents kind of socket container, either plugin or communicator.
type SocketOwner string

const (
    SocketOwnerPlugin SocketOwner       = "plugin"
    SocketOwnerCommunicator SocketOwner = "communicator"
)

// PeerConf represents one side of a link, i.e. owner type (plugin or communicator), name of the owner
// and name of the socket.
type PeerConf struct {
    PeerName string
    Owner    SocketOwner
    Socket   string
}

// LinkConf represents link configuration, i.e. a pair of peers.
type LinkConf struct {
    EndA *PeerConf
    EndZ *PeerConf
}

// ============================================
// ================== COMMON ==================
// ============================================

// CommonConf contains general configuration options.
type CommonConf struct {
    Name string
}

// Conf represents whole bridge configuration.
type Conf struct {
    Discovery     *DiscoveryConf
    Communicators map[string]*CommunicatorConf
    Plugins       map[string]*PluginConf
    Links         []*LinkConf
}
