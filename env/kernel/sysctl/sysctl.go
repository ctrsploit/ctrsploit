package sysctl

import (
	"errors"

	"github.com/ctrsploit/ctrsploit/pkg/sysctl"
	spec "github.com/ctrsploit/sploit-spec/pkg/env/container/kernel/sysctl"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

const CommandName = "sysctl"

func assign[T any](dest *T, fn func() (T, error)) func() error {
	return func() error {
		var err error
		*dest, err = fn()
		return err
	}
}

func runTasks(tasks ...func() error) error {
	var errs []error
	for _, task := range tasks {
		if err := task(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func Sysctl() (spec.Sysctl, error) {
	var s spec.Sysctl
	err := runTasks(
		assign(&s.RouteLocalNet, sysctl.RouteLocalNetEnabled),
		assign(&s.MaxUserNamespaces, sysctl.MaxUserNamespaces),
		assign(&s.UnprivilegedUsernsClone, sysctl.UnprivilegedUsernsCloneEnabled),
		assign(&s.PidMax, sysctl.PidMax),
		assign(&s.ThreadsMax, sysctl.ThreadsMax),
	)
	if err != nil {
		awesome_error.CheckWarning(err)
	}
	return s, err
}
