/**
 * Date: 14.01.13
 * Time: 3:08
 *
 * @author Vladimir Matveev
 */
package util_test

import (
    "testing"
    . "launchpad.net/gocheck"
    "github.com/dpx-infinity/bridge-server/common/util"
)


func Test(t *testing.T) {
    TestingT(t)
}

type MiscSuite struct{}

var _ = Suite(MiscSuite{})

func (_ MiscSuite) TestStringPatternMatching(c *C) {
    p1 := util.StringPattern("a", nil, "b", nil, "c")
    c.Assert(util.MatchStringSlicePattern(p1, []string{"a", "kshs", "b", "skjsdf", "c"}), Equals, true)
    c.Assert(util.MatchStringSlicePattern(p1, []string{"a", "sjkwkwk", "lslkj", "sfsff", "c"}), Equals, false)
    c.Assert(util.MatchStringSlicePattern(p1, []string{"a", "sjkwkwk", "b", "sfsff"}), Equals, false)
    c.Assert(util.MatchStringSlicePattern(p1, []string{"a", "sjkwkwk", "b", "sfsff", "c", "lsksks"}), Equals, false)

    p2 := util.StringPattern("a", nil, "b", nil, "c", nil)
    c.Assert(util.MatchStringSlicePattern(p2, []string{"a", "ddff", "b", "dkdkd", "c"}), Equals, true)
    c.Assert(util.MatchStringSlicePattern(p2, []string{"a", "ddff", "b", "dkdkd", "c", "akkakkd"}), Equals, true)
    c.Assert(util.MatchStringSlicePattern(p2, []string{"a", "ddff", "b", "dkdkd", "c", "akkakkd", "skskks"}),
        Equals, true)
    c.Assert(util.MatchStringSlicePattern(p2, []string{"a", "ddff", "b", "dkdkd"}), Equals, false)
}

