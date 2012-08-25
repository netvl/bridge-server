/**
 * Date: 23.08.12
 * Time: 23:43
 *
 * @author Vladimir Matveev
 */
package net

import (
    "net"
    "bridge/common"
    "bridge/common/msg"
)

type Communicator interface {
    Communicate(node common.Node, message msg.Message) (msg.Message, error)
}

type communicator struct {
    conn net.TCPConn
}

func (comm *communicator) Communicate(node common.Node, message *msg.Message) (*msg.Message, error) {
    return nil, nil
}

