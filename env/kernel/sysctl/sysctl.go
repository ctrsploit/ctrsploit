package sysctl

import (
	"github.com/ctrsploit/ctrsploit/pkg/sysctl"
	"github.com/ctrsploit/sploit-spec/pkg/env/container"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

const CommandName = "sysctl"

func Sysctl() (s container.Sysctl, err error) {
	s.RouteLocalNet, err = sysctl.RouteLocalNetEnabled()
	if err != nil {
		awesome_error.CheckWarning(err)
		return
	}
	s.MaxUserNamespaces, err = sysctl.MaxUserNamespaces()
	if err != nil {
		awesome_error.CheckWarning(err)
		return
	}
	s.UnprivilegedUsernsClone, err = sysctl.UnprivilegedUsernsCloneEnabled()
	if err != nil {
		awesome_error.CheckWarning(err)
		return
	}
	return
}
