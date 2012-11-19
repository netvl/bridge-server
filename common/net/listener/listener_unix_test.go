// +build !windows

/**
 * Date: 30.08.12
 * Time: 22:14
 *
 * @author Vladimir Matveev
 */
package listener_test

import (
    "github.com/dpx-infinity/bridge-server/common/conf"
    "github.com/dpx-infinity/bridge-server/common/net/listener"
    . "launchpad.net/gocheck"
    "net"
    "time"
)

func (_ *ListenerSuite) TestListenUnix(c *C) {
    addr, _ := net.ResolveUnixAddr("unix", "/tmp/go_test_socket")
    cfg := &conf.ListenerConf{
        Name: "local",
        Ports: []*conf.PortConf{
            &conf.PortConf{
                Type: "unix",
                Addr: addr,
            },
        },
    }

    ll := listener.NewListener(cfg)
    ll.SetHandler(echoHandler)
    ll.Start()
    time.Sleep(50 * time.Millisecond)

    cc, err := net.DialUnix("unix", nil, addr)
    c.Assert(err, IsNil)

    testEchoConnection(c, cc)

    cc.Close()

    ll.Stop()
}
