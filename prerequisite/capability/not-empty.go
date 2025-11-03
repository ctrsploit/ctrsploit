package capability

import (
	"fmt"

	"github.com/containerd/containerd/pkg/cap"
	"github.com/ctrsploit/ctrsploit/pkg/capability"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type NotEmpty struct {
	CapType cap.Type
	Pid     []string
	prerequisite.BasePrerequisite
}

var (
	AmbNotEmpty = NotEmpty{
		CapType: cap.Ambient,
		Pid:     []string{"1", "self"},
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "CapAmb Not Empty",
			Info:   "Ambient Capabilities allow capabilities to be added to the Permitted set and Effective set of a process during normal user's execve.",
			ExeEnv: exeenv.InContainer,
		},
	}
)

func (p *NotEmpty) Check() (bool, error) {
	return p.CheckTemplate(func() {
		for _, pid := range p.Pid {
			caps, err := capability.GetCapabilityByPid(pid, p.CapType)
			if err != nil {
				p.Err = p.WrapErr(fmt.Errorf("getting capability for %s: %w", pid, err))
				return
			}
			capsParsed, _ := cap.FromBitmap(caps)
			if len(capsParsed) > 0 {
				p.Satisfied = true
				break
			}
		}
		return
	})
}
