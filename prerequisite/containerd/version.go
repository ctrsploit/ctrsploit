package containerd

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/version/containerd"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type VersionEqualTo struct {
	ExpectedVersion string
	prerequisite.BasePrerequisite
}

var VersionEqualToV2_1_0 = VersionEqualTo{
	ExpectedVersion: "v2.1.0",
	BasePrerequisite: prerequisite.BasePrerequisite{
		Name: "containerd==v2.1.0",
		Info: "Containerd version must be 2.1.0",
	},
}

func (p *VersionEqualTo) Check() (bool, error) {
	return p.CheckTemplate(func() (bool, error) {
		v, err := containerd.GetVersionBySock()
		if err != nil {
			p.Err = fmt.Errorf("failed to check [%s] caused by getting version from containerd: %v", p.GetName(), err)
			return p.Satisfied, p.Err
		}
		p.Satisfied = v.Version == p.ExpectedVersion
		return p.Satisfied, p.Err
	})
}
