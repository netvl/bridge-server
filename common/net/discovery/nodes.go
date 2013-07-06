/**
 * Date: 06.07.13
 * Time: 22:09
 *
 * @author Vladimir Matveev
 */
package discovery

import (
    . "github.com/dpx-infinity/bridge-server/common"
    "net"
    "github.com/dpx-infinity/bridge-server/common/msg"
    "strings"
    "strconv"
)

type node struct {
    name  string
    addr  *net.IPAddr
    ports []int
}

func (n *node) Name() string {
    return n.name
}

func (n *node) Addr() *net.IPAddr {
    return n.addr
}

func (n *node) Ports() []int {
    return n.ports
}

type NodeReadingError struct {
    message string
    cause error
}

func (e *NodeReadingError) Error() string {
    return e.message
}

func (e *NodeReadingError) Cause() error {
    return e.cause
}

func newNodeReadingError(message string) *NodeReadingError {
    return newNodeReadingErrorWithCause(message, nil)
}

func newNodeReadingErrorWithCause(message string, cause error) *NodeReadingError {
    return &NodeReadingError{message: message, cause: cause}
}

func NodeFromMessage(message msg.Message) (Node, error) {
    nodeName := message.Header("Name")
    if nodeName == nil {
        return nil, newNodeReadingError("Node name is not specified in the message")
    }

    nodeAddr := message.Header("Address")
    if nodeAddr == nil {
        return nil, newNodeReadingError("Node address is not specified in the message")
    }
    addr, err := net.ResolveIPAddr("ip", nodeAddr)
    if err != nil {
        return nil, newNodeReadingErrorWithCause("Cannot resolve node address: " + nodeAddr, err)
    }

    nodePorts := message.Header("Ports")
    if nodePorts == nil {
        return nil, newNodeReadingError("Node ports are not specified in the message")
    }

    parts := strings.Split(*nodePorts)
    ports := make([]int, len(parts))
    for i, part := range parts {
        port, err := strconv.Atoi(part)
        if err != nil {
            return nil, newNodeReadingErrorWithCause("Cannot read port value: " + port, err)
        }
        if port < 1 || port > 65535 {
            return nil, newNodeReadingError("Port number is out of range: " + port)
        }
        ports[i] = port
    }

    return &node{name: nodeName, addr: addr, ports: ports}
}
