/**
 * Date: 05.07.13
 * Time: 22:31
 *
 * @author Vladimir Matveev
 */
package communicators

import (
    . "github.com/dpx-infinity/bridge-server/common"
    "sync"
)

type dynamicPeer struct {
    sockets  map[string]ChanPair
    lock     sync.RWMutex
}

func newDynamicPeer() dynamicPeer {
    return dynamicPeer{sockets: make(map[string]ChanPair)}
}

func (p *dynamicPeer) Sockets() map[string]bool {
    p.lock.RLock()
    defer p.lock.RUnlock()

    names := make(map[string]bool, len(p.sockets))
    for name := range p.sockets {
        names[name] = true
    }
    return names
}

func (p *dynamicPeer) Connect(socket string, linkEnd ChanPair) {
    p.lock.Lock()
    defer p.lock.Unlock()

    p.sockets[socket] = linkEnd
}
