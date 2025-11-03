package net

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/sysctl"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type RouteLocalnet struct {
	prerequisite.BasePrerequisite
	Expected bool
}

func (p *RouteLocalnet) Check() (bool, error) {
	return p.CheckTemplate(func() {
		enabled, err := sysctl.RouteLocalNetEnabled()
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("getting sysctl: %w", err))
			return
		}
		p.Satisfied = enabled == p.Expected
		return
	})
}

var (
	RouteLocalNetEnabled = RouteLocalnet{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "route_localnet enabled",
			Info:   "net.ipv4.conf.all.route_localnet = 1",
			ExeEnv: exeenv.InHost | exeenv.InContainer,
		},
		Expected: true,
	}
)
