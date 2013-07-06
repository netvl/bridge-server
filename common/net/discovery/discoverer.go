/**
 * Date: 06.07.13
 * Time: 20:52
 *
 * @author Vladimir Matveev
 */
package discovery

import (
    . "github.com/dpx-infinity/bridge-server/common"
    "github.com/dpx-infinity/bridge-server/common/conf"
    "net"
    "sync"
    "github.com/dpx-infinity/bridge-server/common/util"
    "github.com/dpx-infinity/bridge-server/common/msg"
    "github.com/dpx-infinity/bridge-server/common/log"
)

type discoverer struct {
    nodes    []Node
    ifaces   []string
    networks []*net.IPNet
    statics  []*net.TCPAddr
    ports    []int
    lock     sync.RWMutex
    stopChan util.StopChan
}

func NewDiscoverer(config *conf.DiscoveryConf) Discoverer {
    discoverer := &discoverer{
        nodes:    make([]Node, 2),
        ifaces:   make([]string, len(config.DiscoveryIfaces)),
        networks: make([]*net.IPNet, len(config.DiscoveryNetworks)),
        statics:  make([]*net.TCPAddr, len(config.Statics)),
        ports:    make([]int, len(config.Ports)),
    }
    copy(discoverer.ifaces, config.DiscoveryIfaces)
    copy(discoverer.networks, config.DiscoveryNetworks)
    copy(discoverer.statics, config.Statics)
    copy(discoverer.ports, config.Ports)

    return discoverer
}

func (d *discoverer) Start() error {
    go func() {
        for {
            if d.stopChan.Stopped() {
                break
            }

            nodes := make([]Node, len(d.nodes))

            // Try discovering on interfaces
            // TODO

            // Try discovering on networks
            // TODO

            // Try discovering on statics
            for _, static := range d.statics {
                node, err := readMessage(static)
                if err != nil {
                    // TODO: handle error property
                    log.Errorln("Cannot resolve static node", err)
                } else {
                    nodes = append(nodes, node)
                }
            }

            // Move newly collected nodes to the discoverer
            d.lock.Lock()
            d.nodes = nodes
            d.lock.Unlock()
        }
    }()

    return nil
}

func readMessage(addr *net.TCPAddr) (Node, error) {
    conn, err := net.DialTCP(addr.Network(), nil, addr)
    if err != nil {
        return nil, err
    }

    message, err := msg.Deserialize(conn)
    if err != nil {
        return nil, err
    }

    node, err := NodeFromMessage(message)
    if err != nil {
        return nil, err
    }

    return node, nil
}

func (d *discoverer) Stop() {
    d.stopChan.Stop()
}

func (d *discoverer) Nodes() (result []Node) {
    d.lock.RLock()
    defer d.lock.RUnlock()

    result := make([]Node, len(d.nodes))
    copy(result, d.nodes)
    return
}
