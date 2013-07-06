/**
 * Date: 05.07.13
 * Time: 21:40
 *
 * @author Vladimir Matveev
 */
package ports

import (
    . "github.com/dpx-infinity/bridge-server/common"
    "github.com/dpx-infinity/bridge-server/common/util"
    "github.com/dpx-infinity/bridge-server/common/msg"
    "net"
    "runtime"
    "log"
)

type TcpPort struct {
    sockets  map[string]ChanPair
    net      string
    subnet   *net.IPNet
    addr     *net.TCPAddr
    stopChan util.StopChan
}

// Returns new TCP port defined by CIDR-qualified address addr and TCP port. Type of the address
// (IPv4 or IPv6) is defined by net parameter (either tcp4 or tcp6). Connected sockets are specified by
// sockets parameter.
// This port will listen on the specified address and send messages without destination to broadcast address
// in the current network (determined from the subnet address).
func NewTcpPort(sockets map[string]ChanPair, net string, addr string, port int) *TcpPort {
    ip, subnet, err := net.ParseCIDR(addr)
    if err != nil {
        // TODO: add error result?
        log.Println("Cannot parse CIDR address:", err)
        return nil
    }
    return &TcpPort{sockets: sockets, net: net, subnet: subnet, addr: &net.TCPAddr{IP: ip, Port: port}}
}

func (port *TcpPort) Start() error {
    listener, err := net.ListenTCP(port.net, port.addr)
    if err != nil {
        return err
    }

    // Listening goroutine
    go func() {
        for {
            if (port.stopChan.Stopped()) {
                break
            }

            conn, err := listener.AcceptTCP()
            if err != nil {
                log.Println("Error accepting TCP connection:", err)
                continue
            }

            message, err := msg.Deserialize(conn)
            if err != nil {
                log.Println("Error deserializing message:", err)
                continue
            }

            for _, socket := range port.sockets {
                socket.Sink() <- message
            }
        }
    }()

    // Sending goroutine
    go func() {
        for {
            if (port.stopChan.Stopped()) {
                break
            }

            for _, socket := range port.sockets {
                select {
                case msg := <-socket.Source():
                    if err := port.sendMessageOverTCP(msg); err != nil {
                        // TODO: handle error
                    }
                default:
                    runtime.Gosched()
                }
            }
        }
    }()

    return nil
}

func (port *TcpPort) sendMessageOverTCP(message *msg.Message) error {
    var tcpaddr *net.TCPAddr

    if dest := message.Header("Destination"); dest == nil { // Broadcast
        broadcast := util.OrSlices([]byte(port.subnet.IP), util.NotSlice([]byte(port.subnet.Mask)))

    } else {
        net, addr, err := util.SplitDestination(dest)
        if err != nil {
            return err
        }

        tcpaddr, err = net.ResolveTCPAddr(net, addr)
        if err != nil {
            return err
        }
    }

    conn, err := net.DialTCP(port.net, nil, tcpaddr)
    if err != nil {
        return err
    }
    defer conn.Close()

    if err := msg.Serialize(message, conn); err != nil {
        return err
    }

    return nil
}

func (port *TcpPort) Stop() {
    port.stopChan.Stop()
}

