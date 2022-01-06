package auto

import (
	"github.com/ctrsploit/ctrsploit/env/apparmor"
	"github.com/ctrsploit/ctrsploit/env/cgroups"
	"github.com/ctrsploit/ctrsploit/env/seccomp"
)

func Auto() (err error) {
	err = seccomp.Seccomp()
	if err != nil {
		return
	}
	err = apparmor.Apparmor()
	if err != nil {
		return
	}
	err = cgroups.Version()
	if err != nil {
		return
	}
	err = cgroups.ListSubsystems()
	if err != nil {
		return
	}
	return
}
