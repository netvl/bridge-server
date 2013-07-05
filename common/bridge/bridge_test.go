/**
 * Date: 31.08.12
 * Time: 0:58
 *
 * @author Vladimir Matveev
 */
package bridge_test

import (
    "github.com/dpx-infinity/bridge-server/common/bridge"
    "github.com/dpx-infinity/bridge-server/common/conf"
    "github.com/dpx-infinity/bridge-server/common/msg"
    "github.com/dpx-infinity/bridge-server/common/plugins"
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
    cfg := &conf.Conf{
        Listeners: map[string]*conf.ListenerConf{
            "local": &conf.ListenerConf{
                Name: "local",
                Ports: map[conf.PortType]*conf.PortConf{
                    conf.PortTypeTCP4: &conf.PortConf{
                        Type: conf.PortTypeTCP4,
                        Addr: addr,
                    },
                },
            },
        },
        Plugins: map[string]*conf.PluginConf{
            "echo-plugin": &conf.PluginConf{
                Name:      "echo-plugin",
                Plugin:    "echo",
                Listeners: []string{"local"},
                Mediators: []string{},
                Options:   map[string][]string{"prefix": []string{"Echo"}},
            },
        },
    }

    br := bridge.New(cfg)
    br.AddPlugin("echo-plugin", new(plugins.EchoPlugin))
    br.Start()
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
    m.SetBodyPart("bp1", bps)

    err := msg.Serialize(m, cc)
    printCause(c, err)
    c.Assert(err, IsNil)

    m2, err := msg.Deserialize(cc)
    c.Assert(err, IsNil)
    c.Assert(m2.GetName(), Equals, m.GetName())
    c.Assert(m2.Header("hdr1"), Equals, m.Header("hdr1"))
    c.Assert(m2.Header("hdr2"), Equals, m.Header("hdr2"))
    c.Assert(m2.BodyPart("bp1"), DeepEquals, bps)
}
