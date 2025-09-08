package kernel

import (
	"strings"

	"github.com/ctrsploit/ctrsploit/pkg/kernel/uname"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type Releaser struct {
	ExpectedReleaser string
	prerequisite.BasePrerequisite
}

var (
	ReleasedByLinuxkit = Releaser{
		ExpectedReleaser: "linuxkit",
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name: "linuxkit kernel",
			Info: "kernel released by linuxkit",
		},
	}
)

func (p *Releaser) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	u, err := uname.All()
	if err != nil {
		return false, err
	}
	log.Logger.Debugf("uname: %s", u)
	p.Satisfied = strings.Contains(u, p.ExpectedReleaser)
	p.Checked = true
	return p.Satisfied, nil
}
