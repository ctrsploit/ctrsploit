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
	return p.CheckTemplate(func() (bool, error) {
		p.Satisfied = selinux.IsEnabled() == p.expectedStatus
		return p.Satisfied, p.Err
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
