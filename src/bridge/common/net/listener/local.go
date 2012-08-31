/**
 * Date: 27.08.12
 * Time: 22:26
 *
 * @author Vladimir Matveev
 */
package listener

import (
    "net"
    "bridge/common/conf"
    "log"
)

type localListener struct {
    handler      Handler
    tcpStopChan  StopChan
    tcpListener  net.Listener
    unixStopChan StopChan
    unixListener net.Listener
}

// NewLocalListener creates new local listener instance, that is, that kind of listener
// which uses local part of configuration.
func NewLocalListener() Listener {
    return new(localListener)
}

// SetHandler sets the handler which local listener will use when accepting connections.
func (ll *localListener) SetHandler(handler Handler) {
    ll.handler = handler
}

// Start initializes listener from the given configuration and starts listening
// on configured addresses. When an incoming connection comes, it will be handled
// by a handler set previously with SetHandler.
func (ll *localListener) Start(cfg *conf.Conf) {
    if cfg.Local.TCPEnabled {
        ll.tcpStopChan = make(StopChan)
        go func() {
            // Create TCP socket listener
            tcp, err := net.ListenTCP("tcp", &cfg.Local.TCPAddr)
            if err != nil {
                // TODO: proper error handling
                log.Printf("Error: %v", err)
                return
            }
            ll.tcpListener = tcp

            listenOn(tcp, ll.tcpStopChan, ll.handler)

            // Stop the listener
            if err := tcp.Close(); err != nil {
                // TODO: error handling
            }
        }()
    }

    if cfg.Local.UnixEnabled {
        ll.unixStopChan = make(StopChan)
        go func() {
            // Create UNIX socket listener
            unix, err := net.ListenUnix("unix", &cfg.Local.UnixAddr)
            if err != nil {
                // TODO: proper error handling
                log.Printf("Error: %v", err)
                return
            }
            ll.unixListener = unix

            listenOn(unix, ll.unixStopChan, ll.handler)

            // Stop the listener
            if err := unix.Close(); err != nil {
                // TODO: error handling
            }
        }()
    }
}

// Stop terminates the listening process if it was started.
// It is no-op when called on stopped listener.
func (ll *localListener) Stop() {
    if ll.tcpStopChan != nil {
        ll.tcpStopChan.Stop()
        ll.tcpStopChan = nil
        ll.tcpListener.Close()
    }

    if ll.unixStopChan != nil {
        ll.unixStopChan.Stop()
        ll.unixStopChan = nil
        ll.unixListener.Close()
    }
}
