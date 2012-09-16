/**
 * Date: 23.08.12
 * Time: 23:43
 *
 * @author Vladimir Matveev
 */
package comm

import (
    "bridge/common"
    "bridge/common/msg"
    "net"
)

type Node struct {
    name string
    addr net.IPAddr
}

func NewNode(name string, addr net.IPAddr) common.Node {
    return &Node{name, addr}
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

func NewComm() common.Communicator {
    return new(Communicator)
}

func (comm *Communicator) Send(node common.Node, msg *msg.Message) error {
    return nil
}

func (comm *Communicator) Receive(node common.Node) (*msg.Message, error) {
    return nil, nil
}
