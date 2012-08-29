/**
 * Date: 29.08.12
 * Time: 23:36
 *
 * @author Vladimir Matveev
 */
package listener

import (
    "bridge/common/conf"
)

type remoteListener struct {
    handler Handler
    tcpStopChan StopChan
    udpStopChan StopChan
}

func NewRemoteListener() Listener {
    return new(remoteListener)
}

func (rl *remoteListener) SetHandler(handler Handler) {

}

func (rl *remoteListener) Start(cfg *conf.Conf) {

}

func (rl *remoteListener) Stop() {

}
