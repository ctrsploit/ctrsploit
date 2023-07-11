package prerequisite

import "strings"

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
	// TODO: get kernel release: uname -r
	release := "xxx"
	p.Satisfied = strings.Contains(release, p.ExpectedReleaser)
	return
}
