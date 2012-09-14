package conf_test

import (
    "bridge/common/conf"
    "testing"
    "log"
    "os"
    "strings"
    . "launchpad.net/gocheck"
    "bytes"
    "net"
)

var testConfListeners = `conf listeners {
    conf local {
        conf tcp4 {
            set host localhost
            set port 16128
        }
        conf tcp6 {
            set host "[::1]"
            set port 16128
        }
        conf unix {
            #set path /run/bridge/local.sock
        }
    }
    conf remote {
        conf tcp4 {
            set host 1.2.3.4
            set port 16129
        }
        conf udp4 {
            set host *
            set port 16129
        }
        conf invalid {

        }
    }
}`

type LoaderSuite struct {}
var _ = Suite(&LoaderSuite{})

func TestLoaderSuite(t *testing.T) {
    TestingT(t)
}

func (_ *LoaderSuite) TestLoadingFromFile(c *C) {
    gopath := os.Getenv("GOPATH")
    if gopath == "" {
        c.Fatal("GOPATH variable is not set!")
    }
    var exampleConfPath string
    for _, path := range strings.Split(gopath, ":") {
        confPath := path + "/resources/example.conf"
        if _, err := os.Stat(confPath); err == nil {
            exampleConfPath = confPath
            break
        }
    }
    if exampleConfPath == "" {
        c.Fatal("Unable to find example config!")
    }

    cfg, err := conf.LoadConfigFromFile(exampleConfPath)

    if err != nil {
        if errs, ok := err.(*conf.ConfigErrors); ok {
            log.Println(errs)
        }
    }
    for _, lc := range cfg.Listeners {
//        fmt.Printf("Listener %s:\n", ln)
        for _, _ = range lc.Ports {
//            fmt.Printf("\tPort type %s, port address %v\n", pc.Type, pc.Addr)
        }
    }
}

func (_ *LoaderSuite) TestListeners(c *C) {
    r := bytes.NewReader([]byte(testConfListeners))

    cfg, err := conf.LoadConfigFromReader(r)
    c.Assert(err, NotNil)
    errs, ok := err.(*conf.ConfigErrors)
    c.Assert(ok, Equals, true)
    c.Assert(len(errs.Errors()), Equals, 2)
    c.Log(errs)

    c.Assert(len(cfg.Listeners), Equals, 2)

    local, ok := cfg.Listeners["local"]
    c.Assert(ok, Equals, true)
    c.Assert(len(local.Ports), Equals, 2)

    var local_tcp4ok, local_tcp6ok bool
    for ptype, pconf := range local.Ports {
        if ptype == "tcp4" {
            local_tcp4ok = true
            addr, _ := net.ResolveTCPAddr("tcp4", "localhost:16128")
            c.Assert(pconf.Addr, DeepEquals, addr)
        } else if ptype == "tcp6" {
            local_tcp6ok = true
            addr, _ := net.ResolveTCPAddr("tcp6", "[::1]:16128")
            c.Assert(pconf.Addr, DeepEquals, addr)
        } else {
            c.Fatal("Invalid port type:", ptype)
        }
    }
    c.Assert(local_tcp4ok && local_tcp6ok, Equals, true)

    remote, ok := cfg.Listeners["remote"]
    c.Assert(ok, Equals, true)
    c.Assert(len(local.Ports), Equals, 2)
    var remote_tcp4ok, remote_udp4ok bool
    for ptype, pconf := range remote.Ports {
        if ptype == "tcp4" {
            remote_tcp4ok = true
            addr, _ := net.ResolveTCPAddr("tcp4", "1.2.3.4:16129")
            c.Assert(pconf.Addr, DeepEquals, addr)
        } else if ptype == "udp4" {
            remote_udp4ok = true
            addr, _ := net.ResolveUDPAddr("udp4", "0.0.0.0:16129")
            c.Assert(pconf.Addr, DeepEquals, addr)
        }
    }
    c.Assert(remote_tcp4ok && remote_udp4ok, Equals, true)
}
