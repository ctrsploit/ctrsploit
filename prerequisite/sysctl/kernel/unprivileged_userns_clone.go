package kernel

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/sysctl"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type unprivilegedUsernsClone struct {
	prerequisite.BasePrerequisite
	enabled bool
}

func (p *unprivilegedUsernsClone) Check() (bool, error) {
	return p.CheckTemplate(func() {
		enabled, err := sysctl.UnprivilegedUsernsCloneEnabled()
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("getting sysctl: %w", err))
			return
		}
		p.Satisfied = enabled == p.enabled
		return
	})
}

var (
	UnprivilegedUsernsCloneEnabled = unprivilegedUsernsClone{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "unprivileged userns clone enabled",
			Info:   "kernel.unprivileged_userns_clone=1",
			ExeEnv: exeenv.InContainer | exeenv.InHost,
		},
		enabled: true,
	}
	UnprivilegedUsernsCloneDisabled = unprivilegedUsernsClone{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "unprivileged userns clone disabled",
			Info:   "kernel.unprivileged_userns_clone=0",
			ExeEnv: exeenv.InContainer | exeenv.InHost,
		},
		enabled: false,
	}
)
