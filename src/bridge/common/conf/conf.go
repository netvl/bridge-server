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

type Conf struct {
    TCPListenAddr net.TCPAddr
    UDPListenAddr net.UDPAddr
}
