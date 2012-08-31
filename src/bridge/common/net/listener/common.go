/**
 * Date: 29.08.12
 * Time: 23:36
 *
 * @author Vladimir Matveev
 */
package listener

import (
    "bridge/common/conf"
    "log"
    "net"
    "syscall"
)

// StopChan represents a channel used to tell to some process that it should stop its work.
type StopChan chan interface{}

// Stopped checks whether the given StopChan received a request to stop.
func (ch StopChan) Stopped() bool {
    select {
    case _ = <-ch:
        return true
    default:
        // Still not stopped
    }
    return false
}

// Wait waits until a stop request will have been sent to the given StopChan.
func (ch StopChan) Wait() {
    _ = <-ch
}

// Stop sends a stop request to the given StopChan.
func (ch StopChan) Stop() {
    select {
    case ch <- nil:
    default:
    }
    close(ch)
}

// Handler is a function which is able to handle standard connection.
// It is supposed that the handler itself does not close the connection.
type Handler func (net.Conn)

// Listener is an entity which can start listening for external connections
// and transfer control to provided handler when a connection was received.
// It is also should be possible to stop listener.
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
            // TODO: proper error handling
            // for now, just log and continue
            if operr, ok := err.(*net.OpError); ok {
                if operr.Err == syscall.ECONNABORTED {  // The listener has been closed externally
                    break
                }
            }
            log.Printf("Error accepting connection: %v", err)
            continue
        }

        // Handle the connection
        if handler != nil {
            handler(conn)
        }

        // Close the connection
        if err := conn.Close(); err != nil {
            // TODO: proper error handling
            log.Printf("Error closing connection: %v", err)
        }
    }
}


