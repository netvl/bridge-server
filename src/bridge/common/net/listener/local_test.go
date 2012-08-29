/**
 * Date: 29.08.12
 * Time: 20:50
 *
 * @author Vladimir Matveev
 */
package local_test

import (
    . "launchpad.net/gocheck"
    "testing"
)

func Test(t *testing.T) {
    TestingT(t)
}

type ListenerSuite struct {}
var _ = Suite(&ListenerSuite{})

func (_ *ListenerSuite) TestListenTCP(c *C) {

}

