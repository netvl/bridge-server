/**
 * Date: 29.08.12
 * Time: 20:50
 *
 * @author Vladimir Matveev
 */
package listener_test

import (
    . "launchpad.net/gocheck"
    "bridge/common/net/listener"
    "testing"
    "bridge/common/conf"
    "net"
    "bufio"
    "time"
)

func Test(t *testing.T) {
    TestingT(t)
}

type ListenerSuite struct {}
var _ = Suite(&ListenerSuite{})

func echoHandler(c net.Conn) {
    b := bufio.NewReader(c)
    data, _ := b.ReadBytes(0)
    c.Write(data)
}

func (_ *ListenerSuite) TestListenTCP(c *C) {
    paddr, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:12345")
    addr := *paddr
    cfg := conf.Conf{
        Remote: conf.RemoteListenerConf{},
        Local: conf.LocalListenerConf{
            TCPEnabled: true,
            TCPAddr: addr,
        },
    }

    ll := listener.NewLocalListener()
    ll.SetHandler(echoHandler)
    ll.Start(&cfg)
    time.Sleep(50 * time.Millisecond)

    cc, err := net.DialTCP("tcp4", nil, &addr)
    c.Assert(err, IsNil)

    testEchoConnection(c, cc)

    cc.Close()

    ll.Stop()
}

func testEchoConnection(c *C, cc net.Conn) {
    tpl := []byte{1, 2, 3, 4, 5, 0}
    n, err := cc.Write(tpl)
    c.Assert(n, Equals, len(tpl))
    c.Assert(err, IsNil)

    buf := make([]byte, len(tpl))
    n, err = cc.Read(buf)
    c.Assert(n, Equals, len(tpl))
    c.Assert(buf, DeepEquals, tpl)
    c.Assert(err, IsNil)
}

