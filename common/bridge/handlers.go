/**
 * Date: 29.08.12
 * Time: 23:56
 *
 * @author Vladimir Matveev
 */
package bridge

import (
    . "github.com/dpx-infinity/bridge-server/common"
    "github.com/dpx-infinity/bridge-server/common/msg"
    "log"
    "net"
)

func makePluginsHandler(plugins map[string]Plugin, c Communicator) Handler {
    return func(conn net.Conn) {
        // Load message name and headers
        m, err := msg.DeserializeMessageName(conn)
        if err != nil {
            // TODO: proper error handling
            log.Printf("Bridge local plugin handler message name deserialization error: %t", err)
            return
        }
        if err := msg.DeserializeMessageHeaders(conn, m); err != nil {
            log.Printf("Bridge local plugin handler message headers deserialization error: %v", err)
            return
        }

        for _, p := range plugins {
            doHandle := false
            switch conn.(type) {
            case *net.TCPConn:
                if p.PluginTypes()[PluginTypeTCP] {
                    doHandle = true
                }
            case *net.UDPConn:
                if p.PluginTypes()[PluginTypeUDP] {
                    doHandle = true
                }
            case *net.UnixConn:
                if p.PluginTypes()[PluginTypeUnix] {
                    doHandle = true
                }
            }

            if doHandle && p.SupportsMessage(m.GetName()) {
                if err := msg.DeserializeMessageBodyParts(conn, m, p.DeserializeHook()); err != nil {
                    log.Printf("Bridge local plugin handler message body parts deserialization error: %v", err)
                    break
                }
                rm := p.HandleMessage(m, nil)
                if err := msg.Serialize(rm, conn); err != nil {
                    log.Printf("Bridge local plugin handler message serialization error: %v", err)
                    break
                }
                break
            }
        }
    }
}
