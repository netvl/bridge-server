/**
 * Date: 27.08.12
 * Time: 21:42
 *
 * @author Vladimir Matveev
 */
package conf

import (
    "net"
)

type RemoteListenerConf struct {
    TCPEnabled bool
    TCPAddr net.TCPAddr

    UDPEnabled bool
    UDPAddr net.UDPAddr
}

type LocalListenerConf struct {
    TCPEnabled bool
    TCPAddr net.TCPAddr

    UnixEnabled bool
    UnixAddr net.UnixAddr
}

type Conf struct {
    Remote RemoteListenerConf
    Local LocalListenerConf
}
