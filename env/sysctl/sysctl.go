package sysctl

import (
	"github.com/ctrsploit/ctrsploit/pkg/sysctl"
	"github.com/ctrsploit/sploit-spec/pkg/env/container"
)

const CommandName = "sysctl"

func Sysctl() (container.Sysctl, error) {
	routeLocalNetEnabled, err := sysctl.RouteLocalNetEnabled()
	if err != nil {
		return container.Sysctl{}, err
	}
	return container.Sysctl{
		ProcSysNetIpv4ConfAllRouteLocalNet: routeLocalNetEnabled,
	}, nil
}
