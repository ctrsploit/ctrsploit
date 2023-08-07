package kernel

import (
	"github.com/ctrsploit/ctrsploit/log"
	"github.com/ctrsploit/ctrsploit/pkg/kernel/uname"
	"github.com/ctrsploit/ctrsploit/prerequisite"
	"strings"
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

func (p *Releaser) Check() (err error) {
	u, err := uname.All()
	if err != nil {
		return
	}
	log.Logger.Debugf("uname: %s", u)
	p.Satisfied = strings.Contains(u, p.ExpectedReleaser)
	return
}
