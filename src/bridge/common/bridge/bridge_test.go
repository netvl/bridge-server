/**
 * Date: 31.08.12
 * Time: 0:58
 *
 * @author Vladimir Matveev
 */
package bridge_test

import (
    "bridge/common/bridge"
    "bridge/common/conf"
    "bridge/common/msg"
    "bridge/common/plugins"
    . "launchpad.net/gocheck"
    "net"
    "testing"
    "time"
)

func Test(t *testing.T) {
    TestingT(t)
}

type BridgeSuite struct{}

var _ = Suite(BridgeSuite{})

func (_ BridgeSuite) TestLocalHandling(c *C) {
    addr, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:12345")
    cfg := conf.Conf{
        Remote: conf.RemoteListenerConf{},
        Local: conf.LocalListenerConf{
            TCPEnabled: true,
            TCPAddr:    *addr,
        },
    }

    br := bridge.New()
    br.AddLocalPlugin("echo-plugin", &plugins.EchoPlugin{})
    br.Start(&cfg)
    time.Sleep(50 * time.Millisecond)

    cc, err := net.DialTCP("tcp4", nil, addr)
    c.Assert(err, IsNil)

    testEchoPlugin(cc, c)

    cc.Close()

    br.Stop()
}

func printCause(c *C, err error) {
    if err != nil {
        c.Log(err.(*msg.SerializeError).Cause())
    }
}

func testEchoPlugin(cc net.Conn, c *C) {
    bps := []byte{1, 2, 3, 4, 5}
    m := msg.CreateWithName("test-name")
    m.SetHeader("hdr1", "val1")
    m.SetHeader("hdr2", "val2")
    m.SetBodyPart("bp1", msg.BodyPartFromSlice(bps))

    err := msg.Serialize(m, cc)
    printCause(c, err)
    c.Assert(err, IsNil)

    rm, err := msg.DeserializeMessageName(cc)
    printCause(c, err)
    c.Assert(err, IsNil)
    c.Assert(rm.GetName(), Equals, m.GetName())

    err = msg.DeserializeMessageHeaders(cc, rm)
    printCause(c, err)
    c.Assert(err, IsNil)
    c.Assert(rm.GetHeader("hdr1"), Equals, m.GetHeader("hdr1"))
    c.Assert(rm.GetHeader("hdr2"), Equals, m.GetHeader("hdr2"))

    err = msg.DeserializeMessageBodyParts(cc, rm, msg.EmptyHook)
    c.Assert(err, IsNil)
    printCause(c, err)

    buf := make([]byte, rm.GetBodyPart("bp1").Size())
    n, err := rm.GetBodyPart("bp1").Read(buf)
    c.Assert(err, IsNil)
    c.Assert(n, Equals, len(bps))
    c.Assert(buf, DeepEquals, bps)
}
