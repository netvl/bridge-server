/**
 * Date: 23.08.12
 * Time: 23:43
 *
 * @author Vladimir Matveev
 */
package comm

import (
    "net"
    "bridge/common/msg"
)

type Node struct {
    name string
    addr net.IPAddr
}

func (n *Node) Name() string {
    return n.name
}

func (n *Node) Addr() net.IPAddr {
    return n.addr
}

type Communicator struct {
    conn net.TCPConn
}

func NewCommunicator() *Communicator {
    return new(Communicator)
}

func (comm *Communicator) Communicate(node Node, msg *msg.Message) (*msg.Message, error) {
    return nil, nil
}


