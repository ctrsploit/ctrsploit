package selinux

import (
	"github.com/ctrsploit/ctrsploit/pkg/selinux"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type unconfined struct {
	prerequisite.BasePrerequisite
}

func (p *unconfined) Check() (bool, error) {
	return p.CheckTemplate(func() {
		if selinux.IsEnabled() {
			if selinux.IsSelinuxPrivileged() {
				p.Satisfied = true
				return
			}
		} else {
			p.Satisfied = true
			return
		}
	})
}

var (
	Unconfined = unconfined{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "selinux unconfined",
			Info:   "selinux disabled or use privileged profile",
			ExeEnv: exeenv.InContainer,
		},
	}
)
