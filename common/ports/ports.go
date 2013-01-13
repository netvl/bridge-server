/**
 * Date: 14.01.13
 * Time: 1:35
 *
 * @author Vladimir Matveev
 */
package ports

import (
    . "github.com/dpx-infinity/bridge-server/common"
)

// A container for two channels, source and sink.
// It serves both as a Port (provides an access to First peer and Second peer sides)
// and as a ChanPair for the First peer.
type port struct {
    source Chan
    sink Chan
}

func (p *port) First() ChanPair {
    return p
}

func (p *port) Second() ChanPair {
    return reversed(p)
}

func (p *port) Source() SourceChan {
    return p.source
}

func (p *port) Sink() SinkChan {
    return p.sink
}

// A wrapper around port which serves as a ChanPair for the Second peer.
type reversed port

func (p *reversed) Source() SourceChan {
    return p.sink
}

func (p *reversed) Sink() SinkChan {
    return p.source
}

func NewPort() Port {
    return &port{make(Chan), make(Chan)}
}
