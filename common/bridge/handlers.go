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

func (b *bridge) makePluginsHandler(listener string) Handler {
    return func(conn net.Conn) {
        m, err := msg.Deserialize(conn)
        if err != nil {
            log.Printf("Bridge local plugin handler message deserialization error: %t", err)
            return
        }

        for _, p := range b.plugins {
            if p.SupportsMessage(m.GetName()) {
                rm := p.HandleMessage(m, b)
                if err := msg.Serialize(rm, conn); err != nil {
                    log.Printf("Bridge local plugin handler message serialization error: %v", err)
                    break
                }
                break
            }
        }
    }
}
