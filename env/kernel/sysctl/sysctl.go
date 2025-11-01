package sysctl

import (
	"github.com/ctrsploit/ctrsploit/pkg/sysctl"
	spec "github.com/ctrsploit/sploit-spec/pkg/env/container/kernel/sysctl"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

const CommandName = "sysctl"

func Sysctl() (s spec.Sysctl, err error) {
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
		return
	}
	s.PidMax, err = sysctl.PidMax()
	if err != nil {
		awesome_error.CheckWarning(err)
		return
	}
	return
}
