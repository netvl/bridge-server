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

func (ch StopChan) Stop() {
    ch <- nil
    close(ch)
}

type Handler func (net.Conn)

type Listener interface {
    SetHandler(handler Handler)
    Start(cfg *conf.Conf)
    Stop()
}


