package listener

import (
    . "github.com/dpx-infinity/bridge-server/common"
    "github.com/dpx-infinity/bridge-server/common/conf"
    "log"
    "net"
)

type port struct {
    network  string
    addr     net.Addr
    stopChan StopChan
    listener net.Listener
}

type listener struct {
    name    string
    handler Handler
    ports   []*port
}

func NewListener(cfg *conf.ListenerConf) Listener {
    l := new(listener)
    l.name = cfg.Name
    l.ports = make([]*port, 0, 8)
    for _, pc := range cfg.Ports {
        port := &port{
            network:  string(pc.Type),
            addr:     pc.Addr,
            stopChan: make(StopChan),
        }
        l.ports = append(l.ports, port)
    }
    return l
}

func (l *listener) SetHandler(handler Handler) {
    l.handler = handler
}

func (l *listener) Start() error {
    for _, port := range l.ports {
        pl, err := net.Listen(port.network, port.addr.String())
        if err != nil {
            // TODO: proper error handling
            log.Printf("Error starting listener for port of type %v for Listener %v", port.network, l.name)
            return err
        }
        port.listener = pl

        go func() {
            listenOn(pl, port.stopChan, l.handler)

            if err := pl.Close(); err != nil {
                // TODO: error handling
            }
        }()
    }

    return nil
}

func (l *listener) Stop() {
    for _, port := range l.ports {
        if port.stopChan != nil {
            port.stopChan.Stop()
            port.stopChan = make(StopChan)
            port.listener.Close()
        }
    }
}
