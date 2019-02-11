package nettool

import (
	"net"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

var link netlink.Link

type EIPRoute struct {
	NextIP net.IP
	EIP    net.IPNet
}

func init() {
	link, _ = netlink.LinkByName("eth0")
}

func NewEIPRoute(eip, nextip string, mask int) *EIPRoute {
	EIP := net.IPNet{
		IP:   net.ParseIP(eip),
		Mask: net.CIDRMask(mask, 32),
	}
	n := net.ParseIP(nextip)
	return &EIPRoute{
		NextIP: n,
		EIP:    EIP,
	}
}
func addLocalRule(ip *net.IPNet, mask int) error {
	rule := netlink.NewRule()
	rule.Table = 101
	rule.Dst = ip
	rule.Mask = mask
	return netlink.RuleAdd(rule)
}

func (e *EIPRoute) ToNetlinkRoute() *netlink.Route {
	return &netlink.Route{
		LinkIndex: link.Attrs().Index,
		Src:       e.NextIP,
		Type:      unix.RTN_NAT,
		Dst:       &e.EIP,
	}
}
func (e *EIPRoute) Add() error {
	return netlink.RouteReplace(e.ToNetlinkRoute())
}

func (e *EIPRoute) Delete() error {
	return netlink.RouteDel(e.ToNetlinkRoute())
}

func (e *EIPRoute) IsExist() (bool, error) {
	routes, err := netlink.RouteList(link, netlink.FAMILY_V4)
	if err != nil {
		return false, err
	}
	r_self := e.ToNetlinkRoute()
	for _, r := range routes {
		if r.Equal(*r_self) {
			return true, nil
		}
	}
	return false, nil
}
