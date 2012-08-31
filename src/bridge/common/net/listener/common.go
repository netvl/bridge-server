/**
 * Date: 29.08.12
 * Time: 23:36
 *
 * @author Vladimir Matveev
 */
package listener

import (
    "net"
    "bridge/common/conf"
)

type StopChan chan interface{}

func (ch StopChan) Stopped() bool {
    select {
    case _ = <-ch:
        return true
    default:
        // Still not stopped
    }
    return false
}

func (ch StopChan) Wait() {
    _ = <-ch
}

func (ch StopChan) Stop() {
    select {
    case ch <- nil:
    default:
    }
    close(ch)
}

type Handler func (net.Conn)

type Listener interface {
    SetHandler(handler Handler)
    Start(cfg *conf.Conf)
    Stop()
}

func listenOn(listener net.Listener, stopChan StopChan, handler Handler) {
    for {
        // Check if we have to exit
        if stopChan.Stopped() {
            break
        }

        // Accept a connection
        conn, err := listener.Accept()
        if err != nil {
            // TODO: error handling
            // for now, just continue
            continue
        }

        // Handle the connection
        if handler != nil {
            handler(conn)
        }

        // Close the connection
        if err := conn.Close(); err != nil {
            // TODO: error handling
        }
    }
}


