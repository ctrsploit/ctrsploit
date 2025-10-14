package sysctl

import (
	"fmt"
	"os"
	"strings"
)

const (
	PathRouteLocalNet = "/proc/sys/net/ipv4/conf/all/route_localnet"
)

/*
RouteLocalNetEnabled
route_localnet - BOOLEAN

	Do not consider loopback addresses as martian source or destination
	while routing. This enables the use of 127/8 for local routing purposes.
	default FALSE

https://github.com/kubernetes/kubernetes/issues/90259
*/
func RouteLocalNetEnabled() (bool, error) {
	content, err := os.ReadFile(PathRouteLocalNet)
	if err != nil {
		return false, fmt.Errorf("failed to read %s: %v", PathRouteLocalNet, err)
	}
	return strings.TrimSpace(string(content)) == "1", nil
}
