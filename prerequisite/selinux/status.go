package selinux

import (
	"github.com/ctrsploit/ctrsploit/pkg/selinux"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type status struct {
	prerequisite.BasePrerequisite
	expectedStatus bool
}

func (p *status) Check() (bool, error) {
	return p.CheckTemplate(func() {
		p.Satisfied = selinux.IsEnabled() == p.expectedStatus
		return
	})
}

var (
	Disabled = status{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "selinux disabled",
			Info:   "",
			ExeEnv: exeenv.InContainer,
		},
		expectedStatus: false,
	}
)
