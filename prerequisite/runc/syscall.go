package runc

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/version/runc"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type ensureCloned struct {
	prerequisite.BasePrerequisite
}

func (p *ensureCloned) Check() (bool, error) {
	return p.CheckTemplate(func() (bool, error) {
		var err error
		p.Satisfied, err = runc.StraceFGetSeals()
		if err != nil {
			p.Err = fmt.Errorf("failed to check [%s], caused by strace runc: %w", p.GetName(), err)
		}
		return p.Satisfied, p.Err
	})

}

var (
	// EnsureCloned indicates runc >= v1.0.0-rc7, <= runc-v1.1.15
	EnsureCloned = ensureCloned{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "runc with F_GET_SEALS",
			Info:   "runc called fcntl with F_GET_SEALS, indicates >= v1.0.0-rc7, <= runc-v1.1.15",
			ExeEnv: exeenv.InHost,
		},
	}
	// MaybeVulnerableTo0492BySyscall is a necessary but not sufficient condition for runc to be vulnerable to CVE-2019-5736
	MaybeVulnerableTo0492BySyscall = prerequisite.Not(&EnsureCloned)
)
