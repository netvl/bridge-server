package conf_test

import (
    "bridge/common/conf"
    "testing"
    "fmt"
)

func TestLoading(t *testing.T) {
    cfg := conf.LoadConfig("/home/dpx-infinity/dev/projects/bridge/bridge-server-common/resources/example.conf")
    if cfg == nil {
        return
    }
    for ln, lc := range cfg.Listeners {
        fmt.Printf("Listener %s:\n", ln)
        for _, pc := range lc.Ports {
            fmt.Printf("\tPort type %s, port address %v\n", pc.Type, pc.Addr)
        }
    }
}
