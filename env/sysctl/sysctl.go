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
	maxUserNamespaces, err := sysctl.MaxUserNamespaces()
	if err != nil {
		return container.Sysctl{}, err
	}
	unprivilegedUsernsClone, err := sysctl.UnprivilegedUsernsCloneEnabled()
	if err != nil {
		return container.Sysctl{}, err
	}
	return container.Sysctl{
		Net: container.Net{
			RouteLocalNet: routeLocalNetEnabled,
		},
		User: container.User{
			MaxUserNamespaces: maxUserNamespaces,
		},
		KernelSysctl: container.KernelSysctl{
			UnprivilegedUsernsClone: unprivilegedUsernsClone,
		},
	}, nil
}
