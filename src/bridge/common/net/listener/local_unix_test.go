// +build !windows

/**
 * Date: 30.08.12
 * Time: 22:14
 *
 * @author Vladimir Matveev
 */
package listener_test

import (
    . "launchpad.net/gocheck"
    "net"
    "bridge/common/conf"
    "bridge/common/net/listener"
    "time"
)

func (_ *ListenerSuite) TestListenUnix(c *C) {
    paddr, _ := net.ResolveUnixAddr("unix", "/tmp/go_test_socket")
    addr := *paddr
    cfg := conf.Conf{
        Remote: conf.RemoteListenerConf{},
        Local: conf.LocalListenerConf{
            UnixEnabled: true,
            UnixAddr: addr,
        },
    }

    ll := listener.NewLocalListener()
    ll.SetHandler(echoHandler)
    ll.Start(&cfg)
    time.Sleep(50 * time.Millisecond)

    cc, err := net.DialUnix("unix", nil, &addr)
    c.Assert(err, IsNil)

    testEchoConnection(c, cc)

    cc.Close()

    ll.Stop()
}

