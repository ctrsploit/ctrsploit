package kernel

import (
	"fmt"
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
	return p.CheckTemplate(func() {
		u, err := uname.All()
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("getting uname: %w", err))
			return
		}
		log.Logger.Debugf("uname: %s", u)
		p.Satisfied = strings.Contains(u, p.ExpectedReleaser)
		return
	})
}
