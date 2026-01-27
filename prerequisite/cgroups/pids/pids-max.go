package pids

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/cgroup/pids"
	"github.com/ctrsploit/ctrsploit/pkg/sysctl"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type pidsLimited struct {
	prerequisite.BasePrerequisite
	limited bool
}

func (p *pidsLimited) Check() (bool, error) {
	return p.CheckTemplate(func() {
		pidsMax, err := pids.GetMax()
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("getting max pids: %w", err))
			return
		}
		limited := pidsMax >= 0
		sysctlPidMax, err := sysctl.PidMax()
		if err == nil {
			if pidsMax > int64(sysctlPidMax) {
				limited = false
			}
		}
		sysctlThreadsMax, err := sysctl.ThreadsMax()
		if err == nil {
			if pidsMax > int64(sysctlThreadsMax) {
				limited = false
			}
		}
		p.Satisfied = limited == p.limited
		return
	})
}

var (
	UnlimitedPidsMax = pidsLimited{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "pids.max unlimited",
			Info:   "pids.max != max, < kernel.pids_max, < kernel.threads-max",
			ExeEnv: exeenv.InContainer,
		},
		limited: false,
	}
)
