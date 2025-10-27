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
	return p.CheckTemplate(func() (bool, error) {
		u, err := uname.All()
		if err != nil {
			p.Err = fmt.Errorf("failed to check [%s], caused by getting uname: %w", p.GetName(), err)
			return p.Satisfied, p.Err
		}
		log.Logger.Debugf("uname: %s", u)
		p.Satisfied = strings.Contains(u, p.ExpectedReleaser)
		return p.Satisfied, p.Err
	})
}
