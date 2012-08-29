/**
 * Date: 29.08.12
 * Time: 23:56
 *
 * @author Vladimir Matveev
 */
package bridge

import (
    "net"
    "bridge/common/net/comm"
    "bridge/common/net/listener"
)

func makeLocalPluginsHandler(plugins map[string]LocalPlugin, c *comm.Communicator) listener.Handler {
    return func(conn net.Conn) {

    }
}

func makeRemotePluginsHandler(plugins map[string]RemotePlugin, c *comm.Communicator) listener.Handler {
    return func(conn net.Conn) {

    }
}


