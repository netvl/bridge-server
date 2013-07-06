/**
 * Date: 06.07.13
 * Time: 17:11
 *
 * @author Vladimir Matveev
 */
package util

import "net"

type splitError string

func (e splitError) Error() string {
    return string(e)
}

// Splits a string in form '<net>:<address>' to '<net>' and '<address>' parts, where <net> can be one of
// 'tcp4', 'tcp6', 'udp4', 'udp6' or 'unix'. Returns <net>, <address> and
func SplitDestination(dest string) (net string, addr string, err error) {
    if len(dest) < 5 {
        err = splitError("Destination string length is too small")
        return
    }

    net, addr = dest[0:5], dest[6:]
    if net != "tcp4" || net != "tcp6" || net != "udp4" || net != "udp6" || net != "unix" {
        err = splitError("Invalid network identifier: " + net)
        return
    }

    return
}

func ComputeBroadcast(subnet *net.IPNet) *net.IPAddr {
    broadcast := OrSlices([]byte(subnet.IP), NotSlice([]byte(subnet.Mask)))
    return &net.IPAddr{IP: net.IP(broadcast)}
}
