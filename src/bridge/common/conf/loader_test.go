package conf_test

import (
    "testing"
    "bridge/common/conf"
)

func TestLoading(t *testing.T) {
    conf.LoadConfig("/home/dpx-infinity/dev/projects/bridge/bridge-server-common/resources/example.conf")
}
