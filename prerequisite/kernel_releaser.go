package prerequisite

import (
	"github.com/ctrsploit/ctrsploit/log"
	"github.com/ctrsploit/ctrsploit/pkg/kernel/uname"
	"strings"
)

type KernelReleaser struct {
	ExpectedReleaser string
	BasePrerequisite
}

var (
	KernelReleasedByLinuxkit = KernelReleaser{
		ExpectedReleaser: "linuxkit",
		BasePrerequisite: BasePrerequisite{
			Name: "linuxkit kernel",
			Info: "kernel released by linuxkit",
		},
	}
)

func (p *KernelReleaser) Check() (err error) {
	u, err := uname.All()
	if err != nil {
		return
	}
	log.Logger.Debugf("uname: %s", u)
	p.Satisfied = strings.Contains(u, p.ExpectedReleaser)
	return
}
