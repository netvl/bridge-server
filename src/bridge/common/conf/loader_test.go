package conf_test

import (
    "bridge/common/conf"
    "testing"
)

func TestLoading(t *testing.T) {
    conf.LoadConfig("/home/dpx-infinity/dev/projects/bridge/bridge-server-common/resources/example.conf")
}
